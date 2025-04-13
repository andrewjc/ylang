package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestConditionalParsingIntegration(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedCond   string // String representation of the condition expression
		expectedCons   string // String representation of the consequence block/expr
		expectedAlt    string // String representation of the alternative block/expr, "" if none
		expectedErrors int
		isAltElseIf    bool // Is the alternative expected to be another IfStatement?
	}{
		{
			name:           "Simple If Statement",
			input:          `main() -> { if (x > y) { return x; } }`,
			expectedCond:   "(x > y)",
			expectedCons:   "{\n    return x;\n}",
			expectedAlt:    "",
			expectedErrors: 0,
			isAltElseIf:    false,
		},
		{
			name:           "If/Else Statement",
			input:          `main() -> { if (x < 0) { return -1; } else { return 1; } }`,
			expectedCond:   "(x < 0)",
			expectedCons:   "{\n    return -1;\n}", // Assuming prefix '-' is handled eventually
			expectedAlt:    "{\n    return 1;\n}",
			expectedErrors: 0,
			isAltElseIf:    false,
		},
		{
			name: "If/Else If/Else Statement",
			input: `main() -> {
                if (x == 0) {
                    print("zero");
                } else if (x > 0) {
                    print("positive");
                } else {
                    print("negative");
                }
            }`,
			expectedCond: "(x == 0)", // Needs == operator support
			expectedCons: "{\n    print(\"zero\");\n}",
			// The Alternative *is* another IfStatement node
			expectedAlt: `if (x > 0) {
        print("positive");
    } else {
        print("negative");
    }`, // String of the nested IfStatement
			expectedErrors: 0,
			isAltElseIf:    true,
		},
		{
			name:           "If without braces (Not supported by current parser, should error)",
			input:          `main() -> { if (a > b) return a; }`,
			expectedCond:   "(a > b)",
			expectedCons:   "", // Fails to parse consequence
			expectedAlt:    "",
			expectedErrors: 1, // Error for expecting '{'
			isAltElseIf:    false,
		},
		{
			name:           "Missing condition",
			input:          `main() -> { if () { print("oops"); } }`,
			expectedCond:   "",                           // Fails parsing condition
			expectedCons:   "{\n    print(\"oops\");\n}", // May parse block depending on recovery
			expectedAlt:    "",
			expectedErrors: 1, // Error for missing expression in condition
			isAltElseIf:    false,
		},
		{
			name:           "Missing consequence block",
			input:          `main() -> { if (c > 0) else { print("alt"); } }`,
			expectedCond:   "(c > 0)",
			expectedCons:   "",                          // Fails parsing consequence
			expectedAlt:    "{\n    print(\"alt\");\n}", // May parse alt depending on recovery
			expectedErrors: 1,                           // Error for expecting '{' after condition
			isAltElseIf:    false,
		},
		{
			name: "Nested If Statements",
			input: `main() -> {
                if (a) {
                    if (b) {
                        return 1;
                    } else {
                        return 0;
                    }
                } else {
                    return -1;
                }
            }`,
			expectedCond:   "a",
			expectedCons:   "{\n    if b {\n        return 1;\n    } else {\n        return 0;\n    }\n}", // String representation of nested if
			expectedAlt:    "{\n    return -1;\n}",
			expectedErrors: 0,
			isAltElseIf:    false,
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
				// Don't check AST details if errors were expected
				return
			} else {
				checkParserErrors(t, p) // Fail if unexpected errors occurred
			}

			if program.MainFunction == nil {
				t.Fatalf("ParseProgram() returned nil MainFunction")
			}
			mainBody, ok := program.MainFunction.Body.(*ast.BlockStatement)
			if !ok {
				t.Fatalf("Main function body is not *ast.BlockStatement, got %T", program.MainFunction.Body)
			}
			if len(mainBody.Statements) != 1 {
				t.Fatalf("Expected 1 statement (the if statement) in main body, got %d", len(mainBody.Statements))
			}

			// The If statement is likely wrapped in an ExpressionStatement
			exprStmt, ok := mainBody.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				// Or maybe parseStatement returns the IfStatement directly? Check parseStatement logic.
				// Let's assume it returns the IfStatement directly for now if it implements Statement
				ifStmtDirect, okDirect := mainBody.Statements[0].(*ast.IfStatement)
				if !okDirect {
					t.Fatalf("Expected statement to be *ast.ExpressionStatement or *ast.IfStatement, got %T", mainBody.Statements[0])
				}
				exprStmt = &ast.ExpressionStatement{Expression: ifStmtDirect} // Wrap it for consistent checking
			}

			ifStmt, ok := exprStmt.Expression.(*ast.IfStatement)
			if !ok {
				t.Fatalf("Expression is not *ast.IfStatement, got %T", exprStmt.Expression)
			}

			// Check Condition
			if ifStmt.Condition == nil {
				if tt.expectedCond != "" {
					t.Errorf("Condition is nil, expected %s", tt.expectedCond)
				}
			} else {
				actualCondStr := ifStmt.Condition.String()
				if actualCondStr != tt.expectedCond {
					t.Errorf("Condition mismatch.\nWant: %s\nGot:  %s", tt.expectedCond, actualCondStr)
				}
			}

			// Check Consequence
			if ifStmt.Consequence == nil {
				if tt.expectedCons != "" {
					t.Errorf("Consequence is nil, expected %s", tt.expectedCons)
				}
			} else {
				actualConsStr := ifStmt.Consequence.String()
				if actualConsStr != tt.expectedCons {
					t.Errorf("Consequence mismatch.\nWant: %s\nGot:  %s", tt.expectedCons, actualConsStr)
				}
			}

			// Check Alternative
			if tt.expectedAlt == "" {
				if ifStmt.Alternative != nil {
					t.Errorf("Alternative is not nil (%s), expected nil", ifStmt.Alternative.String())
				}
			} else {
				if ifStmt.Alternative == nil {
					t.Errorf("Alternative is nil, expected %s", tt.expectedAlt)
				} else {
					_, altIsIf := ifStmt.Alternative.(*ast.IfStatement)
					if altIsIf != tt.isAltElseIf {
						t.Errorf("Alternative type mismatch. Expected 'else if'=%v, but got 'else if'=%v (Type: %T)", tt.isAltElseIf, altIsIf, ifStmt.Alternative)
					}

					actualAltStr := ifStmt.Alternative.String()
					// Normalize expected string for nested if/else if
					expectedAltStr := tt.expectedAlt
					if tt.isAltElseIf {
						// The String() method for IfStatement might handle indentation differently
						// when nested within another 'else'. We may need to adjust the expected string
						// or make the String() method more consistent.
						// For now, compare directly.
					}

					if actualAltStr != expectedAltStr {
						t.Errorf("Alternative mismatch.\nWant: %s\nGot:  %s", expectedAltStr, actualAltStr)
					}
				}
			}
		})
	}
}

func TestTernaryExpressionParsingUnit(t *testing.T) {
	tests := []struct {
		name            string
		input           string // Input containing the ternary expression (e.g., in a let or return)
		expectedTernary string // String representation of the ternary AST node
		expectedErrors  int
		nodeExtractor   func(*ast.Program) ast.ExpressionNode // Function to find the ternary node
	}{
		{
			name:            "Traditional Ternary",
			input:           `main() -> { let x = condition ? trueValue : falseValue; }`,
			expectedTernary: "(condition ? trueValue : falseValue)",
			expectedErrors:  0,
			nodeExtractor: func(p *ast.Program) ast.ExpressionNode {
				// Extract from let statement value
				if stmt := p.MainFunction.Body.(*ast.BlockStatement).Statements[0].(*ast.LetStatement); stmt != nil {
					return stmt.Value
				}
				return nil
			},
		},
		{
			name:            "Lambda Style Ternary",
			input:           `main() -> { return check -> resultA : resultB; }`,
			expectedTernary: "(check -> resultA : resultB)",
			expectedErrors:  0,
			nodeExtractor: func(p *ast.Program) ast.ExpressionNode {
				// Extract from return statement value
				if stmt := p.MainFunction.Body.(*ast.BlockStatement).Statements[0].(*ast.ReturnStatement); stmt != nil {
					return stmt.ReturnValue
				}
				return nil
			},
		},
		{
			name:            "Inline If/Else Ternary",
			input:           `main() -> { y = value if flag else defaultValue; }`,
			expectedTernary: "(value if flag else defaultValue)", // Note: this is the expr part, not the assignment
			expectedErrors:  0,
			nodeExtractor: func(p *ast.Program) ast.ExpressionNode {
				// Extract RHS of assignment
				if exprStmt := p.MainFunction.Body.(*ast.BlockStatement).Statements[0].(*ast.ExpressionStatement); exprStmt != nil {
					if assignExpr, ok := exprStmt.Expression.(*ast.AssignmentExpression); ok {
						return assignExpr.Right // The ternary is the right side of the assignment
					}
				}
				return nil
			},
		},
		// Add error cases if needed, e.g., missing parts of the ternary
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
					t.Errorf("Expected at least %d parser errors, but got %d", tt.expectedErrors, len(p.Errors()))
				}
				return
			} else {
				checkParserErrors(t, p)
			}

			ternaryNode := tt.nodeExtractor(program)
			if ternaryNode == nil {
				t.Fatalf("Failed to extract ternary node from AST")
			}

			actualTernaryStr := ternaryNode.String()
			if actualTernaryStr != tt.expectedTernary {
				t.Errorf("Ternary expression string representation mismatch.\nWant: %s\nGot:  %s", tt.expectedTernary, actualTernaryStr)
			}

			// Additionally, verify the node type if possible
			// Example check for traditional ternary:
			// if tt.name == "Traditional Ternary" {
			//  if _, ok := ternaryNode.(*ast.TraditionalTernaryExpression); !ok {
			//      t.Errorf("Expected node type *ast.TraditionalTernaryExpression, got %T", ternaryNode)
			//  }
			// }
			// Add similar checks for other ternary types (LambdaStyleTernaryExpression, InlineIfElseTernaryExpression)
		})
	}
}
