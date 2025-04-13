package ast

import (
	"compiler/lexer"
	"compiler/parser"
	"strings"
	"testing"
)

func TestASTStringificationUnit(t *testing.T) {
	tests := []struct {
		name              string
		input             string // Input code snippet
		expectedStr       string // Expected output of String() on the primary node
		expectedIndentStr string // Expected output of StringIndent(1)
		nodeExtractor     func(*Program) Node
	}{
		{
			name:              "Let Statement",
			input:             `main() -> { let myVar = 10; }`,
			expectedStr:       "let myVar = 10;",
			expectedIndentStr: "    let myVar = 10;", // StringIndent adds indent
			nodeExtractor: func(prog *Program) Node {
				return prog.MainFunction.Body.(*BlockStatement).Statements[0]
			},
		},
		{
			name:              "Return Statement",
			input:             `main() -> { return myVar + 1; }`,
			expectedStr:       "return (myVar + 1);",
			expectedIndentStr: "    return (myVar + 1);",
			nodeExtractor: func(prog *Program) Node {
				return prog.MainFunction.Body.(*BlockStatement).Statements[0]
			},
		},
		{
			name:              "Infix Expression",
			input:             `main() -> { 5 * (10 - 2); }`,
			expectedStr:       "(5 * (10 - 2))",     // ExpressionStatement.String() calls Expression.String()
			expectedIndentStr: "    (5 * (10 - 2))", // StringIndent on ExprStmt adds indent
			nodeExtractor: func(prog *Program) Node {
				return prog.MainFunction.Body.(*BlockStatement).Statements[0] // The ExpressionStatement
			},
		},
		{
			name:              "Block Statement",
			input:             `main() -> {{ let a = 1; return a; }}`,
			expectedStr:       "{\n    let a = 1;\n    return a;\n}",                 // String() calls StringIndent(0)
			expectedIndentStr: "    {\n        let a = 1;\n        return a;\n    }", // StringIndent(1)
			nodeExtractor: func(prog *Program) Node {
				// The inner block is wrapped in an ExpressionStatement
				exprStmt := prog.MainFunction.Body.(*BlockStatement).Statements[0].(*ExpressionStatement)
				return exprStmt.Expression
			},
		},
		{
			name:              "If Statement",
			input:             `main() -> { if (x > 0) { print(x); } else { print("neg"); } }`,
			expectedStr:       "if (x > 0) {\n    print(x);\n} else {\n    print(\"neg\");\n}",
			expectedIndentStr: "    if (x > 0) {\n        print(x);\n    } else {\n        print(\"neg\");\n    }",
			nodeExtractor: func(prog *Program) Node {
				exprStmt := prog.MainFunction.Body.(*BlockStatement).Statements[0].(*ExpressionStatement)
				return exprStmt.Expression // The IfStatement
			},
		},
		{
			name:              "Function Definition",
			input:             `function add(a, b) -> a + b;`,
			expectedStr:       "add(a, b) -> (a + b)",     // No indent for top-level String()
			expectedIndentStr: "    add(a, b) -> (a + b)", // StringIndent adds prefix
			nodeExtractor:     func(prog *Program) Node { return prog.Functions[0] },
		},
		{
			name:              "Lambda Expression",
			input:             `main() -> { let l = (x, y) -> { return x * y; }; }`,
			expectedStr:       "(x, y) -> {\n    return (x * y);\n}",             // String() calls StringIndent(0)
			expectedIndentStr: "    (x, y) -> {\n        return (x * y);\n    }", // StringIndent(1) adds prefix
			nodeExtractor: func(prog *Program) Node {
				letStmt := prog.MainFunction.Body.(*BlockStatement).Statements[0].(*LetStatement)
				return letStmt.Value // The LambdaExpression
			},
		},
		{
			name:              "Array Literal",
			input:             `main() -> { [1, "two", call()]; }`,
			expectedStr:       `[1, "two", call()]`,
			expectedIndentStr: `[1, "two", call()]`, // Array literal doesn't typically indent itself
			nodeExtractor: func(prog *Program) Node {
				exprStmt := prog.MainFunction.Body.(*BlockStatement).Statements[0].(*ExpressionStatement)
				return exprStmt.Expression // The ArrayLiteral
			},
		},
		// Add more tests for Call, Index, Assignment, etc.

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := parser.NewParser(l)
			program := p.ParseProgram()
			checkParserErrorsAST(t, p) // Use AST-specific error check helper

			node := tt.nodeExtractor(program)
			if node == nil {
				t.Fatalf("Node extraction failed for input: %s", tt.input)
			}

			// Test String()
			actualStr := node.String()
			// Normalize whitespace for comparison robustness, especially with block formatting
			normalize := func(s string) string {
				s = strings.ReplaceAll(s, "\r\n", "\n") // Normalize line endings
				// Simple field splitting might be too aggressive, compare with careful whitespace handling
				return s
			}
			if normalize(actualStr) != normalize(tt.expectedStr) {
				t.Errorf("String() mismatch.\nInput: %s\nWant:\n%s\nGot:\n%s", tt.input, tt.expectedStr, actualStr)
			}

			// Test StringIndent(1)
			stringerIndent, ok := node.(interface{ StringIndent(int) string })
			if !ok {
				// If StringIndent isn't implemented, compare String() output with indent added manually
				expectedManualIndent := "    " + tt.expectedStr // Simple indent for non-block nodes
				if normalize(actualStr) != normalize(expectedManualIndent) && tt.expectedIndentStr != "" && normalize(actualStr) != normalize(tt.expectedIndentStr) {
					// Check if the non-indented string matches the expected indented string
					// This handles cases where String() might already include some base indentation
					// or StringIndent is expected but missing.
					t.Logf("StringIndent(1) not implemented on %T, comparing String() output.", node)
					// Decide if this is an error based on whether StringIndent was expected.
					if tt.expectedIndentStr != tt.expectedStr { // If indentation was expected
						t.Errorf("StringIndent(1) expected but not implemented or mismatch.\nWant Indented:\n%s\nGot (from String()):\n%s", tt.expectedIndentStr, actualStr)
					}
				}
			} else {
				actualIndentStr := stringerIndent.StringIndent(1)
				if normalize(actualIndentStr) != normalize(tt.expectedIndentStr) {
					t.Errorf("StringIndent(1) mismatch.\nInput: %s\nWant:\n%s\nGot:\n%s", tt.input, tt.expectedIndentStr, actualIndentStr)
				}
			}
		})
	}
}
