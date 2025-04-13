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

// Helper function to parse input and get the first expression's node
func parseExpr(t *testing.T, input string) ast.ExpressionNode {
	t.Helper()
	// Wrap in main for parsing
	fullInput := fmt.Sprintf("main() -> { %s; }", input)
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
	exprStmt, ok := body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		// Handle direct return/let if needed
		if retStmt, okRet := body.Statements[0].(*ast.ReturnStatement); okRet {
			return retStmt.ReturnValue
		}
		if letStmt, okLet := body.Statements[0].(*ast.LetStatement); okLet {
			return letStmt.Value
		}
		t.Fatalf("First statement is not ExpressionStatement, got %T", body.Statements[0])
	}
	return exprStmt.Expression
}

// Helper function to parse input and get the first statement node
func parseStmt(t *testing.T, input string) ast.Statement {
	t.Helper()
	// Wrap in main for parsing
	fullInput := fmt.Sprintf("main() -> { %s }", input) // No semicolon needed if it's the only stmt
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
	return body.Statements[0]
}

// Helper function to run codegen and check output
func runCodeGen(t *testing.T, node ast.Node) (string, error) {
	t.Helper()
	cg := NewCodeGenerator()

	// Create a dummy main function context for visiting statements/expressions
	mainSig := types.NewFunc(types.I32)
	mainFunc := cg.Module.NewFunc("main", mainSig.RetType)
	entryBlock := mainFunc.NewBlock("entry")
	cg.currentFunc = mainFunc
	cg.Block = entryBlock

	// Create a dummy Program AST node to initiate visit if needed
	// Or directly call Accept on the target node
	err := node.Accept(cg) // Assumes node has Accept method

	// Ensure block is terminated if Accept doesn't add a terminator
	if cg.Block != nil && cg.Block.Term == nil {
		// If the visited node should produce a value, add a dummy return
		// Otherwise, if it's void context (like let), add ret void/i32(0)
		if cg.lastValue != nil && cg.currentFunc.Sig.RetType != types.Void {
			// Basic type matching or casting needed here
			if cg.lastValue.Type().Equal(cg.currentFunc.Sig.RetType) {
				cg.Block.NewRet(cg.lastValue)
			} else {
				// Placeholder: just return 0 if types mismatch
				cg.Block.NewRet(constant.NewInt(types.I32, 0))
			}
		} else {
			// Default return for main
			cg.Block.NewRet(constant.NewInt(types.I32, 0))
		}
	}

	if err != nil {
		return "", err
	}

	return cg.Module.String(), nil
}

func TestCodeGenNumberLiteralUnit(t *testing.T) {
	tests := []struct {
		name             string
		input            string // The literal expression
		expectedIRSubstr string // Substring expected in the generated IR
	}{
		{
			name:             "Integer Literal",
			input:            `5`,
			expectedIRSubstr: `ret i32 5`, // Expecting the literal used in a return
		},
		{
			name:             "Large Integer Literal",
			input:            `1234567890`,
			expectedIRSubstr: `ret i32 1234567890`,
		},
		{
			name:             "Zero Literal",
			input:            `0`,
			expectedIRSubstr: `ret i32 0`,
		},
		// {
		// 	name:            "Float Literal", // LLVM float representation
		// 	input:           `3.14`,
		// 	expectedIRSubstr: `ret float 3.14`, // Adjust type if generator uses double
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := parseExpr(t, tt.input)
			if node == nil {
				t.Fatalf("Parsing failed for input: %s", tt.input)
			}

			// Directly visit the literal node
			numLit, ok := node.(*ast.NumberLiteral)
			if !ok {
				t.Fatalf("Parsed node is not NumberLiteral, got %T", node)
			}

			ir, err := runCodeGen(t, numLit)

			if err != nil {
				t.Errorf("runCodeGen failed: %v", err)
			}
			if !strings.Contains(ir, tt.expectedIRSubstr) {
				t.Errorf("Generated IR mismatch for input %q.\nExpected substring: %s\nGot IR:\n%s", tt.input, tt.expectedIRSubstr, ir)
			}
		})
	}
}

func TestCodeGenStringLiteralUnit(t *testing.T) {
	tests := []struct {
		name              string
		input             string // The literal expression
		expectedGlobalDef string // Expected global string definition (e.g., @str_0 = ... )
		expectedPtrLoad   string // Expected instruction using the pointer (e.g., getelementptr)
	}{
		{
			name:              "Simple String Literal",
			input:             `"hello"`,
			expectedGlobalDef: `@[a-zA-Z_0-9]+ = private unnamed_addr constant \[6 x i8] c"hello\00"`, // Regex for global name
			expectedPtrLoad:   `getelementptr inbounds \[6 x i8], ptr @[a-zA-Z_0-9]+, i32 0, i32 0`,
		},
		{
			name:              "Empty String Literal",
			input:             `""`,
			expectedGlobalDef: `@[a-zA-Z_0-9]+ = private unnamed_addr constant \[1 x i8] c"\00"`,
			expectedPtrLoad:   `getelementptr inbounds \[1 x i8], ptr @[a-zA-Z_0-9]+, i32 0, i32 0`,
		},
		{
			name:              "String Literal with Escape", // Generator must handle escapes correctly for global def
			input:             `"a\nb"`,
			expectedGlobalDef: `@[a-zA-Z_0-9]+ = private unnamed_addr constant \[4 x i8] c"a\0Ab\00"`, // Expect newline char \0A
			expectedPtrLoad:   `getelementptr inbounds \[4 x i8], ptr @[a-zA-Z_0-9]+, i32 0, i32 0`,
		},
	}
	// Reset global counter for predictable names (if feasible)
	stringCounter = 0

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := parseExpr(t, tt.input)
			strLit, ok := node.(*ast.StringLiteral)
			if !ok {
				t.Fatalf("Parsed node is not StringLiteral, got %T", node)
			}

			ir, err := runCodeGen(t, strLit)
			if err != nil {
				t.Errorf("runCodeGen failed: %v", err)
			}

			// Use regex for flexible matching of global names (@str_0, @str_1 etc.)
			reGlobal := regexp.MustCompile(tt.expectedGlobalDef)
			if !reGlobal.MatchString(ir) {
				t.Errorf("Generated IR missing expected global definition for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedGlobalDef, ir)
			}

			reLoad := regexp.MustCompile(tt.expectedPtrLoad)
			if !reLoad.MatchString(ir) {
				// The load might happen inside the return instruction setup
				expectedRetLoad := strings.Replace(tt.expectedPtrLoad, "getelementptr", "ret ptr getelementptr", 1)
				reRetLoad := regexp.MustCompile(expectedRetLoad)

				// Check if the GEP result is returned (codegen might optimize)
				ptrIdentRegex := regexp.MustCompile(`(%[a-zA-Z0-9_.]+) = ` + tt.expectedPtrLoad)
				matches := ptrIdentRegex.FindStringSubmatch(ir)
				returnedDirectly := false
				if len(matches) > 1 {
					ptrIdent := matches[1]
					retRegex := regexp.MustCompile(fmt.Sprintf(`ret ptr %s`, regexp.QuoteMeta(ptrIdent)))
					if retRegex.MatchString(ir) {
						returnedDirectly = true
					}
				}

				if !reRetLoad.MatchString(ir) && !returnedDirectly {
					t.Errorf("Generated IR missing expected pointer load/use for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedPtrLoad, ir)
				}
			}
		})
	}
}

func TestCodeGenLetStatementUnit(t *testing.T) {
	tests := []struct {
		name             string
		input            string // The let statement
		expectedAlloca   string // Expected alloca instruction (type and name)
		expectedStoreVal string // Expected value being stored (e.g., i32 5, ptr @...)
		expectedStore    string // Expected store instruction pattern
	}{
		{
			name:             "Let with Integer Literal",
			input:            `let myVar = 5;`,
			expectedAlloca:   `%myVar = alloca i32`, // Variable name might be different in IR
			expectedStoreVal: `i32 5`,
			expectedStore:    `store i32 5, ptr %[a-zA-Z0-9_.]+`, // Regex for var name ptr
		},
		{
			name:             "Let with String Literal",
			input:            `let myStr = "test";`,
			expectedAlloca:   `%myStr = alloca ptr`,                                                         // Stores i8* pointer
			expectedStoreVal: `ptr getelementptr inbounds \(\[5 x i8], ptr @[a-zA-Z0-9_.]+, i32 0, i32 0\)`, // GEP result
			expectedStore:    `store ptr %[a-zA-Z0-9_.]+, ptr %[a-zA-Z0-9_.]+`,
		},
		// Add test for let with identifier RHS later: let y = x;
		// Add test for let with expression RHS later: let z = a + b;
		// Add test for let with function call RHS later: let res = func();
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := parseStmt(t, tt.input) // Parse the statement directly
			letStmt, ok := node.(*ast.LetStatement)
			if !ok {
				t.Fatalf("Parsed node is not LetStatement, got %T", node)
			}

			ir, err := runCodeGen(t, letStmt) // Visit the LetStatement
			if err != nil {
				t.Errorf("runCodeGen failed: %v", err)
			}

			// Check for alloca instruction (flexible name matching)
			reAlloca := regexp.MustCompile(`%[a-zA-Z0-9_.]+ = alloca ` + strings.SplitN(tt.expectedAlloca, " ", 4)[3]) // Match type
			if !reAlloca.MatchString(ir) {
				t.Errorf("Generated IR missing expected alloca for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedAlloca, ir)
			}

			// Check for store instruction (flexible pointer matching)
			reStore := regexp.MustCompile(tt.expectedStore)
			if !reStore.MatchString(ir) {
				t.Errorf("Generated IR missing expected store for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedStore, ir)
			}

			// Check if the correct *kind* of value was stored (less strict)
			if !strings.Contains(ir, strings.Split(tt.expectedStoreVal, " ")[0]) { // Check type prefix e.g. "i32", "ptr"
				t.Errorf("Generated IR store instruction doesn't seem to store the right type of value for input %q.\nExpected value like: %s\nGot IR:\n%s", tt.input, tt.expectedStoreVal, ir)
			}
		})
	}
}

func TestCodeGenIdentifierUnit(t *testing.T) {
	// Setup: We need a variable defined first (e.g., via Let) so the identifier has something to load from.
	inputSetup := `let myVar = 123; `
	tests := []struct {
		name         string
		inputUse     string // Code that uses the identifier
		expectedLoad string // Expected load instruction pattern
		expectedRet  string // Expected return using the loaded value
	}{
		{
			name:         "Load Integer Identifier",
			inputUse:     `myVar`,                                           // Use the identifier as an expression
			expectedLoad: `%[a-zA-Z0-9_.]+ = load i32, ptr %[a-zA-Z0-9_.]+`, // Load from the alloca ptr
			expectedRet:  `ret i32 %[a-zA-Z0-9_.]+`,                         // Return the loaded value
		},
		// Add tests for loading other types (string ptr, etc.) when Let supports them fully
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullInput := inputSetup + tt.inputUse // Combine setup and usage
			node := parseExpr(t, fullInput)       // Parse the usage expression
			ident, ok := node.(*ast.Identifier)
			if !ok {
				t.Fatalf("Parsed node for usage is not Identifier, got %T", node)
			}

			// Code generation needs to handle the setup *and* the usage
			cg := NewCodeGenerator()
			mainSig := types.NewFunc(types.I32) // Assuming identifier use results in i32 return
			mainFunc := cg.Module.NewFunc("main", mainSig.RetType)
			entryBlock := mainFunc.NewBlock("entry")
			cg.currentFunc = mainFunc
			cg.Block = entryBlock

			// 1. Manually visit the setup statement(s) first to populate context
			setupStmt := parseStmt(t, inputSetup)
			errSetup := setupStmt.Accept(cg)
			if errSetup != nil {
				t.Fatalf("Codegen failed during setup visit: %v", errSetup)
			}
			if val, ok := cg.getVar("myVar"); !ok || val == nil { // Check if setup actually defined the var in context
				t.Fatalf("Variable 'myVar' not found in codegen context after visiting setup")
			}

			// 2. Now visit the identifier usage node
			errUse := ident.Accept(cg)
			if errUse != nil {
				t.Fatalf("Codegen failed during identifier visit: %v", errUse)
			}

			// 3. Add terminator (return the loaded value)
			if cg.Block != nil && cg.Block.Term == nil {
				if cg.lastValue != nil && cg.currentFunc.Sig.RetType != types.Void && cg.lastValue.Type().Equal(cg.currentFunc.Sig.RetType) {
					cg.Block.NewRet(cg.lastValue)
				} else {
					cg.Block.NewRet(constant.NewInt(types.I32, 0)) // Fallback return
				}
			}

			ir := cg.Module.String()

			// Check for load instruction
			reLoad := regexp.MustCompile(tt.expectedLoad)
			if !reLoad.MatchString(ir) {
				t.Errorf("Generated IR missing expected load for identifier %q.\nExpected pattern: %s\nGot IR:\n%s", ident.Value, tt.expectedLoad, ir)
			}

			// Check for return instruction using the loaded value
			reRet := regexp.MustCompile(tt.expectedRet)
			if !reRet.MatchString(ir) {
				t.Errorf("Generated IR missing expected return for identifier %q.\nExpected pattern: %s\nGot IR:\n%s", ident.Value, tt.expectedRet, ir)
			}
		})
	}
}
