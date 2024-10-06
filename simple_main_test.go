package main

import (
	c "compiler/compiler"
	l "compiler/lexer"
	p "compiler/parser"
	"testing"
)

func TestSimpleMain(t *testing.T) {

	input := `main() -> {
		let process = (input) -> {
			return input * 2;
		};
		
		let values = [1, 2, 3, 4, 5];
		values.map(process).forEach(print);
	}`

	lexer, err := l.NewLexerFromString(input)
	if err != nil {
		t.Fatalf("Failed to create lexer: %v", err)
	}

	parser := p.NewParser(lexer)
	program := parser.ParseProgram()
	if len(parser.Errors()) != 0 {
		for i, err := range parser.Errors() {
			t.Logf("Error %d: %v\n", i+1, err)
		}
		t.FailNow()
	}

	compilerInstance := c.NewCompiler(c.LLVM)
	result := compilerInstance.Compile(program)

	if len(result.Errors) != 0 {
		t.Fatalf("Compiler errors: %v", result.Errors)
	}

	if result.Output == "" {
		t.Fatalf("Compiler output is empty")
	}
}
