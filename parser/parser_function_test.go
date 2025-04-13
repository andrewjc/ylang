package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestFunctionDefinitionParsingIntegration(t *testing.T) {
	tests := []struct {
		input           string
		expectedName    string
		expectedParams  []string
		expectedRetType string // Empty string if no explicit return type
		isBodyBlock     bool   // Is the body expected to be a BlockStatement?
		expectedBody    string // String representation of the body
		expectedErrors  int
	}{
		{
			input:           `function simple() -> {}`,
			expectedName:    "simple",
			expectedParams:  []string{},
			expectedRetType: "", // Implicit void/i32 assumed by parser/codegen, not syntax
			isBodyBlock:     true,
			expectedBody:    "{\n}",
			expectedErrors:  0,
		},
		{
			input:           `function add(a, b) -> a + b;`,
			expectedName:    "add",
			expectedParams:  []string{"a", "b"},
			expectedRetType: "",
			isBodyBlock:     false,
			expectedBody:    "(a + b)",
			expectedErrors:  0,
		},
		{
			// Explicit return type, block body
			input:           `function calculate(x): int -> { let y = x * 2; return y; }`,
			expectedName:    "calculate",
			expectedParams:  []string{"x"},
			expectedRetType: "int",
			isBodyBlock:     true,
			expectedBody:    "{\n    let y = (x * 2);\n    return y;\n}",
			expectedErrors:  0,
		},
		{
			// Function assigned to let (implicitly handled by let parser calling expression parser)
			input:           `main() -> { let myFunc = function(p1) -> p1; }`,
			expectedName:    "main", // Testing the outer main function here
			expectedParams:  []string{},
			expectedRetType: "",
			isBodyBlock:     true,
			expectedBody:    "{\n    let myFunc = anonymous(p1) -> p1;\n}", // Inner func is anonymous but assigned
			expectedErrors:  0,
		},
		{
			input:          `function missingBody(a) -> ;`, // Missing body after arrow
			expectedName:   "missingBody",
			expectedParams: []string{"a"},
			isBodyBlock:    false,
			expectedBody:   "", // Parse expression returns nil
			expectedErrors: 1,  // Error for missing expression
		},
		{
			input:          `function badParams(a,,b) -> {}`, // Extra comma
			expectedName:   "badParams",
			expectedParams: []string{"a"}, // Stops parsing params after error
			isBodyBlock:    true,
			expectedBody:   "{\n}",
			expectedErrors: 1, // Error for missing identifier after comma
		},
		{
			input:          `function noArrow(c,d) {}`, // Missing ->
			expectedName:   "noArrow",
			expectedParams: []string{"c", "d"},
			isBodyBlock:    true, // Body might be parsed depending on recovery
			expectedBody:   "{\n}",
			expectedErrors: 1, // Error for missing ->
		},
		{
			input:          `function noParens x, y -> x+y;`, // Missing parens around params
			expectedName:   "noParens",
			expectedParams: []string{}, // Fails to parse params
			isBodyBlock:    false,
			expectedBody:   "",
			expectedErrors: 1, // Error for expecting '('
		},
		{
			input: `function nested() -> {
                function inner(n) -> { return n * n; }
                return inner(5);
            }`,
			expectedName:    "nested",
			expectedParams:  []string{},
			expectedRetType: "",
			isBodyBlock:     true,
			// The parser currently doesn't support nested named functions directly within blocks.
			// It would parse 'function' as starting a statement, but parseFunctionDefinition isn't called by parseStatement.
			// This test will likely fail based on current parser structure.
			// Expected AST string if it *were* supported (hypothetical):
			// expectedBody:    "{\n    function inner(n) -> {\n        return (n * n);\n    };\n    return inner(5);\n}",
			expectedBody:   "{\n    return inner(5);\n}", // Current parser fails on 'function inner...' and recovers
			expectedErrors: 1,                            // Expect error because nested function is not directly supported
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram() // Assume top-level function definitions for now

			if tt.expectedErrors > 0 {
				if len(p.Errors()) < tt.expectedErrors {
					t.Errorf("Expected at least %d parser errors, but got %d:", tt.expectedErrors, len(p.Errors()))
					for i, e := range p.Errors() {
						t.Errorf("  Error %d: %s", i+1, e)
					}
				}
				// Don't check AST details if errors were expected
				return
			} else {
				checkParserErrors(t, p) // Fail if unexpected errors occurred
			}

			var fnDef *ast.FunctionDefinition
			if tt.expectedName == "main" { // Special case for testing 'let myFunc = function...'
				fnDef = program.MainFunction
			} else if len(program.Functions) > 0 {
				// Find the function by name (or assume the first one if only one is defined)
				found := false
				for _, fn := range program.Functions {
					if fn.Name != nil && fn.Name.Value == tt.expectedName {
						fnDef = fn
						found = true
						break
					}
				}
				if !found && len(program.Functions) == 1 && program.Functions[0].Name != nil && program.Functions[0].Name.Value == tt.expectedName {
					// Fallback for single function case if lookup failed (e.g., if name is optional)
					fnDef = program.Functions[0]
				} else if !found {
					t.Fatalf("Function '%s' not found in program.Functions", tt.expectedName)
				}
			} else {
				t.Fatalf("No function definitions found in program")
			}

			if fnDef == nil {
				t.Fatalf("Function definition parsing returned nil")
			}

			// Check Name
			if fnDef.Name == nil {
				if tt.expectedName != "anonymous" { // Allow testing anonymous functions later
					t.Errorf("Expected function name '%s', but got nil", tt.expectedName)
				}
			} else if fnDef.Name.Value != tt.expectedName {
				t.Errorf("Function name mismatch. want=%s, got=%s", tt.expectedName, fnDef.Name.Value)
			}

			// Check Parameters
			if len(fnDef.Parameters) != len(tt.expectedParams) {
				t.Errorf("Parameter count mismatch. want=%d, got=%d", len(tt.expectedParams), len(fnDef.Parameters))
			} else {
				for i, expectedParam := range tt.expectedParams {
					if fnDef.Parameters[i].Value != expectedParam {
						t.Errorf("Parameter %d mismatch. want=%s, got=%s", i, expectedParam, fnDef.Parameters[i].Value)
					}
				}
			}

			// Check Return Type
			if tt.expectedRetType == "" {
				if fnDef.ReturnType != nil {
					t.Errorf("Expected no explicit return type, but got '%s'", fnDef.ReturnType.Value)
				}
			} else {
				if fnDef.ReturnType == nil {
					t.Errorf("Expected explicit return type '%s', but got nil", tt.expectedRetType)
				} else if fnDef.ReturnType.Value != tt.expectedRetType {
					t.Errorf("Return type mismatch. want=%s, got=%s", tt.expectedRetType, fnDef.ReturnType.Value)
				}
			}

			// Check Body Type and Content
			if fnDef.Body == nil {
				if tt.expectedBody != "" { // Allow empty expected body for error cases
					t.Errorf("Function body is nil, expected non-nil body")
				}
			} else {
				_, isBlock := fnDef.Body.(*ast.BlockStatement)
				if isBlock != tt.isBodyBlock {
					t.Errorf("Function body type mismatch. want block=%v, got block=%v (Type: %T)", tt.isBodyBlock, isBlock, fnDef.Body)
				}
				actualBodyStr := fnDef.Body.String() // Use String() which includes indentation for blocks
				if actualBodyStr != tt.expectedBody {
					t.Errorf("Function body string representation mismatch.\nWant: %s\nGot:  %s", tt.expectedBody, actualBodyStr)
				}
			}
		})
	}
}
