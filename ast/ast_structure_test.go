package ast

import (
	"compiler/lexer"
	"compiler/parser"
	"reflect"
	"testing"
)

// checkParserErrorsAST is a helper specifically for AST tests
func checkParserErrorsAST(t *testing.T, p *parser.Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors during AST construction test setup:", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow() // Stop test execution if there are parsing errors
}

func TestASTNodeTypesUnit(t *testing.T) {
	tests := []struct {
		name             string
		input            string              // Input snippet resulting in the target node
		expectedNodeType reflect.Type        // Expected Go type of the AST node
		nodeExtractor    func(*Program) Node // Function to extract the target node
	}{
		{
			name:             "Let Statement",
			input:            `main() -> { let x = 5; }`,
			expectedNodeType: reflect.TypeOf(&LetStatement{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						return body.Statements[0] // Should be the LetStatement
					}
				}
				return nil
			},
		},
		{
			name:             "Return Statement",
			input:            `main() -> { return "hello"; }`,
			expectedNodeType: reflect.TypeOf(&ReturnStatement{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						return body.Statements[0] // Should be the ReturnStatement
					}
				}
				return nil
			},
		},
		{
			name:             "Expression Statement (Infix)",
			input:            `main() -> { x + y; }`,
			expectedNodeType: reflect.TypeOf(&ExpressionStatement{}), // The statement itself
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						return body.Statements[0] // The ExpressionStatement
					}
				}
				return nil
			},
		},
		{
			name:             "Infix Expression",
			input:            `main() -> { a * b; }`,
			expectedNodeType: reflect.TypeOf(&InfixExpression{}), // The expression inside the statement
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The InfixExpression
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Identifier",
			input:            `main() -> { myVar; }`,
			expectedNodeType: reflect.TypeOf(&Identifier{}), // The identifier inside the expr stmt
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The Identifier
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Number Literal",
			input:            `main() -> { 123; }`,
			expectedNodeType: reflect.TypeOf(&NumberLiteral{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The NumberLiteral
						}
					}
				}
				return nil
			},
		},
		{
			name:             "String Literal",
			input:            `main() -> { "abc"; }`,
			expectedNodeType: reflect.TypeOf(&StringLiteral{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The StringLiteral
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Array Literal",
			input:            `main() -> { [1, 2]; }`,
			expectedNodeType: reflect.TypeOf(&ArrayLiteral{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The ArrayLiteral
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Index Expression",
			input:            `main() -> { arr[0]; }`,
			expectedNodeType: reflect.TypeOf(&IndexExpression{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The IndexExpression
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Lambda Expression",
			input:            `main() -> { (x) -> x * x; }`,
			expectedNodeType: reflect.TypeOf(&LambdaExpression{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The LambdaExpression
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Call Expression",
			input:            `main() -> { func(arg1); }`,
			expectedNodeType: reflect.TypeOf(&CallExpression{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The CallExpression
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Assignment Expression",
			input:            `main() -> { a = 10; }`,
			expectedNodeType: reflect.TypeOf(&AssignmentExpression{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The AssignmentExpression
						}
					}
				}
				return nil
			},
		},
		{
			name:             "Block Statement",
			input:            `main() -> { { let inner = 1; } }`,
			expectedNodeType: reflect.TypeOf(&BlockStatement{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						// Assuming the inner block is wrapped in an ExprStmt
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The inner BlockStatement
						}
						// Or maybe it's parsed directly if parseStatement handles '{'
						return body.Statements[0]
					}
				}
				return nil
			},
		},
		{
			name:             "If Statement",
			input:            `main() -> { if (true) {} }`,
			expectedNodeType: reflect.TypeOf(&IfStatement{}),
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						// Assuming If is wrapped in ExprStmt or returned directly by parseStatement
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression
						}
						return body.Statements[0] // The IfStatement
					}
				}
				return nil
			},
		},
		// Add tests for other node types: PrefixExpression, ClassDeclaration, DataStructure, etc. when supported
		{
			name:             "Import Statement",
			input:            `import "foo";`,
			expectedNodeType: reflect.TypeOf(&ImportStatement{}),
			nodeExtractor: func(prog *Program) Node {
				if len(prog.ImportStatements) > 0 {
					return prog.ImportStatements[0]
				}
				return nil
			},
		},
		{
			name:             "Member Access Expression",
			input:            `main() -> { obj.field; }`,
			expectedNodeType: reflect.TypeOf(&MemberAccessExpression{}), // Assuming this node type exists
			nodeExtractor: func(prog *Program) Node {
				if main := prog.MainFunction; main != nil {
					if body, ok := main.Body.(*BlockStatement); ok && len(body.Statements) > 0 {
						if exprStmt, ok := body.Statements[0].(*ExpressionStatement); ok {
							return exprStmt.Expression // The MemberAccessExpression
						}
					}
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := parser.NewParser(l)
			program := p.ParseProgram()
			checkParserErrorsAST(t, p) // Use AST-specific error check helper

			node := tt.nodeExtractor(program)
			if node == nil {
				// Allow nil node if expected type is also nil (e.g., for expected parse failures)
				if tt.expectedNodeType != nil {
					t.Fatalf("Node extraction failed for input: %s", tt.input)
				}
				return // Test passes if node extraction failed as expected
			}

			actualType := reflect.TypeOf(node)

			if actualType != tt.expectedNodeType {
				t.Errorf("Node type mismatch.\nInput: %s\nWant: %v\nGot:  %v", tt.input, tt.expectedNodeType, actualType)
			}

			// Additional Check: Verify children association for InfixExpression
			if infixExpr, ok := node.(*InfixExpression); ok {
				if infixExpr.Left == nil {
					t.Errorf("InfixExpression.Left is nil for input: %s", tt.input)
				}
				if infixExpr.Right == nil {
					t.Errorf("InfixExpression.Right is nil for input: %s", tt.input)
				}
			}
			// Add similar checks for other nodes with required children (CallExpression args, etc.)
		})
	}
}
