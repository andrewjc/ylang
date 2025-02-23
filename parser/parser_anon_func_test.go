package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestAnonymousFunctionParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
		expectedBody   string
		isBlockBody    bool
		expectedError  bool
		description    string
	}{
		{
			input: `
                main() -> {
                    let f = () -> { return 5; };
                }`,
			expectedParams: []string{},
			expectedBody:   "return 5;",
			isBlockBody:    true,
			expectedError:  false,
			description:    "Anonymous function with empty parameters and block body",
		},
		{
			input: `
                main() -> {
                    let f = (a) -> a + 1;
                }`,
			expectedParams: []string{"a"},
			expectedBody:   "a + 1",
			isBlockBody:    false,
			expectedError:  false,
			description:    "Anonymous function with one parameter and expression body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			lexer, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Failed to create lexer: %v", err)
			}

			parser := NewParser(lexer)
			program := parser.ParseProgram()

			if tt.expectedError {
				if len(parser.Errors()) == 0 {
					t.Errorf("Expected parser errors but got none for input: %s", tt.input)
				}
				return
			} else {
				if len(parser.Errors()) != 0 {
					t.Errorf("Parser errors: %v", parser.Errors())
					return
				}
			}

			var lambda *ast.LambdaExpression

			// Since the language requires a main function, we should look for the lambda within it
			if program.MainFunction != nil {
				lambda = extractLambdaFromMainFunction(program.MainFunction)
			}

			if lambda == nil {
				t.Errorf("No lambda expression found in AST for input: %s", tt.input)
				return
			}

			// Check parameters
			if len(lambda.Parameters) != len(tt.expectedParams) {
				t.Errorf("Expected %d parameters but got %d for input: %s", len(tt.expectedParams), len(lambda.Parameters), tt.input)
				return
			}

			for i, param := range lambda.Parameters {
				if param.Value != tt.expectedParams[i] {
					t.Errorf("Expected parameter %d to be %s but got %s for input: %s", i, tt.expectedParams[i], param.Value, tt.input)
				}
			}

			// Check body
			bodyString := lambda.Body.String()

			expectedBody := tt.expectedBody
			if bodyString != expectedBody {
				t.Errorf("Expected lambda body '%s' but got '%s' for input: %s", expectedBody, bodyString, tt.input)
			}
		})
	}
}

// Helper function to extract lambda from the main function
func extractLambdaFromMainFunction(mainFn *ast.FunctionDefinition) *ast.LambdaExpression {
	if mainFn == nil || mainFn.Body == nil {
		return nil
	}

	switch body := mainFn.Body.(type) {
	case *ast.BlockStatement:
		for _, stmt := range body.Statements {
			if letStmt, ok := stmt.(*ast.LetStatement); ok {
				if lamExp, ok := letStmt.Value.(*ast.LambdaExpression); ok {
					return lamExp
				}
			} else if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
				if lamExp, ok := exprStmt.Expression.(*ast.LambdaExpression); ok {
					return lamExp
				}
			}
			// Add other cases as needed
		}
	default:
		return nil
	}

	return nil
}
