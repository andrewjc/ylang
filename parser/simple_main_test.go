package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"testing"
)

func TestSimpleMain(t *testing.T) {
	input := "main() -> {\n    let process = (input) -> {\n      return input * 2;\n    };\n    \n    let values = [1, 2, 3, 4, 5];\n    values.map(process).forEach(print);\n}"
	lexer, _ := lexer.NewLexerFromString(input)
	parser := NewParser(lexer)

	program := parser.ParseProgram()
	if len(program.Statements) != 4 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	// the first statement should be a function definition (main)
	var statement ast.Statement = program.Statements[0]
	assertFunction(t, statement, "main")

}

func assertFunction(t *testing.T, statement ast.Statement, expectedName string) {
	stmt, ok := statement.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", statement)
	}

	fnDef, ok := stmt.Expression.(ast.ExpressionNode)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionDefinition. got=%T", stmt.Expression)
	}

	if fnDef.TokenLiteral() != expectedName {
		t.Fatalf("function.Name.Value not '%s'. got=%s", expectedName, fnDef.TokenLiteral())
	}
}
