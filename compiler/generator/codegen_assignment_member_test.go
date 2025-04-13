package generator

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/parser"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"regexp"
	"testing"
)

func TestCodeGenAssignmentExpressionUnit(t *testing.T) {
	tests := []struct {
		name             string
		input            string // The assignment expression statement
		setupInput       string // Setup code (e.g., defining variables)
		expectedStoreRe  string // Regex for the store instruction (store value, ptr address)
		expectedReturnRe string // Regex for the final return (should return the assigned value)
		expectedLoadRe   string // Regex for optional load if RHS is variable
	}{
		{
			name:             "Assign Literal to Variable",
			input:            `myVar = 42`,
			setupInput:       `let myVar = 0;`, // Variable must exist
			expectedLoadRe:   "",               // No load needed for literal RHS
			expectedStoreRe:  `store i32 42, ptr %[a-zA-Z0-9_.]+`,
			expectedReturnRe: `ret i32 42`, // Assignment returns the RHS
		},
		{
			name:             "Assign Variable to Variable",
			input:            `targetVar = sourceVar`,
			setupInput:       `let targetVar = 0; let sourceVar = 99;`,
			expectedLoadRe:   `load i32, ptr %[a-zA-Z0-9_.]+`, // Load sourceVar
			expectedStoreRe:  `store i32 %[a-zA-Z0-9_.]+, ptr %[a-zA-Z0-9_.]+`,
			expectedReturnRe: `ret i32 %[a-zA-Z0-9_.]+`, // Return the value loaded from sourceVar
		},
		{
			name:           "Assign Expression to Variable",
			input:          `result = x + y`,
			setupInput:     `let result = 0; let x = 5; let y = 6;`,
			expectedLoadRe: `load i32, ptr %[a-zA-Z0-9_.]+`, // Loads for x and y
			// Expect add instruction before store
			expectedStoreRe:  `store i32 %[a-zA-Z0-9_.]+, ptr %[a-zA-Z0-9_.]+`, // Store result of add
			expectedReturnRe: `ret i32 %[a-zA-Z0-9_.]+`,                        // Return result of add
		},
		// Add tests for assigning to index expressions `arr[i] = val` (Covered partly by Index tests)
		// Add tests for assigning to member access `obj.field = val`
		{
			name:             "Assign to Member Access",
			input:            `myArray.length = 10`,               // Assuming Array struct { length: i32, data: ptr }
			setupInput:       `let myArray = [1,2];`,              // Creates an Array struct instance
			expectedLoadRe:   "",                                  // No load for literal RHS
			expectedStoreRe:  `store i32 10, ptr %[a-zA-Z0-9_.]+`, // Store into the GEP address of length field
			expectedReturnRe: `ret i32 10`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the whole block
			fullInput := fmt.Sprintf("main() -> { %s %s; }", tt.setupInput, tt.input)
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
			if err := setupArrayStructType(cg); err != nil { // Needed for member access test
				t.Fatalf("Failed to setup Array struct type: %v", err)
			}
			retType := types.I32 // Assignment expr returns RHS, assume i32 for test
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

			// Check for Load (if RHS is variable)
			if tt.expectedLoadRe != "" {
				reLoad := regexp.MustCompile(tt.expectedLoadRe)
				if !reLoad.MatchString(ir) {
					t.Errorf("Generated IR missing expected Load instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedLoadRe, ir)
				}
			}

			// Check for Store instruction
			reStore := regexp.MustCompile(tt.expectedStoreRe)
			if !reStore.MatchString(ir) {
				t.Errorf("Generated IR missing expected Store instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedStoreRe, ir)
			}

			// Check for Return instruction
			reRet := regexp.MustCompile(tt.expectedReturnRe)
			if !reRet.MatchString(ir) {
				t.Errorf("Generated IR missing expected Return instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedReturnRe, ir)
			}
		})
	}
}

func TestCodeGenMemberAccessExpressionUnit(t *testing.T) {
	tests := []struct {
		name             string
		input            string // The member access expression
		setupInput       string // Setup code
		expectedGEPRe    string // Regex for the GEP instruction to get member address
		expectedLoadRe   string // Regex for the load instruction from member address
		expectedReturnRe string // Regex for returning the loaded value
	}{
		{
			name:             "Access Struct Member (Length)",
			input:            `myArray.length`,
			setupInput:       `let myArray = [100, 200];`,                                        // Creates %Array{i32 2, ptr ...}
			expectedGEPRe:    `getelementptr inbounds %Array, ptr %[a-zA-Z0-9_.]+, i32 0, i32 0`, // GEP for length field (index 0)
			expectedLoadRe:   `load i32, ptr %[a-zA-Z0-9_.]+`,                                    // Load i32 from the GEP address
			expectedReturnRe: `ret i32 %[a-zA-Z0-9_.]+`,                                          // Return the loaded i32
		},
		{
			name:             "Access Struct Member (Data Ptr)",
			input:            `myArray.data`,
			setupInput:       `let myArray = [100];`,                                             // Creates %Array{i32 1, ptr ...}
			expectedGEPRe:    `getelementptr inbounds %Array, ptr %[a-zA-Z0-9_.]+, i32 0, i32 1`, // GEP for data field (index 1)
			expectedLoadRe:   `load ptr, ptr %[a-zA-Z0-9_.]+`,                                    // Load ptr (i32*) from the GEP address
			expectedReturnRe: `ret ptr %[a-zA-Z0-9_.]+`,                                          // Return the loaded ptr
		},
		// Add tests for nested access like `obj.inner.field` when supported
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the whole block
			fullInput := fmt.Sprintf("main() -> { %s return %s; }", tt.setupInput, tt.input) // Use expression in return
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
			// Target node is the member access expression inside the return statement
			targetStmt := body.Statements[len(body.Statements)-1] // Assume return is last
			var memberExpr ast.ExpressionNode                     // Can be MemberAccessExpression or DotOperator
			if retStmt, okRet := targetStmt.(*ast.ReturnStatement); okRet {
				memberExpr = retStmt.ReturnValue
				if _, okMae := memberExpr.(*ast.MemberAccessExpression); !okMae {
					if _, okDot := memberExpr.(*ast.DotOperator); !okDot {
						t.Fatalf("Return value is not MemberAccessExpression or DotOperator, got %T", memberExpr)
					}
				}
			} else {
				t.Fatalf("Target statement is not ReturnStatement, got %T", targetStmt)
			}

			// Setup codegen context
			cg := NewCodeGenerator()
			if err := setupArrayStructType(cg); err != nil {
				t.Fatalf("Failed to setup Array struct type: %v", err)
			}

			// Determine expected return type based on test case
			var retType types.Type = types.I32 // Default assumption
			if tt.name == "Access Struct Member (Data Ptr)" {
				retType = types.NewPointer(types.I32) // Expecting i32*
			}
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

			// Terminator should have been added by VisitReturnStatement
			if cg.Block != nil && cg.Block.Term == nil {
				t.Errorf("Block was not terminated after visiting return statement")
				// Add fallback just in case
				cg.Block.NewRet(constant.NewZeroInitializer(retType))
			}

			ir := cg.Module.String()
			t.Logf("Generated IR for %q:\n%s", tt.input, ir)

			// Check for GEP instruction
			reGEP := regexp.MustCompile(tt.expectedGEPRe)
			if !reGEP.MatchString(ir) {
				t.Errorf("Generated IR missing expected GEP instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedGEPRe, ir)
			}

			// Check for Load instruction
			reLoad := regexp.MustCompile(tt.expectedLoadRe)
			if !reLoad.MatchString(ir) {
				t.Errorf("Generated IR missing expected Load instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedLoadRe, ir)
			}

			// Check for Return instruction
			reRet := regexp.MustCompile(tt.expectedReturnRe)
			if !reRet.MatchString(ir) {
				t.Errorf("Generated IR missing expected Return instruction.\nExpected pattern: %s\nGot IR:\n%s", tt.expectedReturnRe, ir)
			}
		})
	}
}
