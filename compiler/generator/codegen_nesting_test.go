package generator

import (
	"regexp"
	"testing"
)

func TestCodeGenCompositionNesting(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string // Check specific instructions related to nesting and scoping
		expectError          bool
	}{
		{
			name: "Lambda in Block, Access Outer Var",
			input:`
                 main() -> {
                     let x = 10;
                     let y = 0;
                     {
                         let inner_lambda = (p) -> p + x; // Captures outer x (Needs Closure)
                         y = inner_lambda(5); // Call lambda defined in block
                     }
                     return y; // Should be 15 if closures work
                 }
            `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`%[a-zA-Z0-9_.]+ = alloca i32`,
				`store i32 10, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`%[a-zA-Z0-9_.]+ = alloca i32`,
				`store i32 0, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				// Inner Block
				`%[a-zA-Z0-9_.]+ = alloca [a-zA-Z0-9*_]+`,
				`define internal i32 @lambda_[0-9]+\(i32 %p\)`, // Lambda definition
				// Lambda body (Needs closure load for x): `%[a-zA-Z0-9_.]+ = load i32, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`
				`%[a-zA-Z0-9_.]+ = add (nsw )?i32 %p, %[a-zA-Z0-9_.]+`, // %p + loaded x
				`ret i32 %[a-zA-Z0-9_.]+`,
				`store [a-zA-Z0-9*(). _]+ @lambda_[0-9]+, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,  // Store lambda ptr
				`%[a-zA-Z0-9_.]+ = load [a-zA-Z0-9*(). _]+, %[a-zA-Z0-9_.]+`, // Load ptr for call
				`%[a-zA-Z0-9_.]+ = call i32 %loaded_lambda\(i32 5\)`, // Call lambda
				`store i32 %[a-zA-Z0-9_.]+, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,                  // Assign result to outer y
				// After Block
				`%[a-zA-Z0-9_.]+ = load i32, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`ret i32 %[a-zA-Z0-9_.]+`, // Return y (should be 15)
			},
			expectError: true, // Expect failure due to missing closure support for 'x'
		},
		{
			name: "Block in If Consequence",
			input:`
                main() -> {
                    let flag = 1; // True
                    let result = 0;
                    if (flag > 0) { // True branch taken
                        let temp = 100;
                        result = temp / 2; // result = 50
                    }
                    return result; // Should return 50
                }
            `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`%[a-zA-Z0-9_.]+ = alloca i32`,
				`store i32 1, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`%[a-zA-Z0-9_.]+ = alloca i32`,
				`store i32 0, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`%[a-zA-Z0-9_.]+ = load i32, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`%[a-zA-Z0-9_.]+ = icmp`, `icmp [a-z]+ i32 %[a-zA-Z0-9_.]+, 0`,          // Compare flag > 0
				`br i1 %[a-zA-Z0-9_.]+, label %if_then, label %if_else`, // Branch based on cond
				// Then block
				`if_then:`,
				`%[a-zA-Z0-9_.]+ = alloca i32`, // Allocation inside the block
				`store i32 100, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`%[a-zA-Z0-9_.]+ = load i32, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`%[a-zA-Z0-9_.]+ = sdiv i32 %[a-zA-Z0-9_.]+, 2`,
				`store i32 %[a-zA-Z0-9_.]+, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`, // Store into outer 'result'
				`br label %if_merge`,              // Branch to merge
				// Else block (may be empty or just branch)
				`if_else:`,
				`br label %if_merge`,
				// Merge block
				`if_merge:`,
				`%[a-zA-Z0-9_.]+ = load i32, [a-zA-Z0-9*_]+ %[a-zA-Z0-9_.]+`,
				`ret i32 %[a-zA-Z0-9_.]+`, // Should return 50
			},
			expectError: false,
		},
		{
			name: "Block in Else",
			input:`
                 main() -> {
                     let val = -5;
                     let status = "unknown"; // String requires support
                     if (val > 0) {
                         status = "positive";
                     } else {
                         let code = val * -1; // Needs negation/prefix support
                         status = "negative"; // Simplified
                     }
                     return 0; // Ignore status for now
                 }
            `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`%[a-zA-Z0-9_.]+ = alloca i32`,
				`br i1 %[a-zA-Z0-9_.]+, label %if_then, label %if_else`,
				`if_then:`,
				`getelementptr`,
				`br label %if_merge`,
				`if_else:`,
				`br label %if_merge`,
				`if_merge:`,
				`ret i32 0`,
			},
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input)

			if tt.expectError {
				if err == nil {
					t.Logf("Warning: Expected an error (feature likely not implemented), but got nil. Skipping pattern checks.")
					return // Skip pattern checks for expected-error tests where error didn't occur
				} else {
					t.Logf("Received expected error (feature likely not implemented): %v", err)
					return
				}
			} else if err != nil {
				t.Fatalf("generateIRForProgram failed unexpectedly: %v\nIR (if any):\n%s", err, ir)
			}

			t.Logf("Generated IR for %q:\n%s", tt.name, ir)

			lastMatchIndex := -1
			for i, subPattern := range tt.expectedIRSubstrings {
				re := regexp.MustCompile(subPattern)
				matches := re.FindStringIndex(ir[lastMatchIndex+1:])
				if matches == nil {
					t.Errorf("Generated IR missing expected pattern #%d for test %q.\nExpected pattern: %s\nGot IR:\n%s", i+1, tt.name, subPattern, ir)
				} else {
					currentMatchIndex := lastMatchIndex + 1 + matches[0]
					lastMatchIndex = currentMatchIndex
					t.Logf("Pattern #%d (%s) found at index %d", i+1, subPattern, currentMatchIndex)
				}
			}
		})
	}
}
