package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"strings"
	"testing"
)

func TestArrayLiteralParsingIntegration(t *testing.T) {
	tests := []struct {
		name             string
		input            string   // Input containing the array literal (e.g., in a let)
		expectedElements []string // String representations of expected element expressions
		expectedErrors   int
	}{
		{
			name:             "Empty Array",
			input:            `main() -> { let a = []; }`,
			expectedElements: []string{},
			expectedErrors:   0,
		},
		{
			name:             "Array of Integers",
			input:            `main() -> { let b = [1, 2, 3]; }`,
			expectedElements: []string{"1", "2", "3"},
			expectedErrors:   0,
		},
		{
			name:             "Array of Strings",
			input:            `main() -> { let c = ["one", "two"]; }`,
			expectedElements: []string{`"one"`, `"two"`},
			expectedErrors:   0,
		},
		{
			name:             "Array of Expressions",
			input:            `main() -> { let d = [1 + 2, x * y, getValue()]; }`,
			expectedElements: []string{"(1 + 2)", "(x * y)", "getValue()"},
			expectedErrors:   0,
		},
		{
			name:             "Array with Trailing Comma (Error)",
			input:            `main() -> { let e = [1, 2,]; }`, // Trailing comma currently causes error
			expectedElements: []string{"1", "2"},               // Parses elements before error
			expectedErrors:   1,                                // Error for unexpected trailing comma
		},
		{
			name:             "Array Missing Comma",
			input:            `main() -> { let f = [1 2]; }`,
			expectedElements: []string{"1"}, // Parses first element, expects comma or ']'
			expectedErrors:   1,             // Error for expecting ']' or ',' but got '2'
		},
		{
			name:             "Array Missing Closing Bracket",
			input:            `main() -> { let g = [1, 2; }`,
			expectedElements: []string{"1", "2"}, // Parses elements
			expectedErrors:   1,                  // Error for expecting ']' but got ';'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()

			if tt.expectedErrors > 0 {
				if len(p.Errors()) < tt.expectedErrors {
					t.Errorf("Expected at least %d parser errors, but got %d:", tt.expectedErrors, len(p.Errors()))
					for i, e := range p.Errors() {
						t.Errorf("  Error %d: %s", i+1, e)
					}
				}
				// Attempt to check parsed elements even if errors occurred, as recovery might allow partial parsing
			} else {
				checkParserErrors(t, p) // Fail on unexpected errors
			}

			if program.MainFunction == nil {
				// If the input itself was just the array, this check needs adjustment
				if !strings.HasPrefix(tt.input, "main()") {
					t.Logf("Skipping main function check for direct array input.")
				} else {
					t.Fatalf("ParseProgram() returned nil MainFunction")
				}
			}

			var arrayLit *ast.ArrayLiteral
			// Extract array literal - assumes it's the value in the first let statement of main
			if program.MainFunction != nil && program.MainFunction.Body != nil {
				if body, ok := program.MainFunction.Body.(*ast.BlockStatement); ok && len(body.Statements) > 0 {
					if letStmt, ok := body.Statements[0].(*ast.LetStatement); ok {
						if arr, ok := letStmt.Value.(*ast.ArrayLiteral); ok {
							arrayLit = arr
						}
					}
				}
			}

			if arrayLit == nil {
				// Maybe it was parsed directly as an expression statement?
				if program.MainFunction != nil && program.MainFunction.Body != nil {
					if body, ok := program.MainFunction.Body.(*ast.BlockStatement); ok && len(body.Statements) > 0 {
						if expStmt, ok := body.Statements[0].(*ast.ExpressionStatement); ok {
							if arr, ok := expStmt.Expression.(*ast.ArrayLiteral); ok {
								arrayLit = arr
							}
						}
					}
				}

				if arrayLit == nil {
					t.Fatalf("Failed to find *ast.ArrayLiteral node in AST")
				}
			}

			if len(arrayLit.Elements) != len(tt.expectedElements) {
				t.Errorf("Array element count mismatch. want=%d, got=%d", len(tt.expectedElements), len(arrayLit.Elements))
			} else {
				for i, expectedElemStr := range tt.expectedElements {
					actualElemStr := arrayLit.Elements[i].String()
					if actualElemStr != expectedElemStr {
						t.Errorf("Array element %d mismatch.\nWant: %s\nGot:  %s", i, expectedElemStr, actualElemStr)
					}
				}
			}
		})
	}
}

func TestIndexExpressionParsingIntegration(t *testing.T) {
	tests := []struct {
		name           string
		input          string // Input containing the index expression
		expectedLeft   string // String representation of the expression being indexed
		expectedIndex  string // String representation of the index expression
		expectedErrors int
	}{
		{
			name:           "Index with Integer",
			input:          `main() -> { myArray[0]; }`,
			expectedLeft:   "myArray",
			expectedIndex:  "0",
			expectedErrors: 0,
		},
		{
			name:           "Index with Identifier",
			input:          `main() -> { data[key]; }`,
			expectedLeft:   "data",
			expectedIndex:  "key",
			expectedErrors: 0,
		},
		{
			name:           "Index with Expression",
			input:          `main() -> { values[i + 1]; }`,
			expectedLeft:   "values",
			expectedIndex:  "(i + 1)",
			expectedErrors: 0,
		},
		{
			name:           "Index with Function Call",
			input:          `main() -> { results[getIndex()]; }`,
			expectedLeft:   "results",
			expectedIndex:  "getIndex()",
			expectedErrors: 0,
		},
		{
			name:           "Chained Indexing",
			input:          `main() -> { matrix[row][col]; }`,
			expectedLeft:   "(matrix[row])", // Left side of the *outer* index is the inner index expr
			expectedIndex:  "col",
			expectedErrors: 0,
		},
		{
			name:           "Missing Index Expression",
			input:          `main() -> { array[]; }`,
			expectedLeft:   "array",
			expectedIndex:  "", // Fails to parse index
			expectedErrors: 1,  // Error for missing expression inside []
		},
		{
			name:           "Missing Closing Bracket",
			input:          `main() -> { array[0; }`,
			expectedLeft:   "array",
			expectedIndex:  "0", // Parses index
			expectedErrors: 1,   // Error for expecting ']'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()

			if tt.expectedErrors > 0 {
				if len(p.Errors()) < tt.expectedErrors {
					t.Errorf("Expected at least %d parser errors, but got %d:", tt.expectedErrors, len(p.Errors()))
					for i, e := range p.Errors() {
						t.Errorf("  Error %d: %s", i+1, e)
					}
				}
				return // Don't check AST if errors expected
			} else {
				checkParserErrors(t, p) // Fail on unexpected errors
			}

			if program.MainFunction == nil {
				t.Fatalf("ParseProgram() returned nil MainFunction")
			}
			mainBody, ok := program.MainFunction.Body.(*ast.BlockStatement)
			if !ok || len(mainBody.Statements) != 1 {
				t.Fatalf("Expected main body with 1 statement, got %T with %d statements", program.MainFunction.Body, len(mainBody.Statements))
			}
			exprStmt, ok := mainBody.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Expected statement to be *ast.ExpressionStatement, got %T", mainBody.Statements[0])
			}
			indexExpr, ok := exprStmt.Expression.(*ast.IndexExpression)
			if !ok {
				t.Fatalf("Expression is not *ast.IndexExpression, got %T", exprStmt.Expression)
			}

			// Check Left Expression
			if indexExpr.Left == nil {
				t.Errorf("IndexExpression.Left is nil")
			} else {
				actualLeftStr := indexExpr.Left.String()
				if actualLeftStr != tt.expectedLeft {
					t.Errorf("IndexExpression.Left mismatch.\nWant: %s\nGot:  %s", tt.expectedLeft, actualLeftStr)
				}
			}

			// Check Index Expression
			if indexExpr.Index == nil {
				if tt.expectedIndex != "" {
					t.Errorf("IndexExpression.Index is nil, expected %s", tt.expectedIndex)
				}
			} else {
				actualIndexStr := indexExpr.Index.String()
				if actualIndexStr != tt.expectedIndex {
					t.Errorf("IndexExpression.Index mismatch.\nWant: %s\nGot:  %s", tt.expectedIndex, actualIndexStr)
				}
			}
		})
	}
}
