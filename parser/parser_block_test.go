package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestBlockStatementParsingUnit(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedNumStmts int
		expectedStmts    []string // Expected string representation of statements inside the block
	}{
		{
			name:             "Empty Block",
			input:            `main() -> {{}}`, // Main requires a block, then test nested empty block
			expectedNumStmts: 0,
			expectedStmts:    []string{},
		},
		{
			name:             "Single Statement Block",
			input:            `main() -> {{ let x = 5; }}`,
			expectedNumStmts: 1,
			expectedStmts:    []string{"let x = 5;"},
		},
		{
			name:             "Multiple Statement Block",
			input:            `main() -> {{ let x = 5; return x + 1; }}`,
			expectedNumStmts: 2,
			expectedStmts:    []string{"let x = 5;", "return (x + 1);"},
		},
		{
			name: "Nested Blocks",
			input: `main() -> {
                {
                    let outer = 1;
                    {
                        let inner = 2;
                        outer + inner;
                    }
                    return outer;
                }
            }`,
			expectedNumStmts: 3, // outer let, inner block, outer return
			expectedStmts: []string{
				"let outer = 1;",
				"{\n        let inner = 2;\n        (outer + inner);\n    }", // StringIndent applied
				"return outer;",
			},
		},
		{
			name:             "Block without Semicolons (where optional)",
			input:            `main() -> {{ let x = 1\n x = x + 1\n return x }}`, // Assumes newline can act as terminator sometimes
			expectedNumStmts: 3,
			expectedStmts:    []string{"let x = 1;", "x = (x + 1);", "return x;"}, // Semicolons are added by String() where needed
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
			checkParserErrors(t, p)

			if program.MainFunction == nil {
				t.Fatalf("ParseProgram() returned nil MainFunction")
			}

			mainBody, ok := program.MainFunction.Body.(*ast.BlockStatement)
			if !ok {
				t.Fatalf("Main function body is not *ast.BlockStatement, got %T", program.MainFunction.Body)
			}

			// Find the target block (could be the main body itself or nested)
			var targetBlock *ast.BlockStatement
			if tt.name == "Empty Block" || tt.name == "Single Statement Block" || tt.name == "Multiple Statement Block" || tt.name == "Block without Semicolons (where optional)" {
				if len(mainBody.Statements) != 1 {
					t.Fatalf("Expected 1 statement in main body (the nested block), got %d", len(mainBody.Statements))
				}
				targetBlock, ok = mainBody.Statements[0].(*ast.BlockStatement)
				if !ok {
					t.Fatalf("Expected main body's first statement to be *ast.BlockStatement, got %T", mainBody.Statements[0])
				}
			} else if tt.name == "Nested Blocks" {
				// The target is the first block within the main block
				if len(mainBody.Statements) != 1 {
					t.Fatalf("Expected 1 statement in main body (the outer block), got %d", len(mainBody.Statements))
				}
				exprStmt, ok := mainBody.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf("Expected ExpressionStatement wrapping block, got %T", mainBody.Statements[0])
				}
				targetBlock, ok = exprStmt.Expression.(*ast.BlockStatement)
				if !ok {
					t.Fatalf("Expected main body's first statement's expression to be *ast.BlockStatement, got %T", exprStmt.Expression)
				}
			} else {
				t.Fatalf("Test case '%s' not handled for target block selection", tt.name)
			}

			if len(targetBlock.Statements) != tt.expectedNumStmts {
				t.Errorf("Block statement count mismatch. want=%d, got=%d", tt.expectedNumStmts, len(targetBlock.Statements))
				t.Logf("Got Statements:\n")
				for _, s := range targetBlock.Statements {
					t.Logf("- %s (%T)\n", s.String(), s)
				}

			}

			for i, expectedStmtStr := range tt.expectedStmts {
				if i >= len(targetBlock.Statements) {
					t.Errorf("Missing expected statement %d: %s", i, expectedStmtStr)
					continue
				}
				actualStmtStr := targetBlock.Statements[i].String() // Use default String() which should match expected format closely
				if actualStmtStr != expectedStmtStr {
					t.Errorf("Statement %d mismatch.\nWant: %s\nGot:  %s", i, expectedStmtStr, actualStmtStr)
				}
			}
		})
	}
}
