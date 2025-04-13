// compiler/generator/codegen_ir_structure_test.go
package generator

import (
	"compiler/lexer" // Need imports for the test setup helpers
	"compiler/parser"
	"fmt"
	// "github.com/llir/llvm/ir" // Keep ir import for types/constants etc.
	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"strings"
	"testing"
)

func TestCodeGenIRStructureBasic(t *testing.T) {
	// Test cases remain the same...
	tests := []struct {
		name                  string
		input                 string
		expectGenerationError bool   // Whether generateIRForProgram itself is expected to error
		expectStringerError   bool   // Whether module.String() is expected to error/panic due to invalid IR
		errorSubstring        string // Optional substring to check in the error message
	}{
		{
			name:                  "Simple Valid Function",
			input:                 `main() -> { return 1; }`,
			expectGenerationError: false,
			expectStringerError:   false,
		},
		{
			name:                  "Valid If/Else with Returns",
			input:                 `main() -> { if (1 > 0) { return 1; } else { return 0; } }`,
			expectGenerationError: false,
			expectStringerError:   false, // Both paths terminate
		},
		{
			name:                  "Missing Return in Main", // generateIR helper should add one
			input:                 `main() -> { let x = 5; }`,
			expectGenerationError: false,
			expectStringerError:   false, // Helper adds default return
		},
		{
			name: "Missing Return in If Branch (Causes Unterminated Block)",
			input: `
                 main() -> {
                     if (1 > 0) {
                         let y = 10; // No return here, block is unterminated before merge
                     } else {
                         return 0; // This path is terminated
                     }
                     // Code after the if/else implicitly needs a return if reachable
                     // return 999; // Adding this would make it valid again
                 }
             `,
			expectGenerationError: false,                // Generator visitor might succeed
			expectStringerError:   true,                 // String() should fail/panic on unterminated block
			errorSubstring:        "missing terminator", // Expect error msg to mention terminator
		},
		{
			name: "Function Declared but Not Defined",
			input: `
                  function externalFunc(a); // Declaration syntax might fail parser
                  main() -> { return externalFunc(5); }
              `,
			expectGenerationError: true,                                         // Expect error from parser or codegen visitor
			expectStringerError:   false,                                        // Won't reach stringer if generation fails
			errorSubstring:        "function externalFunc was not pre-declared", // Or parser error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a modified helper that returns the module object *before* stringifying
			// to separate generation errors from stringer errors.

			cg := NewCodeGenerator() // Create fresh generator for each test
			var genErr error
			var parseErrs []string

			// --- Parsing ---
			l, errLex := lexer.NewLexerFromString(tt.input)
			if errLex != nil {
				genErr = fmt.Errorf("lexer error: %v", errLex)
			} else {
				p := parser.NewParser(l)
				prog := p.ParseProgram()
				parseErrs = p.Errors()
				// Allow parser errors if test expects generation error (often related)
				if len(parseErrs) > 0 && !tt.expectGenerationError {
					genErr = fmt.Errorf("parser errors: %v", parseErrs)
				} else if len(parseErrs) > 0 && tt.expectGenerationError {
					t.Logf("Parser errors occurred as potentially expected: %v", parseErrs)
					// Set genErr to potentially check error message later
					genErr = fmt.Errorf("parser errors: %v", parseErrs)
				} else if prog == nil {
					genErr = fmt.Errorf("parser returned nil program")
				} else {
					// --- Code Generation (Visiting) ---
					if errSetup := setupArrayStructType(cg); errSetup != nil {
						genErr = fmt.Errorf("Failed to setup Array struct type: %v", errSetup)
					} else {
						acceptErr := prog.Accept(cg) // Visit the program
						if acceptErr != nil {
							genErr = fmt.Errorf("visitor error: %v", acceptErr)
						}
					}
				}
			}

			// --- Check Generation Errors ---
			if tt.expectGenerationError {
				if genErr == nil {
					t.Errorf("Expected an error during parsing/generation, but got nil.")
				} else {
					t.Logf("Received expected error during parsing/generation: %v", genErr)
					if tt.errorSubstring != "" && !strings.Contains(genErr.Error(), tt.errorSubstring) {
						t.Errorf("Generation error message mismatch.\nExpected substring: %s\nGot error: %v", tt.errorSubstring, genErr)
					}
				}
				return // Stop if generation error was expected
			} else if genErr != nil {
				t.Fatalf("Parsing/Generation failed unexpectedly: %v", genErr)
			}

			// --- Add Terminators (Test Harness Logic) ---
			if mainFunc, ok := cg.Functions["main"]; ok {
				if len(mainFunc.Blocks) > 0 {
					lastBlock := mainFunc.Blocks[len(mainFunc.Blocks)-1]
					if lastBlock.Term == nil {
						t.Logf("Test harness adding default terminator to main's last block (%s)", lastBlock.LocalIdent.Name())
						// Add default return based on main's signature
						if mainFunc.Sig.RetType.Equal(types.I32) {
							lastBlock.NewRet(constant.NewInt(types.I32, 0))
						} else if mainFunc.Sig.RetType.Equal(types.Void) {
							lastBlock.NewRet(nil)
						} else {
							lastBlock.NewRet(constant.NewZeroInitializer(mainFunc.Sig.RetType))
						}
					}
				}
			}

			// --- IR Stringification and Validation ---
			var irString string
			var stringerErr error
			panicPayload := Catcher(func() {
				irString = cg.Module.String()
			})()

			if panicPayload != nil {
				stringerErr = fmt.Errorf("panic during Module.String(): %v", panicPayload)
			}

			if tt.expectStringerError {
				if stringerErr == nil {
					t.Errorf("Expected Module.String() to error or panic due to invalid IR, but it did not.\nIR:\n%s", irString)
				} else {
					t.Logf("Received expected error/panic from Module.String(): %v", stringerErr)
					if tt.errorSubstring != "" && !strings.Contains(stringerErr.Error(), tt.errorSubstring) {
						t.Errorf("Stringer error message mismatch.\nExpected substring: %s\nGot error: %v", tt.errorSubstring, stringerErr)
					}
				}
			} else {
				// Expected valid IR, so String() should succeed
				if stringerErr != nil {
					t.Errorf("Module.String() failed unexpectedly: %v\nIR generated before error (if any):\n%s", stringerErr, irString)
				} else {
					t.Logf("Module.String() succeeded as expected.")
					// t.Logf("Generated Valid IR:\n%s", irString)

					// --- Use asm.ParseAssembly for final validation ---
					_, errParse := asm.Parse("", strings.NewReader(irString)) // Use asm.Parse
					if errParse != nil {
						t.Errorf("Failed to parse the supposedly valid IR generated by String() back using asm.Parse: %v\nIR:\n%s", errParse, irString)
					} else {
						t.Logf("Generated IR successfully parsed back via asm.Parse.")
					}
				}
			}
		})
	}
}

// Catcher remains the same
func Catcher(f func()) func() interface{} {
	return func() (payload interface{}) {
		defer func() {
			payload = recover()
		}()
		f()
		return
	}
}
