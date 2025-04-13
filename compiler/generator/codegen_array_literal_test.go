package generator

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/parser"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

// Helper to setup a basic Array struct type in the codegen context for tests
func setupArrayStructType(cg *CodeGenerator) error {
	if _, exists := cg.Structs["Array"]; exists {
		return nil // Already defined
	}
	// Define based on stdlib/array.y: { length: int, data: *int }
	// Assuming int -> i32, *int -> i32*
	llvmIntType, errInt := cg.mapType("int")
	llvmIntPtrType, errIntPtr := cg.mapType("*int")
	if errInt != nil || errIntPtr != nil {
		return fmt.Errorf("failed mapType for Array fields: %v, %v", errInt, errIntPtr)
	}
	if llvmIntType == nil || llvmIntPtrType == nil {
		return fmt.Errorf("mapType returned nil for Array fields")
	}

	arrayStructType := types.NewStruct(llvmIntType, llvmIntPtrType)
	cg.Module.NewTypeDef("Array", arrayStructType) // Define the type globally
	cg.Structs["Array"] = arrayStructType          // Store for lookup
	return nil
}

func TestCodeGenArrayLiteralUnit(t *testing.T) {
	tests := []struct {
		name                   string
		input                  string // The array literal expression
		expectedArrayTypeRe    string // Regex for the underlying data array type (e.g., `\[3 x i32]`)
		expectedStructAlloca   string // Expect allocation of the Array struct (e.g., `%array_struct = alloca %Array`)
		expectedDataAllocaRe   string // Expect allocation for the array data (e.g., `%array_data = alloca \[3 x i32]`)
		expectedLengthStoreRe  string // Expect store of length (e.g., `store i32 3, ptr %len_addr`)
		expectedDataPtrStoreRe string // Expect store of data pointer (e.g., `store ptr %first_elem, ptr %data_addr`)
		expectedConstValuesRe  string // Regex for the constant array initializer values (e.g., `\[i32 1, i32 2, i32 3]`)
	}{
		{
			name:                   "Integer Array Literal",
			input:                  `[1, 2, 3]`,
			expectedArrayTypeRe:    `\[3 x i32]`,
			expectedStructAlloca:   `%[a-zA-Z0-9_.]+ = alloca %Array`,
			expectedDataAllocaRe:   `%[a-zA-Z0-9_.]+ = alloca \[3 x i32]`,
			expectedLengthStoreRe:  `store i32 3, ptr %[a-zA-Z0-9_.]+`,
			expectedDataPtrStoreRe: `store ptr %[a-zA-Z0-9_.]+, ptr %[a-zA-Z0-9_.]+`,
			expectedConstValuesRe:  `\[i32 1, i32 2, i32 3]`,
		},
		{
			name:                   "Empty Array Literal",
			input:                  `[]`,
			expectedArrayTypeRe:    "", // No underlying data array allocated in this specific way
			expectedStructAlloca:   `%[a-zA-Z0-9_.]+ = alloca %Array`,
			expectedDataAllocaRe:   "",                                    // No specific data allocation like above
			expectedLengthStoreRe:  `store i32 0, ptr %[a-zA-Z0-9_.]+`,    // Store length 0
			expectedDataPtrStoreRe: `store ptr null, ptr %[a-zA-Z0-9_.]+`, // Store null data pointer
			expectedConstValuesRe:  "",                                    // No constant values
		},
		// Add test for array of strings? Requires String literal codegen to be stable.
		// Add test for array of expressions? Requires expression codegen.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := parseExpr(t, tt.input) // Parse the array literal expression
			arrayLit, ok := node.(*ast.ArrayLiteral)
			if !ok {
				t.Fatalf("Parsed node is not ArrayLiteral, got %T", node)
			}

			// Setup codegen context
			cg := NewCodeGenerator()
			// Setup the %Array type definition needed by the generator
			if err := setupArrayStructType(cg); err != nil {
				t.Fatalf("Failed to setup Array struct type: %v", err)
			}

			// Create dummy main context
			retType := types.NewPointer(cg.Structs["Array"]) // Array literal evaluates to ptr to struct
			mainSig := types.NewFunc(retType)
			mainFunc := cg.Module.NewFunc("main", mainSig.RetType)
			entryBlock := mainFunc.NewBlock("entry")
			cg.currentFunc = mainFunc
			cg.Block = entryBlock

			// Visit the array literal node
			errVisit := arrayLit.Accept(cg)
			if errVisit != nil {
				// Check if it's the expected error for non-constant elements if applicable
				// if strings.Contains(errVisit.Error(), "non-constant element found") && tt.name == "Array of Expressions" {
				//     // Expected error for this case in current simple implementation
				//     return
				// }
				t.Fatalf("Codegen failed during array literal visit for input %q: %v", tt.input, errVisit)
			}

			// Add terminator (return the pointer to the array struct)
			if cg.Block != nil && cg.Block.Term == nil {
				if cg.lastValue != nil && cg.lastValue.Type().Equal(retType) {
					cg.Block.NewRet(cg.lastValue)
				} else {
					t.Logf("Warning: cg.lastValue (%v) type (%s) mismatch or nil after visiting array lit %q. Returning default null.", cg.lastValue, cg.lastValue.Type(), tt.input)
					cg.Block.NewRet(constant.NewNull(retType)) // Return null pointer as fallback
				}
			} else if cg.Block == nil {
				t.Fatalf("Codegen block became nil unexpectedly for input %q", tt.input)
			}

			ir := cg.Module.String()
			t.Logf("Generated IR for %q:\n%s", tt.input, ir) // Log IR for debugging

			// Check for struct allocation
			reStructAlloca := regexp.MustCompile(tt.expectedStructAlloca)
			if !reStructAlloca.MatchString(ir) {
				t.Errorf("Generated IR missing expected Array struct allocation.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedStructAlloca, ir)
			}

			// Checks specific to non-empty arrays
			if tt.input != "[]" {
				reDataAlloca := regexp.MustCompile(tt.expectedDataAllocaRe)
				if !reDataAlloca.MatchString(ir) {
					t.Errorf("Generated IR missing expected array data allocation.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedDataAllocaRe, ir)
				}
				// Check for constant array definition if values were constant
				if tt.expectedConstValuesRe != "" {
					reConst := regexp.MustCompile(`constant ` + regexp.QuoteMeta(tt.expectedArrayTypeRe) + ` ` + regexp.QuoteMeta(tt.expectedConstValuesRe))
					// This const might be part of the 'store' or a global, check broadly
					if !reConst.MatchString(ir) && !strings.Contains(ir, tt.expectedConstValuesRe) {
						t.Errorf("Generated IR missing expected constant array values.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedConstValuesRe, ir)
					}
				}
			}

			// Check length store
			reLenStore := regexp.MustCompile(tt.expectedLengthStoreRe)
			if !reLenStore.MatchString(ir) {
				t.Errorf("Generated IR missing expected length store.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedLengthStoreRe, ir)
			}

			// Check data pointer store
			reDataPtrStore := regexp.MustCompile(tt.expectedDataPtrStoreRe)
			if !reDataPtrStore.MatchString(ir) {
				t.Errorf("Generated IR missing expected data pointer store.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedDataPtrStoreRe, ir)
			}
		})
	}
}

func TestCodeGenIndexExpressionUnit(t *testing.T) {
	tests := []struct {
		name             string
		input            string // The expression using indexing
		setupInput       string // Setup code (e.g., defining the array)
		isLoad           bool   // True if testing load (RHS), false for store (LHS)
		expectedGEPRe    string // Regex for the GetElementPtr instruction
		expectedLoadRe   string // Regex for the load instruction (if isLoad is true)
		expectedStoreRe  string // Regex for the store instruction (if isLoad is false)
		expectedReturnRe string // Regex for the final return (loaded value or stored value)
	}{
		{
			name:       "Index Load (Integer Index)",
			input:      `myArr[1]`,
			setupInput: `let myArr = [10, 20, 30];`,
			isLoad:     true,
			// GEP from Array struct ptr -> data field ptr -> element ptr
			// 1. GEP to get data field (%Array*, idx 0, idx 1 -> ptr*)
			// 2. Load data field (ptr* -> ptr)
			// 3. GEP from data ptr (ptr, index -> ptr)
			expectedGEPRe:    `getelementptr i32, ptr %[a-zA-Z0-9_.]+, i64 1`, // GEP from data array pointer with index 1
			expectedLoadRe:   `load i32, ptr %[a-zA-Z0-9_.]+`,                 // Load from the element address
			expectedStoreRe:  "",
			expectedReturnRe: `ret i32 %[a-zA-Z0-9_.]+`, // Return the loaded i32 value
		},
		{
			name:       "Index Load (Variable Index)",
			input:      `myArr[idx]`,
			setupInput: `let myArr = [10, 20, 30]; let idx = 2;`,
			isLoad:     true,
			// Expect load of idx before GEP
			expectedGEPRe:    `getelementptr i32, ptr %[a-zA-Z0-9_.]+, i64 %[a-zA-Z0-9_.]+`, // GEP uses loaded index
			expectedLoadRe:   `load i32, ptr %[a-zA-Z0-9_.]+`,
			expectedStoreRe:  "",
			expectedReturnRe: `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:             "Index Store (Integer Index)",
			input:            `myArr[0] = 99`,
			setupInput:       `let myArr = [10, 20, 30];`,
			isLoad:           false,
			expectedGEPRe:    `getelementptr i32, ptr %[a-zA-Z0-9_.]+, i64 0`, // GEP to get address of element 0
			expectedLoadRe:   "",
			expectedStoreRe:  `store i32 99, ptr %[a-zA-Z0-9_.]+`, // Store 99 into the element address
			expectedReturnRe: `ret i32 99`,                        // Assignment expression returns the RHS value
		},
		{
			name:       "Index Store (Variable Index, Variable Value)",
			input:      `myArr[idx] = newValue`,
			setupInput: `let myArr = [10, 20, 30]; let idx = 1; let newValue = 55;`,
			isLoad:     false,
			// Expect loads for idx and newValue
			expectedGEPRe:    `getelementptr i32, ptr %[a-zA-Z0-9_.]+, i64 %[a-zA-Z0-9_.]+`,
			expectedLoadRe:   "",
			expectedStoreRe:  `store i32 %[a-zA-Z0-9_.]+, ptr %[a-zA-Z0-9_.]+`, // Store loaded newValue into GEP address
			expectedReturnRe: `ret i32 %[a-zA-Z0-9_.]+`,                        // Return the loaded newValue
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the whole block to handle setup and use
			fullInput := fmt.Sprintf("main() -> { %s %s; }", tt.setupInput, tt.input) // Add semicolon for use expr
			l, err := lexer.NewLexerFromString(fullInput)
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}
			p := parser.NewParser(l)
			prog := p.ParseProgram()
			if len(p.Errors()) > 0 {
				t.Fatalf("Parser errors: %v", p.Errors())
			}
			if prog.MainFunction == nil || prog.MainFunction.Body == nil {
				t.Fatalf("Failed to parse main function")
			}
			body, ok := prog.MainFunction.Body.(*ast.BlockStatement)
			if !ok || len(body.Statements) == 0 {
				t.Fatalf("Main body is not block or empty")
			}

			// Setup codegen context
			cg := NewCodeGenerator()
			if err := setupArrayStructType(cg); err != nil { // Ensure %Array is defined
				t.Fatalf("Failed to setup Array struct type: %v", err)
			}
			retType := types.I32 // Assume result/return is i32 for simplicity
			mainSig := types.NewFunc(retType)
			mainFunc := cg.Module.NewFunc("main", mainSig.RetType)
			entryBlock := mainFunc.NewBlock("entry")
			cg.currentFunc = mainFunc
			cg.Block = entryBlock

			// Visit the entire function body (setup + use)
			errVisit := body.Accept(cg)
			if errVisit != nil {
				t.Errorf("runCodeGen failed: %v", errVisit)
			}

			// Ensure termination
			if cg.Block != nil && cg.Block.Term == nil {
				if cg.lastValue != nil && cg.lastValue.Type().Equal(retType) {
					cg.Block.NewRet(cg.lastValue)
				} else {
					cg.Block.NewRet(constant.NewInt(types.I32, 0))
				}
			}

			ir := cg.Module.String()
			t.Logf("Generated IR for %q:\n%s", tt.input, ir)

			// Check for GetElementPtr
			reGEP := regexp.MustCompile(tt.expectedGEPRe)
			if !reGEP.MatchString(ir) {
				t.Errorf("Generated IR missing expected GEP instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedGEPRe, ir)
			}

			// Check for Load or Store
			if tt.isLoad {
				reLoad := regexp.MustCompile(tt.expectedLoadRe)
				if !reLoad.MatchString(ir) {
					t.Errorf("Generated IR missing expected Load instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedLoadRe, ir)
				}
			} else {
				reStore := regexp.MustCompile(tt.expectedStoreRe)
				if !reStore.MatchString(ir) {
					t.Errorf("Generated IR missing expected Store instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedStoreRe, ir)
				}
			}

			// Check for Return
			reRet := regexp.MustCompile(tt.expectedReturnRe)
			if !reRet.MatchString(ir) {
				t.Errorf("Generated IR missing expected Return instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedReturnRe, ir)
			}
		})
	}
}
