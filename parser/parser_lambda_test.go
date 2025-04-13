package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestLambdaExpressionParsingIntegration(t *testing.T) {
	tests := []struct {
		name            string
		input           string // Input snippet, usually within a `main` function context
		expectedParams  []string
		isBodyBlock     bool
		expectedBody    string
		expectedErrors  int
		lambdaExtractor func(*ast.Program) *ast.LambdaExpression // Function to find the lambda in the parsed AST
	}{
		{
			name:           "Lambda assigned to variable - no params, block body",
			input:          `main() -> { let f = () -> { return 1; }; }`,
			expectedParams: []string{},
			isBodyBlock:    true,
			expectedBody:   "{\n    return 1;\n}",
			expectedErrors: 0,
			lambdaExtractor: func(prog *ast.Program) *ast.LambdaExpression {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*ast.BlockStatement); ok && len(body.Statements) == 1 {
						if letStmt, ok := body.Statements[0].(*ast.LetStatement); ok {
							if lam, ok := letStmt.Value.(*ast.LambdaExpression); ok {
								return lam
							}
						}
					}
				}
				return nil
			},
		},
		{
			name:           "Lambda assigned to variable - one param, expr body",
			input:          `main() -> { let double = (x) -> x * 2; }`,
			expectedParams: []string{"x"},
			isBodyBlock:    false,
			expectedBody:   "(x * 2)",
			expectedErrors: 0,
			lambdaExtractor: func(prog *ast.Program) *ast.LambdaExpression {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*ast.BlockStatement); ok && len(body.Statements) == 1 {
						if letStmt, ok := body.Statements[0].(*ast.LetStatement); ok {
							if lam, ok := letStmt.Value.(*ast.LambdaExpression); ok {
								return lam
							}
						}
					}
				}
				return nil
			},
		},
		{
			name:           "Lambda assigned to variable - multi params, block body",
			input:          `main() -> { let add = (a, b) -> { return a + b; }; }`,
			expectedParams: []string{"a", "b"},
			isBodyBlock:    true,
			expectedBody:   "{\n    return (a + b);\n}",
			expectedErrors: 0,
			lambdaExtractor: func(prog *ast.Program) *ast.LambdaExpression {
				// Same extraction logic as above
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*ast.BlockStatement); ok && len(body.Statements) == 1 {
						if letStmt, ok := body.Statements[0].(*ast.LetStatement); ok {
							if lam, ok := letStmt.Value.(*ast.LambdaExpression); ok {
								return lam
							}
						}
					}
				}
				return nil
			},
		},
		{
			name:           "Lambda as call argument",
			input:          `main() -> { process((n) -> n > 0); }`,
			expectedParams: []string{"n"},
			isBodyBlock:    false,
			expectedBody:   "(n > 0)",
			expectedErrors: 0,
			lambdaExtractor: func(prog *ast.Program) *ast.LambdaExpression {
				// Extract lambda from call expression argument
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*ast.BlockStatement); ok && len(body.Statements) == 1 {
						if exprStmt, ok := body.Statements[0].(*ast.ExpressionStatement); ok {
							if callExpr, ok := exprStmt.Expression.(*ast.CallExpression); ok && len(callExpr.Arguments) == 1 {
								if lam, ok := callExpr.Arguments[0].(*ast.LambdaExpression); ok {
									return lam
								}
							}
						}
					}
				}
				return nil
			},
		},
		{
			name:           "Lambda as return value (requires nested parsing test)",
			input:          `function makeAdder() -> { return (x) -> { return (y) -> x + y; }; }`,
			expectedParams: []string{"x"}, // Outer lambda
			isBodyBlock:    true,
			expectedBody:   "{\n    return (y) -> (x + y);\n}", // Outer lambda's body contains inner lambda
			expectedErrors: 0,
			lambdaExtractor: func(prog *ast.Program) *ast.LambdaExpression {
				// Extract lambda from return statement inside the function
				if len(prog.Functions) == 1 {
					fn := prog.Functions[0]
					if body, ok := fn.Body.(*ast.BlockStatement); ok && len(body.Statements) == 1 {
						if retStmt, ok := body.Statements[0].(*ast.ReturnStatement); ok {
							if lam, ok := retStmt.ReturnValue.(*ast.LambdaExpression); ok {
								return lam
							}
						}
					}
				}
				return nil
			},
		},
		{
			name:            "Malformed lambda - missing arrow",
			input:           `main() -> { let f = (a, b) { return a; }; }`,
			expectedParams:  []string{"a", "b"}, // Params might be parsed
			isBodyBlock:     true,               // Body might be parsed depending on recovery
			expectedBody:    "{\n    return a;\n}",
			expectedErrors:  1,                                                            // Error for missing ->
			lambdaExtractor: func(prog *ast.Program) *ast.LambdaExpression { return nil }, // Don't check AST on error
		},
		{
			name:            "Malformed lambda - missing params paren",
			input:           `main() -> { let f = a, b) -> a + b; }`,
			expectedParams:  []string{},
			isBodyBlock:     false,
			expectedBody:    "",
			expectedErrors:  1, // Error for missing '('
			lambdaExtractor: func(prog *ast.Program) *ast.LambdaExpression { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()

			if tt.expectedErrors > 0 {
				if len(p.Errors()) < tt.expectedErrors {
					t.Errorf("Expected at least %d parser errors, but got %d:", tt.expectedErrors, len(p.Errors()))
					for i, e := range p.Errors() {
						t.Errorf("  Error %d: %s", i+1, e)
					}
				}
				return // Don't check AST if errors expected
			} else {
				checkParserErrors(t, p) // Fail on unexpected errors
			}

			lambda := tt.lambdaExtractor(program)
			if lambda == nil {
				t.Fatalf("Failed to extract lambda expression from parsed AST for input: %s", tt.input)
			}

			// Check Parameters
			if len(lambda.Parameters) != len(tt.expectedParams) {
				t.Errorf("Parameter count mismatch. want=%d, got=%d", len(tt.expectedParams), len(lambda.Parameters))
			} else {
				for i, expectedParam := range tt.expectedParams {
					if lambda.Parameters[i].Value != expectedParam {
						t.Errorf("Parameter %d mismatch. want=%s, got=%s", i, expectedParam, lambda.Parameters[i].Value)
					}
				}
			}

			// Check Body Type and Content
			if lambda.Body == nil {
				t.Fatalf("Lambda body is nil")
			}
			_, isBlock := lambda.Body.(*ast.BlockStatement)
			if isBlock != tt.isBodyBlock {
				t.Errorf("Lambda body type mismatch. want block=%v, got block=%v (Type: %T)", tt.isBodyBlock, isBlock, lambda.Body)
			}
			actualBodyStr := lambda.Body.String() // Use String() which includes indentation for blocks
			if actualBodyStr != tt.expectedBody {
				t.Errorf("Lambda body string representation mismatch.\nWant: %s\nGot:  %s", tt.expectedBody, actualBodyStr)
			}
		})
	}
}
