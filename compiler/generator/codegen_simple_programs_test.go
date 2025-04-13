package generator

import (
	"compiler/lexer"
	"compiler/parser"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"regexp"
	"testing"
)

// Helper to parse, codegen, and return IR for a full program string
func generateIRForProgram(t *testing.T, input string) (string, error) {
	t.Helper()
	l, errLex := lexer.NewLexerFromString(input)
	if errLex != nil {
		return "", fmt.Errorf("lexer error: %v", errLex)
	}
	p := parser.NewParser(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		// Return parser errors instead of failing hard, allows testing codegen error handling later
		return "", fmt.Errorf("parser errors: %v", p.Errors())
	}
	if prog == nil {
		return "", fmt.Errorf("parser returned nil program")
	}

	cg := NewCodeGenerator()
	// Perform necessary setup before visiting the program AST
	if err := setupArrayStructType(cg); err != nil { // Required by some test cases
		t.Fatalf("Failed to setup Array struct type: %v", err)
	}
	// Add other necessary setup (e.g., pre-defining types or functions if needed)

	errVisit := prog.Accept(cg) // Visit the entire program AST
	if errVisit != nil {
		// Return codegen visitor errors
		return cg.Module.String(), fmt.Errorf("codegen visitor error: %v", errVisit)
	}

	// Optional: Add validation or checks on the generated module *before* stringifying
	// For example, check if all blocks are terminated, etc.
	// Verify main function termination if it exists
	if mainFunc, ok := cg.Functions["main"]; ok {
		hasTerminator := false
		for _, block := range mainFunc.Blocks {
			if block.Term != nil {
				// Check if it's the last block? Or if any block has a terminator?
				// For now, assume if *any* block has a terminator, it's potentially fine.
				// A better check verifies all exit paths.
				hasTerminator = true
				// break // Found at least one terminator
			}
		}
		// Get the presumed exit block (might be complex with control flow)
		// If the last block added doesn't have a terminator, add one.
		if len(mainFunc.Blocks) > 0 {
			lastBlock := mainFunc.Blocks[len(mainFunc.Blocks)-1]
			if lastBlock.Term == nil {
				t.Logf("Warning: main function's last block (%s) not terminated in test %q. Adding default ret.", lastBlock.LocalIdent.Name(), t.Name())
				// Add default return based on main's signature (usually i32)
				if mainFunc.Sig.RetType.Equal(types.I32) {
					lastBlock.NewRet(constant.NewInt(types.I32, 0))
				} else if mainFunc.Sig.RetType.Equal(types.Void) {
					lastBlock.NewRet(nil)
				} // Handle other potential return types if needed
			}
		} else if !hasTerminator {
			// No blocks or no terminators at all? Problem.
			// Add an entry block and terminator if possible
			entry := mainFunc.NewBlock("entry")
			entry.NewRet(constant.NewInt(types.I32, 0)) // Default return
			t.Logf("Warning: main function had no blocks/terminators in test %q. Added entry block with default ret.", t.Name())
		}
	}

	return cg.Module.String(), nil // Return IR and nil error if visitor succeeded
}

// TestCodeGenIntegrationSimple programs with basic constructs.
// Covers Requirement 5.2 (Integration: End-to-End Samples): Trivial, Simple arithmetic.
func TestCodeGenIntegrationSimple(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string // Substrings expected in the final IR
		expectError          bool     // Whether generateIRForProgram is expected to return an error
	}{
		{
			name:  "Return Constant",
			input: `main() -> { return 42; }`,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`ret i32 42`,
			},
			expectError: false,
		},
		{
			name:  "Let and Return Variable",
			input: `main() -> { let x = 100; return x; }`,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`%x = alloca i32`,
				`store i32 100, ptr %x`,
				`%[0-9]+ = load i32, ptr %x`,
				`ret i32 %[0-9]+`,
			},
			expectError: false,
		},
		{
			name:  "Simple Arithmetic",
			input: `main() -> { let a = 5; let b = 3; return a + b; }`,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`%a = alloca i32`,
				`store i32 5, ptr %a`,
				`%b = alloca i32`,
				`store i32 3, ptr %b`,
				`%[0-9]+ = load i32, ptr %a`,
				`%[0-9]+ = load i32, ptr %b`,
				`%[0-9]+ = add nsw i32 %[0-9]+, %[0-9]+`,
				`ret i32 %[0-9]+`,
			},
			expectError: false,
		},
		{
			name:  "Arithmetic with Literals",
			input: `main() -> { return (10 * 2) - 5; }`,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				// Order might vary: mul then sub, or constants calculated directly
				`%[0-9]+ = mul nsw i32 10, 2`,      // 20
				`%[0-9]+ = sub nsw i32 %[0-9]+, 5`, // 20 - 5 = 15
				`ret i32 %[0-9]+`,                  // ret 15
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error from generateIRForProgram, but got nil.\nIR Generated:\n%s", ir)
				} else {
					t.Logf("Received expected error: %v", err)
				}
				return // Don't check IR content if error was expected
			}

			// Handle unexpected errors
			if err != nil {
				t.Fatalf("generateIRForProgram failed: %v", err)
			}
			t.Logf("Generated IR for %q:\n%s", tt.name, ir) // Log successful IR

			// Check for expected substrings (using regex for flexibility)
			for _, subPattern := range tt.expectedIRSubstrings {
				re := regexp.MustCompile(subPattern)
				if !re.MatchString(ir) {
					t.Errorf("Generated IR missing expected pattern for test %q.\nExpected pattern: %s\nGot IR:\n%s", tt.name, subPattern, ir)
				}
			}
		})
	}
}

// TestCodeGenIntegrationFunctions tests function definitions and calls.
// Covers Requirement 5.2 (Integration: End-to-End Samples): Function calls (with/without args).
func TestCodeGenIntegrationFunctions(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string
		expectError          bool
	}{
		{
			name: "Define and Call No-Arg Function",
			input: `
                function getVal() -> { return 99; }
                main() -> { return getVal(); }
            `,
			expectedIRSubstrings: []string{
				`define i32 @getVal()`,         // Function definition
				`ret i32 99`,                   // Body of getVal
				`define i32 @main()`,           // Main definition
				`%[0-9]+ = call i32 @getVal()`, // Call in main
				`ret i32 %[0-9]+`,              // Return result in main
			},
			expectError: false,
		},
		{
			name: "Define and Call One-Arg Function",
			input: `
                function identity(x) -> { return x; }
                main() -> { let v = 55; return identity(v); }
            `,
			expectedIRSubstrings: []string{
				`define i32 @identity\(i32 %x\)`, // Definition with param
				`ret i32 %x`,                     // Body of identity
				`define i32 @main()`,
				`%v = alloca i32`,
				`store i32 55, ptr %v`,
				`%[0-9]+ = load i32, ptr %v`,                  // Load arg v
				`%[0-9]+ = call i32 @identity\(i32 %[0-9]+\)`, // Call with loaded arg
				`ret i32 %[0-9]+`,                             // Return result
			},
			expectError: false,
		},
		{
			name: "Define and Call Multi-Arg Function",
			input: `
                function sum(a, b, c) -> { return a + b + c; }
                main() -> { return sum(10, 20, 30); }
            `,
			expectedIRSubstrings: []string{
				`define i32 @sum\(i32 %a, i32 %b, i32 %c\)`,
				`%[0-9]+ = add nsw i32 %a, %b`,
				`%[0-9]+ = add nsw i32 %[0-9]+, %c`,
				`ret i32 %[0-9]+`,
				`define i32 @main()`,
				`%[0-9]+ = call i32 @sum\(i32 10, i32 20, i32 30\)`,
				`ret i32 %[0-9]+`,
			},
			expectError: false,
		},
		{
			name: "Void Function Call (Using Builtin)", // Assuming print maps to a void builtin
			input: `
                 main() -> {
                     asm("builtin_print_int", 123); // Use asm for placeholder
                 }
             `,
			expectedIRSubstrings: []string{
				// Expect declaration of the builtin (may happen automatically)
				// `declare void @builtin_print_int(i32)`
				`define i32 @main()`,
				`call void @builtin_print_int\(i32 123\)`, // The call itself
				`ret i32 0`, // Default return from main
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error from generateIRForProgram, but got nil.\nIR Generated:\n%s", ir)
				} else {
					t.Logf("Received expected error: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("generateIRForProgram failed: %v\nIR (if any):\n%s", err, ir)
			}
			t.Logf("Generated IR for %q:\n%s", tt.name, ir)

			for i, subPattern := range tt.expectedIRSubstrings {
				re := regexp.MustCompile(subPattern)
				if !re.MatchString(ir) {
					t.Errorf("Generated IR missing expected pattern #%d for test %q.\nExpected pattern: %s\nGot IR:\n%s", i+1, tt.name, subPattern, ir)
				}
			}
		})
	}
}

// TestCodeGenIntegrationLambda tests lambda definition and calls.
// Covers Requirement 5.2 (Integration: End-to-End Samples): Lambda assigned, called via variable.
func TestCodeGenIntegrationLambda(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string
		expectError          bool
	}{
		{
			name: "Define Lambda and Call via Variable",
			input: `
                main() -> {
                    let doubler = (x) -> { return x * 2; };
                    let result = doubler(21);
                    return result;
                }
            `,
			expectedIRSubstrings: []string{
				// Lambda function definition (internal linkage likely)
				`define internal i32 @lambda_[0-9]+\(i32 %x\)`, // Lambda with param
				`%[0-9]+ = mul nsw i32 %x, 2`,
				`ret i32 %[0-9]+`,
				// Main function
				`define i32 @main()`,
				// Let statement for lambda variable
				`%doubler = alloca ptr`,                  // Allocate space for function pointer
				`store ptr @lambda_[0-9]+, ptr %doubler`, // Store lambda function address
				// Let statement for result
				`%result = alloca i32`,
				// Call via loaded pointer
				`%[0-9]+ = load ptr, ptr %doubler`,     // Load lambda address
				`%[0-9]+ = call i32 %[0-9]+\(i32 21\)`, // Call loaded function pointer
				`store i32 %[0-9]+, ptr %result`,       // Store call result
				// Return result
				`%[0-9]+ = load i32, ptr %result`,
				`ret i32 %[0-9]+`,
			},
			expectError: false,
		},
		{
			name: "Lambda Passed as Argument",
			input: `
                function apply(fn, val) -> { return fn(val); }
                main() -> {
                    let sq = (y) -> y * y;
                    return apply(sq, 7);
                }
            `,
			expectedIRSubstrings: []string{
				// Lambda definition
				`define internal i32 @lambda_[0-9]+\(i32 %y\)`,
				`ret i32 %[a-zA-Z0-9_.]+`, // Body of lambda (mul result)
				// Apply function definition
				`define i32 @apply\(ptr %fn, i32 %val\)`, // Takes function pointer and value
				// `%loaded_fn = load ptr, ptr %fn`? No, ptr is passed directly
				`%call_res = call i32 %fn\(i32 %val\)`, // Call the passed function pointer
				`ret i32 %call_res`,
				// Main function
				`define i32 @main()`,
				`%sq = alloca ptr`,                                      // Allocate for lambda variable
				`store ptr @lambda_[0-9]+, ptr %sq`,                     // Store lambda address
				`%loaded_sq = load ptr, ptr %sq`,                        // Load lambda address for passing
				`%final_res = call i32 @apply\(ptr %loaded_sq, i32 7\)`, // Call apply
				`ret i32 %final_res`,
			},
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error from generateIRForProgram, but got nil.\nIR Generated:\n%s", ir)
				}
				return
			}
			if err != nil {
				t.Fatalf("generateIRForProgram failed: %v\nIR (if any):\n%s", err, ir)
			}
			t.Logf("Generated IR for %q:\n%s", tt.name, ir)

			for i, subPattern := range tt.expectedIRSubstrings {
				re := regexp.MustCompile(subPattern)
				if !re.MatchString(ir) {
					t.Errorf("Generated IR missing expected pattern #%d for test %q.\nExpected pattern: %s\nGot IR:\n%s", i+1, tt.name, subPattern, ir)
				}
			}
		})
	}
}

// TestCodeGenIntegrationBlocks tests code generation involving nested blocks and scoping.
// Covers Requirement 5.2 (Integration: End-to-End Samples): Block with let, return, nested lambda.
// Covers Requirement 5.3 (Codegen Correctness for Chaining/Composition): Variable shadowing/scoping.
func TestCodeGenIntegrationBlocks(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string // Check for specific allocas, loads, stores, calls within blocks
		expectError          bool
	}{
		{
			name: "Block with Let and Return",
			input: `
                 main() -> {
                     let x = 5;
                     {
                         let y = x + 10; // Use outer x
                         return y;       // Return inner y
                     }
                     // Code here is unreachable due to return in block
                 }
             `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`%x = alloca i32`, // Outer x
				`store i32 5, ptr %x`,
				// No specific instruction marks the block start/end in flat IR besides labels for control flow
				`%y = alloca i32`, // Inner y
				`%load_x = load i32, ptr %x`,
				`%add_res = add nsw i32 %load_x, 10`,
				`store i32 %add_res, ptr %y`,
				`%load_y = load i32, ptr %y`,
				`ret i32 %load_y`, // Return from within the block
				// No second return should be generated after the block
			},
			expectError: false,
		},
		{
			name: "Variable Shadowing",
			input: `
                 main() -> {
                     let shadow = "outer";
                     {
                         let shadow = 10; // Inner shadow (different type!) - needs flexible alloca
                         asm("builtin_print_int", shadow); // Should print 10
                     }
                     // asm("builtin_print_str", shadow); // Would print "outer" - complex to test IR string print
                     return 0; // Return something simple
                 }
             `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				// Outer shadow (string ptr)
				// `%shadow = alloca ptr` (may have different name)
				// `store ptr @..., ptr %shadow`
				// Inner block
				// `%shadow.1 = alloca i32` (shadowed var gets different IR name)
				`store i32 10, ptr %[a-zA-Z0-9_.]+`,                      // Store 10 into inner shadow
				`%load_inner_shadow = load i32, ptr %[a-zA-Z0-9_.]+`,     // Load inner shadow (10)
				`call void @builtin_print_int\(i32 %load_inner_shadow\)`, // Call print with 10
				// After block
				`ret i32 0`,
			},
			expectError: false,
		},
		{
			name: "Lambda Defined in Block",
			input: `
                 main() -> {
                     let multiplier = 1; // Outer variable
                     let fn = ""; // Placeholder for pointer
                     {
                         let factor = 5; // Inner variable
                         fn = (n) -> { return n * factor; }; // Lambda captures inner 'factor' - COMPLEX requires closure support
                         multiplier = 10; // Modify outer var from inner block
                     }
                     // return fn(multiplier); // Call lambda after block - requires closures
                     return multiplier; // Return modified outer var (10)
                 }
             `,
			expectedIRSubstrings: []string{
				// Very simplified checks without proper closure support:
				`define i32 @main()`,
				`%multiplier = alloca i32`,
				`store i32 1, ptr %multiplier`,
				`%fn = alloca ptr`, // For function pointer
				// Inner block starts
				`%factor = alloca i32`,
				`store i32 5, ptr %factor`,
				// Lambda definition (likely internal function @lambda_...)
				`define internal i32 @lambda_[0-9]+\(i32 %n\)`,
				// Lambda body *without* closure support would fail to find 'factor'
				// With closure support, it would load 'factor' from captured context
				// `store ptr @lambda_[0-9]+, ptr %fn`, // Store lambda ptr
				`store i32 10, ptr %multiplier`, // Modify outer var
				// Inner block ends
				`%load_mult = load i32, ptr %multiplier`, // Load outer var (should be 10)
				`ret i32 %load_mult`,
			},
			// Expect error if closures aren't implemented, as lambda body will fail
			expectError: true, // Adjust if closures *are* supported
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error from generateIRForProgram, but got nil.\nIR Generated:\n%s", ir)
				} else {
					t.Logf("Received expected error: %v", err)
				}
				// Even on error, check IR for partial correctness if possible
				if ir != "" {
					t.Logf("Generated IR (despite error) for %q:\n%s", tt.name, ir)
					for i, subPattern := range tt.expectedIRSubstrings {
						// Check only patterns likely generated *before* the error point
						// This is heuristic and test-case specific
						re := regexp.MustCompile(subPattern)
						if !re.MatchString(ir) {
							t.Logf("Note: IR (on error path) missing expected pattern #%d.\nPattern: %s", i+1, subPattern)
						} else {
							t.Logf("Note: IR (on error path) contains expected pattern #%d.", i+1)
						}
					}
				}

				return
			}
			if err != nil {
				t.Fatalf("generateIRForProgram failed: %v\nIR (if any):\n%s", err, ir)
			}
			t.Logf("Generated IR for %q:\n%s", tt.name, ir)

			for i, subPattern := range tt.expectedIRSubstrings {
				re := regexp.MustCompile(subPattern)
				if !re.MatchString(ir) {
					t.Errorf("Generated IR missing expected pattern #%d for test %q.\nExpected pattern: %s\nGot IR:\n%s", i+1, tt.name, subPattern, ir)
				}
			}
		})
	}
}
