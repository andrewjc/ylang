package parser

import (
	"compiler/lexer"
	"testing"
)

func TestSimpleMain(t *testing.T) {
	input := "main() -> {\n    let process = (input) -> {\n      return input * 2;\n    };\n    \n    let values = [1, 2, 3, 4, 5];\n    values.map(process).forEach(print);\n}"
	lexer, _ := lexer.NewLexerFromString(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

}
