package parser

import (
	"compiler/lexer"
	"testing"
)

func TestAssignmentExpression(t *testing.T) {
	input := `main() -> {let x = 10; x = x + 5;}`

	lex, err := lexer.NewLexerFromString(input)
	if err != nil {
		t.Errorf("Failed to create lexer: %v", err)
	}

	parse := NewParser(lex)
	_ = parse.ParseProgram()
	if len(parse.Errors()) != 0 {
		t.Errorf("Parser errors: %v", parse.Errors())
	}

}
