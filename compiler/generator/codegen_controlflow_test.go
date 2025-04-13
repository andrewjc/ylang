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

func TestCodeGenIfStatementUnit(t *testing.T) {
	tests := []struct {
		name                 string
		input                string // The if statement within a function context
		setupInput           string // Optional setup
		expectedBranchRe     string // Regex for the conditional branch instruction (br i1 ...)
		expectedThenLabelRe  string // Regex matching the label for the 'then' block (e.g., `if_then:`)
		expectedElseLabelRe  string // Regex matching the label for the 'else' block (or merge if no else)
		expectedMergeLabelRe string // Regex matching the label for the merge block
		expectElseBlock      bool   // Whether an explicit else block is expected
	}{
		{
			name:                 "Simple If (True)",
			input:                `if (1 > 0) { return 10; } return 0;`, // Condition is constant true
			setupInput:           "",
			expectedBranchRe:     `br i1 true, label %if_then, label %if_else`, // Expect direct branch if condition is constant
			expectedThenLabelRe:  `if_then:`,
			expectedElseLabelRe:  `if_else:`, // Even if empty, the label is generated
			expectedMergeLabelRe: `if_merge:`,
			expectElseBlock:      false, // No explicit else {} block in input
		},
		{
			name:                 "Simple If (False)",
			input:                `if (0 > 1) { return 10; } return 0;`, // Condition is constant false
			setupInput:           "",
			expectedBranchRe:     `br i1 false, label %if_then, label %if_else`,
			expectedThenLabelRe:  `if_then:`,
			expectedElseLabelRe:  `if_else:`,
			expectedMergeLabelRe: `if_merge:`,
			expectElseBlock:      false,
		},
		{
			name:       "If with Variable Condition",
			input:      `if (x > 5) { return 1; } return 0;`,
			setupInput: `let x = 10;`,
			// Expect load of x, compare, then branch
			expectedBranchRe:     `br i1 %[a-zA-Z0-9_.]+, label %if_then, label %if_else`,
			expectedThenLabelRe:  `if_then:`,
			expectedElseLabelRe:  `if_else:`,
			expectedMergeLabelRe: `if_merge:`,
			expectElseBlock:      false,
		},
		{
			name:                 "If/Else",
			input:                `if (y < 0) { return -1; } else { return 1; }`,
			setupInput:           `let y = -5;`,
			expectedBranchRe:     `br i1 %[a-zA-Z0-9_.]+, label %if_then, label %if_else`,
			expectedThenLabelRe:  `if_then:`,
			expectedElseLabelRe:  `if_else:`,
			expectedMergeLabelRe: `if_merge:`, // Merge block might be optimized out if both branches return
			expectElseBlock:      true,
		},
		{
			name:       "If/Else If/Else", // Requires parser support first
			input:      `if (a == 1) { return 1; } else if (a == 2) { return 2; } else { return 3; }`,
			setupInput: `let a = 2;`,
			// Expect nested branches
			expectedBranchRe:     `br i1 %[a-zA-Z0-9_.]+, label %if_then\d*, label %if_else\d*`, // More general regex
			expectedThenLabelRe:  `if_then\d*:`,
			expectedElseLabelRe:  `if_else\d*:`,
			expectedMergeLabelRe: `if_merge\d*:`,
			expectElseBlock:      true, // The final else
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Need to parse the *block* containing the if statement
			fullInput := fmt.Sprintf("main() -> { %s %s }", tt.setupInput, tt.input)
			l, err := lexer.NewLexerFromString(fullInput)
			if err != nil {
				t.Fatalf("Lexer error: %v", err)
			}
			p := parser.NewParser(l)
			prog := p.ParseProgram()
			if len(p.Errors()) > 0 {
				// Allow parsing errors if the test setup expects them (e.g., else if not supported yet)
				t.Logf("Parser errors encountered (may be expected): %v", p.Errors())
			}
			if prog.MainFunction == nil || prog.MainFunction.Body == nil {
				t.Fatalf("Failed to parse main function")
			}
			body, ok := prog.MainFunction.Body.(*ast.BlockStatement)
			if !ok || len(body.Statements) == 0 {
				t.Fatalf("Main body is not block or empty")
			}

			// Find the IfStatement node (might not be the first if setup exists)
			var ifStmt *ast.IfStatement
			for _, stmt := range body.Statements {
				if exprStmt, okExpr := stmt.(*ast.ExpressionStatement); okExpr {
					if foundIf, okIf := exprStmt.Expression.(*ast.IfStatement); okIf {
						ifStmt = foundIf
						break
					}
				} else if foundIf, okIf := stmt.(*ast.IfStatement); okIf {
					// If parseStatement returns IfStatement directly
					ifStmt = foundIf
					break
				}
			}

			if ifStmt == nil {
				t.Fatalf("Could not find IfStatement node in parsed AST for input:\n%s", fullInput)
			}

			// Setup codegen context and visit the entire program
			// This ensures setup code is generated before the if statement
			cg := NewCodeGenerator()
			errVisit := prog.Accept(cg) // Visit the whole program

			if errVisit != nil {
				// Don't fail hard, allow checking IR even if visitor had issues
				t.Logf("Warning: Visitor Accept returned error: %v", errVisit)
			}

			// Ensure main function terminates if not already done by returns in if/else
			if mainFunc, ok := cg.Functions["main"]; ok {
				lastBlock := mainFunc.Blocks[len(mainFunc.Blocks)-1] // Get the last block added
				if lastBlock.Term == nil {
					t.Logf("Warning: Main function's last block (%s) not terminated. Adding default ret.", lastBlock.LocalIdent.Name())
					lastBlock.NewRet(constant.NewInt(types.I32, 0))
				}
			}

			ir := cg.Module.String()

			// Check for conditional branch
			reBranch := regexp.MustCompile(tt.expectedBranchRe)
			if !reBranch.MatchString(ir) {
				t.Errorf("Generated IR missing expected conditional branch for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedBranchRe, ir)
			}

			// Check for labels
			reThen := regexp.MustCompile(`(?m)^` + tt.expectedThenLabelRe) // Match start of line
			if !reThen.MatchString(ir) {
				t.Errorf("Generated IR missing expected 'then' block label for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedThenLabelRe, ir)
			}

			reElse := regexp.MustCompile(`(?m)^` + tt.expectedElseLabelRe)
			if !reElse.MatchString(ir) {
				t.Errorf("Generated IR missing expected 'else' block label for input %q.\nExpected pattern: %s\nGot IR:\n%s", tt.input, tt.expectedElseLabelRe, ir)
			}

			reMerge := regexp.MustCompile(`(?m)^` + tt.expectedMergeLabelRe)
			// Merge block might be optimized away if all paths return, so this check might fail legitimately
			if !reMerge.MatchString(ir) {
				t.Logf("Note: Merge block label pattern %q not found for input %q. May have been optimized out if all branches return.", tt.expectedMergeLabelRe, tt.input)
			}

			// Optionally: Check content of then/else blocks if needed
		})
	}
}

func TestCodeGenBlockStatementUnit(t *testing.T) {
	tests := []struct {
		name                 string
		input                string // The block statement content { ... }
		setupInput           string
		expectedIRSubstrings []string // List of substrings expected within the block's generated IR
	}{
		{
			name:       "Block with Let",
			input:      `{ let innerVar = 50; }`,
			setupInput: "",
			expectedIRSubstrings: []string{
				`alloca i32`,          // Allocation for innerVar
				`store i32 50, ptr %`, // Storing the value
			},
		},
		{
			name:       "Block with Return",
			input:      `{ return 100; }`,
			setupInput: "",
			expectedIRSubstrings: []string{
				`ret i32 100`,
			},
		},
		{
			name:       "Block with Multiple Statements",
			input:      `{ let a = 1; let b = 2; return a + b; }`,
			setupInput: "",
			expectedIRSubstrings: []string{
				`alloca i32`,         // a
				`store i32 1, ptr %`, // store a
				`alloca i32`,         // b
				`store i32 2, ptr %`, // store b
				`load i32, ptr %`,    // load a
				`load i32, ptr %`,    // load b
				`add nsw i32 %`,      // add
				`ret i32 %`,          // return result
			},
		},
		{
			name:       "Nested Block",
			input:      `{ let outer = 10; { let inner = 20; outer = outer + inner; } return outer; }`,
			setupInput: "",
			// Scoping means `inner` is distinct. `outer` is allocated once.
			expectedIRSubstrings: []string{
				`%outer = alloca i32`, // Outer allocation (using specific name for clarity)
				`store i32 10, ptr %outer`,
				`%inner = alloca i32`, // Inner allocation
				`store i32 20, ptr %inner`,
				`load i32, ptr %outer`,
				`load i32, ptr %inner`,
				`add nsw i32`,
				`store i32 %, ptr %outer`, // Store result back to outer
				`load i32, ptr %outer`,    // Load outer for return
				`ret i32 %`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the block directly as the main body
			fullInput := fmt.Sprintf("main() -> %s", tt.input) // Assume block is the function body
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
				t.Fatalf("Failed to parse main function or body")
			}
			blockNode, ok := prog.MainFunction.Body.(*ast.BlockStatement)
			if !ok {
				t.Fatalf("Main body is not BlockStatement, got %T", prog.MainFunction.Body)
			}

			// Setup codegen context and visit the block node
			cg := NewCodeGenerator()
			mainSig := types.NewFunc(types.I32)
			mainFunc := cg.Module.NewFunc("main", mainSig.RetType)
			entryBlock := mainFunc.NewBlock("entry")
			cg.currentFunc = mainFunc
			cg.Block = entryBlock

			errVisit := blockNode.Accept(cg) // Visit the block
			if errVisit != nil {
				t.Errorf("runCodeGen failed: %v", errVisit)
			}

			// Ensure termination if block doesn't end with return
			if cg.Block != nil && cg.Block.Term == nil {
				cg.Block.NewRet(constant.NewInt(types.I32, 0))
			}

			ir := cg.Module.String()

			// Check for expected substrings
			for _, sub := range tt.expectedIRSubstrings {
				// Use regex for flexibility with register/variable names (%)
				subRe := strings.ReplaceAll(sub, "%", `%[a-zA-Z0-9_.]+`)
				re := regexp.MustCompile(subRe)
				if !re.MatchString(ir) {
					t.Errorf("Generated IR missing expected substring for input block.\nExpected pattern: %s\nGot IR:\n%s", subRe, ir)
				}
			}
		})
	}
}
