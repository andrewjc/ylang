package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestVariableDeclarationNumberLiteral(t *testing.T) {
	input := "main() -> {let x = 5;}"
	lexer, _ := lexer.NewLexerFromString(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram() // Assuming you have a method to parse the entire program
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if program.MainFunction == nil {
		t.Fatalf("ParseProgram() returned nil MainFunction")
	}

	body, ok := program.MainFunction.Body.(*ast.BlockStatement)
	if !ok {
		t.Fatalf("MainFunction.Body is not a BlockStatement")
	}

	if len(body.Statements) != 1 {
		t.Fatalf("body.Statements does not contain 1 statement. got=%d", len(body.Statements))
	}

	stmt, ok := body.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("body.Statements[0] is not *ast.LetStatement. got=%T", body.Statements[0])
	}

	if stmt.TokenLiteral() != "let" {
		t.Fatalf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
	}

	if stmt.Name.Value != "x" {
		t.Fatalf("stmt.Name.Value not 'x'. got=%q", stmt.Name.Value)
	}
}
