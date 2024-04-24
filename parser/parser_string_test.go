package parser

import (
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
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
	}
}
