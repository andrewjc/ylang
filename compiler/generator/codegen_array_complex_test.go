package generator

import (
	"regexp"
	"testing"
)

func TestCodeGenIntegrationArrays(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string
		expectError          bool
	}{
		{
			name: "Allocate and Return Empty Array",
			input: `
                main() -> Array { // Assume parser allows type hint for return
                     let arr = [];
                     return arr;
                }
            `,
			expectedIRSubstrings: []string{
				`define ptr @main()`,                      // Returns ptr to %Array struct
				`%arr = alloca ptr`,                       // Allocate space for the 'arr' variable (holding ptr to %Array)
				`%empty_array_struct = alloca %Array`,     // Allocate the struct itself
				`store i32 0, ptr %`,                      // Store length 0
				`store ptr null, ptr %`,                   // Store data null
				`store ptr %empty_array_struct, ptr %arr`, // Store struct addr into 'arr' variable
				`%retval = load ptr, ptr %arr`,            // Load struct addr for return
				`ret ptr %retval`,
			},
			expectError: false,
		},
		{
			name: "Allocate Int Array, Index, Return Element",
			input: `
                main() -> {
                    let vals = [11, 22, 33];
                    let idx = 1;
                    return vals[idx];
                }
            `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				// Array Literal Generation (simplified checks)
				`%vals = alloca ptr`, // Var holds ptr to %Array struct
				`%array_struct = alloca %Array`,
				`%array_data = alloca \[3 x i32]`,
				`store i32 3, ptr %`,                 // Store len
				`store ptr %`,                        // Store data ptr
				`store ptr %array_struct, ptr %vals`, // Store struct addr in 'vals' var
				// Index Variable
				`%idx = alloca i32`,
				`store i32 1, ptr %idx`,
				// Indexing Operation
				`%loaded_struct_ptr = load ptr, ptr %vals`, // Load %Array* from 'vals'
				// `%data_addr_ptr = getelementptr %Array, ptr %loaded_struct_ptr, i32 0, i32 1`, // GEP to data field addr
				// `%data_ptr = load ptr, ptr %data_addr_ptr`, // Load the i32* data pointer
				`%loaded_idx = load i32, ptr %idx`,       // Load index value (1)
				`%idx_i64 = sext i32 %loaded_idx to i64`, // Extend index to i64 for GEP
				// `%element_addr = getelementptr i32, ptr %data_ptr, i64 %idx_i64`, // GEP to element address
				`%element_val = load i32, ptr %[a-zA-Z0-9_.]+`, // Load the value (22) from element address
				`ret i32 %element_val`,
			},
			expectError: false,
		},
		{
			name: "Allocate Int Array, Index Assign, Return Assigned",
			input: `
                 main() -> {
                     let data = [0, 0, 0];
                     data[1] = 77;
                     return data[1];
                 }
             `,
			expectedIRSubstrings: []string{
				`define i32 @main()`,
				// Array Literal Generation for data
				`%data = alloca ptr`,
				// Indexing and Assignment
				`%assign_lhs_gep = getelementptr i32, ptr %[a-zA-Z0-9_.]+, i64 1`, // GEP to element 1 address
				`store i32 77, ptr %assign_lhs_gep`,                               // Store 77 into element 1
				// Indexing and Load for Return
				`%ret_gep = getelementptr i32, ptr %[a-zA-Z0-9_.]+, i64 1`, // GEP again (or reuse)
				`%ret_val = load i32, ptr %ret_gep`,                        // Load element 1 (should be 77)
				`ret i32 %ret_val`,
			},
			expectError: false,
		},
		// Add tests for method calls like map/forEach when codegen supports them
		// {
		//     name: "Array Map Call (Conceptual)",
		//     input: `
		//         main() -> {
		//             let nums = [1, 2, 3];
		//             let doubler = (n) -> n * 2;
		//             let doubled = nums.map(doubler);
		//             return doubled[0]; // Expect 2
		//         }
		//     `,
		//     expectedIRSubstrings: []string{
		//         // ... array setup ...
		//         // ... lambda setup ...
		//         `call %Array\* @Array_map\(ptr %nums_struct_ptr, ptr %doubler_lambda_ptr\)`, // Conceptual call
		//         // ... indexing into 'doubled' result ...
		//         `ret i32 2`, // Simplified expected result
		//     },
		//      expectError: false, // Will fail until method calls are implemented
		// },
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
				// Optionally check partial IR on error if needed
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

func TestCodeGenIntegrationMethodCalls(t *testing.T) {
	tests := []struct {
		name                 string
		input                string // Assume Array type is defined similar to stdlib/array.y
		expectedIRSubstrings []string
		expectError          bool // Expected to fail until methods are fully implemented
	}{
		{
			name: "Array Map Call",
			input: `
                 type Array { // Simplified definition for test context
                     let length: int;
                     let data: *int;
                     function map(self: Array, fn) -> Array {
                         // Dummy body for now, focus is on call site
                         return self;
                     }
                 }
                 main() -> {
                     let nums = [1, 2];
                     let doubler = (n) -> n * 2;
                     let result = nums.map(doubler); // The method call
                     return 0; // Ignore result for now
                 }
             `,
			expectedIRSubstrings: []string{
				// Array struct definition
				`%Array = type { i32, ptr }`,
				// map method definition
				`define ptr @Array_map\(ptr %self, ptr %fn\)`, // Mangled name, self is first arg (ptr to %Array)
				// main function
				`define i32 @main()`,
				// nums array literal setup
				`%nums = alloca ptr`,
				`store ptr %`, // store struct ptr in nums var
				// doubler lambda setup
				`%doubler = alloca ptr`,
				`store ptr @lambda`,
				// result variable setup
				`%result = alloca ptr`, // Will hold ptr to the result %Array struct
				// Prepare args for method call
				`%nums_ptr = load ptr, ptr %nums`,      // Load ptr to nums struct (*%Array)
				`%lambda_ptr = load ptr, ptr %doubler`, // Load ptr to lambda (*func)
				// The actual method call
				`%call_map = call ptr @Array_map\(ptr %nums_ptr, ptr %lambda_ptr\)`,
				`store ptr %call_map, ptr %result`, // Store the returned struct ptr
				`ret i32 0`,
			},
			expectError: true, // Will fail until type/method codegen is complete
		},
		{
			name: "Array ForEach Call",
			input: `
                 type Array {
                     let length: int; let data: *int;
                     function forEach(self: Array, fn) -> Array { return self; }
                 }
                  main() -> {
                      let items = ["a", "b"];
                      let printFn = (s) -> { asm("builtin_print", s); }; // Placeholder print
                      items.forEach(printFn);
                      return 0;
                  }
             `,
			expectedIRSubstrings: []string{
				`%Array = type { i32, ptr }`,
				`define ptr @Array_forEach\(ptr %self, ptr %fn\)`,
				`define i32 @main()`,
				// items array setup (array of i8*)
				`%items = alloca ptr`,
				`store ptr %`,
				// printFn lambda setup
				`%printFn = alloca ptr`,
				`store ptr @lambda`,
				// Prepare args
				`%items_ptr = load ptr, ptr %items`,
				`%lambda_ptr = load ptr, ptr %printFn`,
				// The call (returns void or self ptr, assume self ptr for now)
				`%call_forEach = call ptr @Array_forEach\(ptr %items_ptr, ptr %lambda_ptr\)`,
				`ret i32 0`,
			},
			expectError: true, // Will fail until type/method codegen is complete
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input)

			if tt.expectError {
				// If it *doesn't* error, that might mean the codegen stubs worked unexpectedly well, or the test is wrong.
				if err == nil {
					t.Logf("Warning: Expected an error (feature likely not implemented), but got nil. Checking IR anyway.")
					// Proceed to check IR patterns that *should* exist if it worked
				} else {
					t.Logf("Received expected error (feature likely not implemented): %v", err)
					return // Stop here if error is expected and received
				}
			} else if err != nil {
				t.Fatalf("generateIRForProgram failed unexpectedly: %v\nIR (if any):\n%s", err, ir)
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

func TestCodeGenIntegrationStdlib(t *testing.T) {
	// Mark as skipped if stdlib isn't ready or configured
	// t.Skip("Skipping stdlib integration test until module loading is stable.")

	tests := []struct {
		name                 string
		input                string
		expectedIRSubstrings []string // Check for imported function calls etc.
		expectError          bool
	}{
		{
			name: "Import Core and Use Print (Conceptual)",
			input: `
                import "stdlib/core/print"; // Assume this defines print(any)
                main() -> {
                    print(123);
                }
            `,
			expectedIRSubstrings: []string{
				// Expect print function definition/declaration from the imported module
				// `define void @print(...)` or `declare void @print(...)`
				// Expect call to print in main
				`call void @print\(.*\)`, // Very general call pattern
			},
			expectError: true, // Will fail until import and print are fully functional
		},
		{
			name: "Import Array and Use Map (Conceptual)",
			input: `
                import "stdlib/array"; // Assume this defines Array type and methods
                 main() -> {
                     let a = [1, 2];
                     let f = (x) -> x + 1;
                     let b = a.map(f); // Method call
                     return 0;
                 }
             `,
			expectedIRSubstrings: []string{
				// Expect Array type definition from import
				`%Array = type { i32, ptr }`,
				// Expect Array_map function definition/declaration
				`define ptr @Array_map\(ptr %self, ptr %fn\)`,
				// Expect call to Array_map in main
				`call ptr @Array_map\(ptr %.*, ptr %.*\)`,
			},
			expectError: true, // Will fail until import and array methods work
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir, err := generateIRForProgram(t, tt.input) // generateIR includes parsing

			if tt.expectError {
				if err == nil {
					t.Logf("Warning: Expected an error (stdlib/import likely not implemented), but got nil.")
				} else {
					t.Logf("Received expected error (stdlib/import likely not implemented): %v", err)
					return
				}
			} else if err != nil {
				t.Fatalf("generateIRForProgram failed unexpectedly: %v\nIR (if any):\n%s", err, ir)
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
