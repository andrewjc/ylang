package parser

import (
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

	if len(program.ClassDeclarations) != 1 {
		t.Fatalf("program.ClassDeclarations does not contain 1 statement. got=%d", len(program.ClassDeclarations))
	}

	stmt := program.ClassDeclarations[0]
	if stmt.TokenLiteral() != "let" {
		t.Fatalf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
	}

	if stmt.Name.Value != "x" {
		t.Fatalf("stmt.Name.Value not 'x'. got=%q", stmt.Name.Value)
	}
}
