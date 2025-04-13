package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"strings"
	"testing"
)

func TestStatementRecoveryUnit(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectedErrorCount    int
		expectedErrorSubstr   []string // Substrings to look for in error messages
		expectedNumValidStmts int      // Number of statements expected in the *recovered* block/program
	}{
		{
			name: "Missing Semicolon Recovery",
			input: `main() -> {
                let x = 5 // Missing semicolon
                let y = 10;
                return y;
            }`,
			expectedErrorCount:    1,                                                                          // Expect an error for the missing '=' after let x
			expectedErrorSubstr:   []string{"Expected '=' operator after let statement identifier near line"}, // Error from let parsing failure
			expectedNumValidStmts: 2,                                                                          // Should parse 'let y = 10;' and 'return y;'
		},
		{
			name: "Invalid Expression Start Recovery",
			input: `main() -> {
                let a = 1;
                + 5; // Invalid start of expression
                let b = 2;
            }`,
			expectedErrorCount:    1, // Error for the '+'
			expectedErrorSubstr:   []string{"Operator '+' cannot start an expression"},
			expectedNumValidStmts: 2, // Should parse 'let a = 1;' and 'let b = 2;'
		},
		{
			name: "Mismatched Parentheses Recovery",
			input: `main() -> {
                let c = (5 + 3; // Missing closing parenthesis
                let d = 4;
            }`,
			expectedErrorCount:    1, // Error for missing ')'
			expectedErrorSubstr:   []string{"Expected ')' after grouped expression"},
			expectedNumValidStmts: 1, // Only 'let d = 4;' should be parsed cleanly after recovery
		},
		{
			name: "Bad Let Statement Recovery",
			input: `main() -> {
                let = 6; // Missing identifier
                let e = 7;
            }`,
			expectedErrorCount:    1, // Error for missing identifier
			expectedErrorSubstr:   []string{"expected token to be Identifier, got Assignment instead"},
			expectedNumValidStmts: 1, // Should parse 'let e = 7;'
		},
		{
			name: "Multiple Errors Recovery",
			input: `main() -> {
                let f = 10
                let g // Missing assignment
                let h = f + ; // Missing operand
                return h;
            }`,
			expectedErrorCount:    3,                                                                                                // Missing ';', missing '=', missing operand
			expectedErrorSubstr:   []string{"Expected '=' operator", "Unexpected token ';' (Semicolon) cannot start an expression"}, // Errors might cascade or manifest differently
			expectedNumValidStmts: 1,                                                                                                // Only the return statement is likely fully valid after recovery
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

			// 1. Check Error Count
			if len(p.Errors()) != tt.expectedErrorCount {
				t.Errorf("Expected %d parser errors, but got %d:", tt.expectedErrorCount, len(p.Errors()))
				for i, e := range p.Errors() {
					t.Errorf("  Error %d: %s", i+1, e)
				}
			}

			// 2. Check Error Content (if specific substrings are expected)
			if len(tt.expectedErrorSubstr) > 0 {
				errorString := strings.Join(p.Errors(), "\n")
				for _, sub := range tt.expectedErrorSubstr {
					if !strings.Contains(errorString, sub) {
						t.Errorf("Expected error message to contain '%s', but got:\n%s", sub, errorString)
					}
				}
			}

			// 3. Check Number of Successfully Parsed Statements in the main block
			if program == nil || program.MainFunction == nil || program.MainFunction.Body == nil {
				if tt.expectedNumValidStmts > 0 {
					t.Errorf("Expected %d valid statements, but program/main/body is nil", tt.expectedNumValidStmts)
				}
				// If everything is nil and we expected 0 statements, that's okay.
			} else {
				mainBody, ok := program.MainFunction.Body.(*ast.BlockStatement)
				if !ok {
					t.Errorf("Main function body is not a BlockStatement, got %T", program.MainFunction.Body)
				} else {
					actualValidStmts := 0
					if mainBody != nil { // Check if mainBody itself isn't nil
						actualValidStmts = len(mainBody.Statements)
					}

					if actualValidStmts != tt.expectedNumValidStmts {
						t.Errorf("Expected %d valid statements after recovery, but got %d", tt.expectedNumValidStmts, actualValidStmts)
						t.Logf("Recovered Statements (%d):", actualValidStmts)
						if mainBody != nil {
							for i, stmt := range mainBody.Statements {
								t.Logf("  %d: %s (%T)", i+1, stmt.String(), stmt)
							}
						}

					}
				}
			}
		})
	}
}
