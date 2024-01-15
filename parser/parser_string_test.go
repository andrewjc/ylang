package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestVariableDeclarationStringLiteral(t *testing.T) {
	testCases := []struct {
		input         string
		expectedValue string
		expectedToken string
	}{
		{"let x = 'Andrew C';", "Andrew C", "Andrew C"},
		{"let y = \"Another Example\";", "Another Example", "Another Example"},
		{"let z = '';", "", ""},
	}

	for _, tt := range testCases {
		lexer, _ := lexer.NewLexerFromString(tt.input)
		parser := NewParser(lexer)

		program := parser.ParseProgram() // Assuming method to parse the entire program
		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.VariableDeclaration)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.VariableDeclaration. got=%T", program.Statements[0])
		}

		literal, ok := stmt.Value.(*ast.StringLiteral)
		if !ok {
			t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Value)
		}
		if literal.Value != tt.expectedValue {
			t.Errorf("literal.Value not %s. got=%s", tt.expectedValue, literal.Value)
		}
		if literal.TokenLiteral() != tt.expectedToken {
			t.Errorf("literal.TokenLiteral not %s. got=%s", tt.expectedToken, literal.TokenLiteral())
		}
	}
}
