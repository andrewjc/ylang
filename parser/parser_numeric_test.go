package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestVariableDeclarationNumberLiteral(t *testing.T) {
	input := "let x = 5;"
	lexer, _ := lexer.NewLexerFromString(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram() // Assuming you have a method to parse the entire program
	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VariableDeclaration)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VariableDeclaration. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Value.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("exp not *ast.NumberLiteral. got=%T", stmt.Value)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %f. got=%f", 5.0, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}
