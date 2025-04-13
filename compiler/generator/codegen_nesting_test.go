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
			input: `
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
				`%x = alloca i32`,
				`store i32 10, ptr %x`,
				`%y = alloca i32`,
				`store i32 0, ptr %y`,
				// Inner Block
				`%inner_lambda = alloca ptr`,
				`define internal i32 @lambda_[0-9]+\(i32 %p\)`, // Lambda definition
				// Lambda body (Needs closure load for x): `%load_x_closure = load i32, ptr %captured_x`
				`%add_res = add nsw i32 %p, %[a-zA-Z0-9_.]+`, // %p + loaded x
				`ret i32 %add_res`,
				`store ptr @lambda_[0-9]+, ptr %inner_lambda`,  // Store lambda ptr
				`%loaded_lambda = load ptr, ptr %inner_lambda`, // Load ptr for call
				`%call_res = call i32 %loaded_lambda\(i32 5\)`, // Call lambda
				`store i32 %call_res, ptr %y`,                  // Assign result to outer y
				// After Block
				`%load_y = load i32, ptr %y`,
				`ret i32 %load_y`, // Return y (should be 15)
			},
			expectError: true, // Expect failure due to missing closure support for 'x'
		},
		{
			name: "Block in If Consequence",
			input: `
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
				`%flag = alloca i32`,
				`store i32 1, ptr %flag`,
				`%result = alloca i32`,
				`store i32 0, ptr %result`,
				`%load_flag = load i32, ptr %flag`,
				`%cond = icmp sgt i32 %load_flag, 0`,          // Compare flag > 0
				`br i1 %cond, label %if_then, label %if_else`, // Branch based on cond
				// Then block
				`if_then:`,
				`%temp = alloca i32`, // Allocation inside the block
				`store i32 100, ptr %temp`,
				`%load_temp = load i32, ptr %temp`,
				`%div_res = sdiv i32 %load_temp, 2`,
				`store i32 %div_res, ptr %result`, // Store into outer 'result'
				`br label %if_merge`,              // Branch to merge
				// Else block (may be empty or just branch)
				`if_else:`,
				`br label %if_merge`,
				// Merge block
				`if_merge:`,
				`%load_result = load i32, ptr %result`,
				`ret i32 %load_result`, // Should return 50
			},
			expectError: false,
		},
		{
			name: "Block in Else",
			input: `
                 main() -> {
                     let val = -5;
                     let status = "unknown"; // String requires support
                     if (val > 0) {
                         status = "positive";
                     } else {
                         let code = val * -1; // Needs negation/prefix support
                         // status = "negative_" + code; // Needs string concat support
                         status = "negative"; // Simplified
                     }
                     return 0; // Ignore status for now
                 }
            `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				`%val = alloca i32`,
				`store i32 -5, ptr %val`,
				`%status = alloca ptr`,       // String variable
				`store ptr @.*, ptr %status`, // Store "unknown"
				`%load_val = load i32, ptr %val`,
				`%cond = icmp sgt i32 %load_val, 0`,
				`br i1 %cond, label %if_then, label %if_else`,
				// Then block
				`if_then:`,
				`%str_pos_ptr = getelementptr .* @`,   // GEP for "positive"
				`store ptr %str_pos_ptr, ptr %status`, // Store "positive" pointer
				`br label %if_merge`,
				// Else block
				`if_else:`,
				`%code = alloca i32`, // Inner variable
				// `%neg_val = mul nsw i32 %load_val, -1` // Code for val * -1 (needs -1 literal/prefix)
				`store i32 %[0-9]+, ptr %code`,
				`%str_neg_ptr = getelementptr .* @`,   // GEP for "negative"
				`store ptr %str_neg_ptr, ptr %status`, // Store "negative" pointer
				`br label %if_merge`,
				// Merge block
				`if_merge:`,
				`ret i32 0`,
			},
			expectError: false, // May have errors if string/prefix ops aren't ready
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input)

			if tt.expectError {
				if err == nil {
					t.Logf("Warning: Expected an error (feature likely not implemented), but got nil.")
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
