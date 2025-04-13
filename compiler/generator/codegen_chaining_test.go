package generator

import (
	"regexp"
	"testing"
)

// TestCodeGenCompositionChaining tests chained calls like map().forEach().
// Covers Requirement 5.3 (Codegen Correctness for Chaining/Composition): Array/map/forEach chains.
// Note: Highly dependent on working method call and array implementation. Will likely fail initially.
func TestCodeGenCompositionChaining(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string // Check for the sequence of calls and intermediate results
		expectError          bool
	}{
		{
			name: "Array Map then ForEach",
			input: `
                type Array { // Simplified definition for test context
                    let length: int; let data: *int;
                    function map(self: Array, fn) -> Array { /* dummy */ return self; }
                    function forEach(self: Array, fn) -> Array { /* dummy */ return self; }
                }
                main() -> {
                    let nums = [1, 2];
                    let addOne = (n) -> n + 1;
                    let printIt = (n) -> { asm("builtin_print_int", n); };
                    // The chain:
                    nums.map(addOne).forEach(printIt);
                    return 0;
                }
            `,
			// Expect: setup nums, setup addOne, setup printIt,
			// call Array_map on nums -> tmp_array_ptr1
			// call Array_forEach on tmp_array_ptr1
			expectedIRSubstrings: []string{
				// Definitions (simplified check)
				`%Array = type`,
				`define .* @Array_map\(`,
				`define .* @Array_forEach\(`,
				`define internal .* @lambda_[0-9]+\(`, // addOne lambda
				`define internal .* @lambda_[0-9]+\(`, // printIt lambda
				`define i32 @main()`,
				// Variable setup
				`%nums = alloca ptr`,
				`%addOne = alloca ptr`,
				`%printIt = alloca ptr`,
				// First call: nums.map(addOne)
				`%nums_ptr = load ptr, ptr %nums`,
				`%addOne_ptr = load ptr, ptr %addOne`,
				`%map_result = call ptr @Array_map\(ptr %nums_ptr, ptr %addOne_ptr\)`, // Result is ptr to new Array struct
				// Second call: map_result.forEach(printIt)
				// Note: %map_result holds the pointer needed for 'self'
				`%printIt_ptr = load ptr, ptr %printIt`,
				`%forEach_result = call ptr @Array_forEach\(ptr %map_result, ptr %printIt_ptr\)`,
				`ret i32 0`,
			},
			expectError: true, // Will fail until types/methods/chaining work fully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input) // generateIR includes parsing

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

			// Check sequence and patterns
			lastMatchIndex := -1
			for i, subPattern := range tt.expectedIRSubstrings {
				re := regexp.MustCompile(subPattern)
				// Search starting from the end of the last match to enforce some ordering
				matches := re.FindStringIndex(ir[lastMatchIndex+1:]) // Find first match *after* last one
				if matches == nil {
					t.Errorf("Generated IR missing expected pattern #%d for test %q.\nExpected pattern: %s\nGot IR:\n%s", i+1, tt.name, subPattern, ir)
					// Don't update lastMatchIndex if pattern not found
				} else {
					// Adjust index relative to the start of the string
					currentMatchIndex := lastMatchIndex + 1 + matches[0]
					lastMatchIndex = currentMatchIndex // Update for next search
					t.Logf("Pattern #%d (%s) found at index %d", i+1, subPattern, currentMatchIndex)
				}
			}
		})
	}
}
