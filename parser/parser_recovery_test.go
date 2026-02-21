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
			// The improved parser tolerates a missing semicolon and parses all three statements cleanly.
			expectedErrorCount:    0,
			expectedErrorSubstr:   []string{},
			expectedNumValidStmts: 3, // let x = 5, let y = 10, return y
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
			// After fixing the double-error bug only the primary "expected RightParenthesis" error is emitted.
			expectedErrorCount:    1,
			expectedErrorSubstr:   []string{"expected next token to be RightParenthesis"},
			expectedNumValidStmts: 1, // Only 'let d = 4;' should be parsed cleanly after recovery
		},
		{
			name: "Bad Let Statement Recovery",
			input: `main() -> {
                let = 6; // Missing identifier
                let e = 7;
            }`,
			expectedErrorCount:    1, // Error for missing identifier
			expectedErrorSubstr:   []string{"expected next token to be Identifier, got Assignment"},
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
			// Improved recovery: let f parses fine (no semicolon needed); let g fails; let h is
			// skipped by block-level recovery; return h is parsed.
			expectedErrorCount:    1,
			expectedErrorSubstr:   []string{"expected next token to be Assignment"},
			expectedNumValidStmts: 2, // let f = 10 and return h
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
