package generator

import (
	"compiler/ast"
	"compiler/lexer"
	"compiler/parser"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"regexp"
	"strings"
	"testing"
)

func TestCodeGenInfixExpressionUnit(t *testing.T) {
	tests := []struct {
		name              string
		input             string // The infix expression
		setupInput        string // Optional setup code (like let statements)
		expectedOperation string // Expected LLVM IR operation (e.g., "add nsw", "sub nsw", "mul nsw", "sdiv")
		expectedResultRe  string // Regex for the result assignment (e.g., `%[0-9]+ = add nsw i32 %.*, %.*`)
		expectedRetRe     string // Regex for returning the result (e.g., `ret i32 %[0-9]+`)
	}{
		{
			name:              "Integer Addition",
			input:             `5 + 3`,
			setupInput:        "",
			expectedOperation: "add nsw", // nsw = No Signed Wrap
			expectedResultRe:  `(%[a-zA-Z0-9_.]+) = add nsw i32 5, 3`,
			expectedRetRe:     `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:              "Integer Subtraction",
			input:             `10 - 4`,
			setupInput:        "",
			expectedOperation: "sub nsw",
			expectedResultRe:  `(%[a-zA-Z0-9_.]+) = sub nsw i32 10, 4`,
			expectedRetRe:     `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:              "Integer Multiplication",
			input:             `6 * 7`,
			setupInput:        "",
			expectedOperation: "mul nsw",
			expectedResultRe:  `(%[a-zA-Z0-9_.]+) = mul nsw i32 6, 7`,
			expectedRetRe:     `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:              "Integer Division",
			input:             `20 / 5`,
			setupInput:        "",
			expectedOperation: "sdiv", // Signed division
			expectedResultRe:  `(%[a-zA-Z0-9_.]+) = sdiv i32 20, 5`,
			expectedRetRe:     `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:              "Identifier Addition",
			input:             `a + b`,
			setupInput:        `let a = 100; let b = 200;`,
			expectedOperation: "add nsw",
			// Expect loads before the add
			expectedResultRe: `(%[a-zA-Z0-9_.]+) = add nsw i32 %[a-zA-Z0-9_.]+, %[a-zA-Z0-9_.]+`,
			expectedRetRe:    `ret i32 %[a-zA-Z0-9_.]+`,
		},
		// Add tests for comparison operators (>, <, ==, !=) when implemented
		// {
		//  name:            "Greater Than",
		//  input:           `a > 5`,
		//  setupInput:      `let a = 10;`,
		//  expectedOperation: "icmp sgt", // Signed Greater Than
		//  expectedResultRe: `(%[a-zA-Z0-9_.]+) = icmp sgt i32 %[a-zA-Z0-9_.]+, 5`,
		//  expectedRetRe:    `ret i1 %[a-zA-Z0-9_.]+`, // Comparison returns boolean (i1)
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullInput := tt.setupInput + tt.input // Combine setup and expression
			node := parseExpr(t, fullInput)       // Parse the main expression part
			infixExpr, ok := node.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("Parsed node for expression is not InfixExpression, got %T", node)
			}

			// Setup codegen context
			cg := NewCodeGenerator()
			mainSig := types.NewFunc(types.I32) // Assume infix results in i32 (adjust for comparisons)
			mainFunc := cg.Module.NewFunc("main", mainSig.RetType)
			entryBlock := mainFunc.NewBlock("entry")
			cg.currentFunc = mainFunc
			cg.Block = entryBlock

			// Visit setup code first if any
			if tt.setupInput != "" {
				// Need to parse setup separately and visit each statement
				setupStmtsInput := fmt.Sprintf("main() -> { %s }", tt.setupInput)
				lSetup, _ := lexer.NewLexerFromString(setupStmtsInput)
				pSetup := parser.NewParser(lSetup)
				progSetup := pSetup.ParseProgram()
				// No error check - assume setup is valid for the test
				bodySetup, _ := progSetup.MainFunction.Body.(*ast.BlockStatement)
				for _, stmt := range bodySetup.Statements {
					errVisit := stmt.Accept(cg)
					if errVisit != nil {
						t.Fatalf("Codegen failed during setup visit for input %q: %v", tt.input, errVisit)
					}
				}
			}

			// Visit the target infix expression node
			errExpr := infixExpr.Accept(cg)
			if errExpr != nil {
				t.Fatalf("Codegen failed during infix expression visit for input %q: %v", tt.input, errExpr)
			}

			// Add terminator (return the result of the infix operation)
			if cg.Block != nil && cg.Block.Term == nil {
				if cg.lastValue != nil { // && cg.lastValue.Type().Equal(mainSig.RetType) { // Add type check back later
					cg.Block.NewRet(cg.lastValue)
				} else {
					t.Logf("Warning: cg.lastValue is nil after visiting infix expression %q. Adding default return.", tt.input)
					cg.Block.NewRet(constant.NewInt(types.I32, 0)) // Fallback
				}
			} else if cg.Block == nil {
				t.Fatalf("Codegen block became nil unexpectedly for input %q", tt.input)
			}

			ir := cg.Module.String()

			// Check if the expected operation exists
			if !strings.Contains(ir, tt.expectedOperation) {
				t.Errorf("Generated IR missing expected operation %q for input %q.\nGot IR:\n%s", tt.expectedOperation, tt.input, ir)
			}

			// Check for the result assignment instruction
			reResult := regexp.MustCompile(tt.expectedResultRe)
			if !reResult.MatchString(ir) {
				t.Errorf("Generated IR missing expected result assignment for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedResultRe, ir)
			}

			// Check for the return instruction
			reRet := regexp.MustCompile(tt.expectedRetRe)
			if !reRet.MatchString(ir) {
				t.Errorf("Generated IR missing expected return instruction for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedRetRe, ir)
			}
		})
	}
}

func TestCodeGenCallExpressionUnit(t *testing.T) {
	tests := []struct {
		name            string
		input           string // The call expression
		setupInput      string // Code to define functions/variables used
		expectedFuncSig string // Expected LLVM signature of the function being declared/called
		expectedCallRe  string // Regex for the LLVM call instruction
		expectedRetRe   string // Regex for returning the call result (if not void)
	}{
		{
			name:            "Call No-Arg Function",
			input:           `getId()`,
			setupInput:      `function getId() -> { return 42; }`,
			expectedFuncSig: `define i32 @getId()`,
			expectedCallRe:  `(%[a-zA-Z0-9_.]+) = call i32 @getId\(\)`,
			expectedRetRe:   `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:            "Call One-Arg Function (Literal)",
			input:           `double(5)`,
			setupInput:      `function double(n) -> { return n * 2; }`,
			expectedFuncSig: `define i32 @double\(i32 %n\)`, // Param name might differ
			expectedCallRe:  `(%[a-zA-Z0-9_.]+) = call i32 @double\(i32 5\)`,
			expectedRetRe:   `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:            "Call Multi-Arg Function (Literals)",
			input:           `add(10, 20)`,
			setupInput:      `function add(a, b) -> { return a + b; }`,
			expectedFuncSig: `define i32 @add\(i32 %a, i32 %b\)`,
			expectedCallRe:  `(%[a-zA-Z0-9_.]+) = call i32 @add\(i32 10, i32 20\)`,
			expectedRetRe:   `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:            "Call Function with Variable Args",
			input:           `add(x, y)`,
			setupInput:      `let x = 1; let y = 2; function add(a, b) -> { return a + b; }`,
			expectedFuncSig: `define i32 @add\(i32 %a, i32 %b\)`,
			// Expect loads for x and y before the call
			expectedCallRe: `(%[a-zA-Z0-9_.]+) = call i32 @add\(i32 %[a-zA-Z0-9_.]+, i32 %[a-zA-Z0-9_.]+\)`, // Args are loaded values
			expectedRetRe:  `ret i32 %[a-zA-Z0-9_.]+`,
		},
		{
			name:            "Call Function Assigned to Variable",
			input:           `fnPtr(50)`,
			setupInput:      `function multiply(n) -> n*2; let fnPtr = multiply;`,
			expectedFuncSig: `define internal i32 @multiply\(i32 %n\)`, // Internal linkage likely
			// Expect alloc for fnPtr, store of @multiply, load of ptr, then call via ptr
			expectedCallRe: `(%[a-zA-Z0-9_.]+) = call i32 %[a-zA-Z0-9_.]+\(i32 50\)`, // Call via loaded pointer
			expectedRetRe:  `ret i32 %[a-zA-Z0-9_.]+`,
		},
		// Add test for void function call
		// {
		//  name:           "Call Void Function",
		//  input:          `print("hello")`,
		//  setupInput:     `function print(s) -> { /* syscall or similar */ }`, // Assume print returns void
		//  expectedFuncSig: `define void @print\(ptr %s\)`, // Assuming string is ptr
		//  expectedCallRe: `call void @print\(ptr %[a-zA-Z0-9_.]+\)`, // No result assignment
		//  expectedRetRe:  `ret i32 0`, // Default return 0 from main
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullInput := tt.setupInput + tt.input // Combine setup and call
			node := parseExpr(t, fullInput)       // Parse the call expression part
			callExpr, ok := node.(*ast.CallExpression)
			if !ok {
				t.Fatalf("Parsed node for expression is not CallExpression, got %T", node)
			}

			// Setup codegen context
			cg := NewCodeGenerator()
			// Result type depends on called function, assume i32 for non-void for now
			mainSig := types.NewFunc(types.I32)
			mainFunc := cg.Module.NewFunc("main", mainSig.RetType)
			entryBlock := mainFunc.NewBlock("entry")
			cg.currentFunc = mainFunc
			cg.Block = entryBlock

			// Visit setup code first
			if tt.setupInput != "" {
				// Need to parse setup into a Program and visit it
				lSetup, _ := lexer.NewLexerFromString(tt.setupInput)
				pSetup := parser.NewParser(lSetup)
				progSetup := pSetup.ParseProgram()
				// No error check - assume setup is valid for the test
				errVisit := progSetup.Accept(cg) // Visit the whole setup program
				if errVisit != nil {
					t.Fatalf("Codegen failed during setup visit for input %q: %v", tt.input, errVisit)
				}
				// Reset block pointer as visiting setup might change it
				cg.Block = entryBlock
			}

			// Visit the target call expression node
			errCall := callExpr.Accept(cg)
			if errCall != nil {
				t.Fatalf("Codegen failed during call expression visit for input %q: %v", tt.input, errCall)
			}

			// Add terminator (return the result or default)
			if cg.Block != nil && cg.Block.Term == nil {
				// Determine if the called function was expected to be void based on the regex
				isVoidReturn := strings.Contains(tt.expectedCallRe, "call void")

				// The dummy main in the test always returns i32.
				// If the called function was *not* void AND cg.lastValue holds the result,
				// return that result (assuming it's compatible with i32 for this test).
				// Otherwise, return the default 0 for the dummy main.
				if !isVoidReturn && cg.lastValue != nil {
					// Basic check: if lastValue is i32, return it.
					// More robust checks/casting might be needed if the called function
					// could return other types convertible to i32.
					if cg.lastValue.Type().Equal(types.I32) {
						cg.Block.NewRet(cg.lastValue)
					} else {
						// If the type isn't i32, just return 0 for the test main for now.
						// Or potentially log an error if the test expected a specific return type.
						t.Logf("Warning: Call result type (%s) doesn't match test main's i32 return. Returning default 0.", cg.lastValue.Type())
						cg.Block.NewRet(constant.NewInt(types.I32, 0))
					}
				} else {
					// Return default 0 if the called function was void,
					// or if cg.lastValue was nil (e.g., codegen error).
					cg.Block.NewRet(constant.NewInt(types.I32, 0))
				}
			} else if cg.Block == nil {
				t.Fatalf("Codegen block became nil unexpectedly for input %q", tt.input)
			}

			ir := cg.Module.String()

			// Check if the function was declared/defined with expected signature
			reSig := regexp.MustCompile(tt.expectedFuncSig)
			if !reSig.MatchString(ir) {
				t.Errorf("Generated IR missing expected function signature for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedFuncSig, ir)
			}

			// Check for the call instruction
			reCall := regexp.MustCompile(tt.expectedCallRe)
			if !reCall.MatchString(ir) {
				t.Errorf("Generated IR missing expected call instruction for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedCallRe, ir)
			}

			// Check for the return instruction
			reRet := regexp.MustCompile(tt.expectedRetRe)
			if !reRet.MatchString(ir) {
				t.Errorf("Generated IR missing expected return instruction for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedRetRe, ir)
			}
		})
	}
}
