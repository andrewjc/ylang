package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"fmt"
	"testing"
)

// checkParserErrors prints any errors found during parsing.
func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors:", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}
	t.FailNow() // Stop test execution if there are parsing errors
}

func TestLetStatementUnit(t *testing.T) {
	testCases := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{} // Can be int64, string, bool etc.
	}{
		{"main() -> {let x = 5;}", "x", int64(5)},
		{"main() -> {let y = 100;}", "y", int64(100)},
		{"main() -> {let foo = 12345;}", "foo", int64(12345)},
		{"main() -> {let name = \"bar\";}", "name", "bar"},
		{"main() -> {let empty = '';}", "empty", ""},
		// Add tests for other literal types (bool, etc.) if/when supported
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input_%s", tc.input), func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tc.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if program.MainFunction == nil {
				t.Fatalf("ParseProgram() returned nil MainFunction")
			}

			if len(program.MainFunction.Body.(*ast.BlockStatement).Statements) != 1 {
				t.Fatalf("MainFunction.Body does not contain 1 statement. got=%d", len(program.MainFunction.Body.(*ast.BlockStatement).Statements))
			}

			stmt, ok := program.MainFunction.Body.(*ast.BlockStatement).Statements[0].(*ast.LetStatement)
			if !ok {
				t.Fatalf("MainFunction.Body.Statements[0] is not *ast.LetStatement. got=%T", program.MainFunction.Body.(*ast.BlockStatement).Statements[0])
			}

			if stmt.Name.Value != tc.expectedIdentifier {
				t.Errorf("LetStatement.Name.Value not '%s'. got=%s", tc.expectedIdentifier, stmt.Name.Value)
			}

			if stmt.Name.TokenLiteral() != tc.expectedIdentifier {
				t.Errorf("LetStatement.Name.TokenLiteral() not '%s'. got=%s", tc.expectedIdentifier, stmt.Name.TokenLiteral())
			}

			// Test the value
			testLiteralExpression(t, stmt.Value, tc.expectedValue)
		})
	}
}

func TestReturnStatementUnit(t *testing.T) {
	testCases := []struct {
		input         string
		expectedValue interface{} // Can be int64, string, identifier etc.
	}{
		{"main() -> {return 5;}", int64(5)},
		{"main() -> {return 100;}", int64(100)},
		{"main() -> {return result;}", "result"}, // Returning an identifier
		{"main() -> {return \"hello\";}", "hello"},
		{"main() -> {return;}", nil}, // Return without value
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input_%s", tc.input), func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tc.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if program.MainFunction == nil {
				t.Fatalf("ParseProgram() returned nil MainFunction")
			}

			if len(program.MainFunction.Body.(*ast.BlockStatement).Statements) != 1 {
				t.Fatalf("MainFunction.Body does not contain 1 statement. got=%d", len(program.MainFunction.Body.(*ast.BlockStatement).Statements))
			}

			stmt, ok := program.MainFunction.Body.(*ast.BlockStatement).Statements[0].(*ast.ReturnStatement)
			if !ok {
				t.Fatalf("MainFunction.Body.Statements[0] is not *ast.ReturnStatement. got=%T", program.MainFunction.Body.(*ast.BlockStatement).Statements[0])
			}

			if stmt.TokenLiteral() != "return" {
				t.Errorf("ReturnStatement.TokenLiteral() not 'return'. got=%q", stmt.TokenLiteral())
			}

			// Test the return value (if any)
			if tc.expectedValue == nil {
				if stmt.ReturnValue != nil {
					t.Errorf("Expected ReturnValue to be nil, got %s", stmt.ReturnValue.String())
				}
			} else {
				if stmt.ReturnValue == nil {
					t.Errorf("Expected ReturnValue to be %v, got nil", tc.expectedValue)
				} else {
					testLiteralExpression(t, stmt.ReturnValue, tc.expectedValue)
				}
			}
		})
	}
}

func TestImportStatementUnit(t *testing.T) {
	testCases := []struct {
		input        string
		expectedPath string
	}{
		{`import "stdlib/core";`, "stdlib/core"},
		{`import "my_module/utils";`, "my_module/utils"},
		{`import "single_file";`, "single_file"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input_%s", tc.input), func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tc.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.ImportStatements) != 1 {
				t.Fatalf("Program.ImportStatements does not contain 1 statement. got=%d", len(program.ImportStatements))
			}

			stmt2 := program.ImportStatements[0]

			if stmt2.TokenLiteral() != "import" {
				t.Errorf("ImportStatement.TokenLiteral() not 'import'. got=%q", stmt2.TokenLiteral())
			}

			if stmt2.Path != tc.expectedPath {
				t.Errorf("ImportStatement.Path not '%s'. got=%s", tc.expectedPath, stmt2.Path)
			}
		})
	}
}

// Helper to test literal expression nodes
func testLiteralExpression(t *testing.T, expr ast.ExpressionNode, expected interface{}) {
	t.Helper()
	switch v := expected.(type) {
	case int:
		testNumberLiteral(t, expr, int64(v))
	case int64:
		testNumberLiteral(t, expr, v)
	case float64: // Although parser currently uses float64 internally for NumberLiteral
		testNumberLiteral(t, expr, v)
	case string:
		// Could be an identifier or a string literal
		ident, ok := expr.(*ast.Identifier)
		if ok {
			if ident.Value != v {
				t.Errorf("Identifier.Value not %s. got=%s", v, ident.Value)
			}
			if ident.TokenLiteral() != v {
				t.Errorf("Identifier.TokenLiteral not %s. got=%s", v, ident.TokenLiteral())
			}
		} else {
			testStringLiteral(t, expr, v)
		}
	case bool:
		// testBooleanLiteral(t, expr, v) // Add when boolean literals are supported
	default:
		t.Errorf("type of expr not handled: %T", expr)
	}
}

// Helper for number literals
func testNumberLiteral(t *testing.T, expr ast.ExpressionNode, expectedValue interface{}) {
	t.Helper()
	numLit, ok := expr.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("expr not *ast.NumberLiteral. got=%T", expr)
	}

	var expectedFloat float64
	switch val := expectedValue.(type) {
	case int64:
		expectedFloat = float64(val)
	case float64:
		expectedFloat = val
	default:
		t.Fatalf("Unexpected type for expected number value: %T", expectedValue)
	}

	if numLit.Value != expectedFloat {
		t.Errorf("NumberLiteral.Value not %f. got=%f", expectedFloat, numLit.Value)
	}

	// Check TokenLiteral (might be int or float string representation)
	expectedLiteral := fmt.Sprintf("%v", expectedValue) // Simple conversion for comparison
	if numLit.TokenLiteral() != expectedLiteral {
		// Allow for float comparison nuances if necessary, e.g., "5" vs "5.0"
		// For now, require exact match based on parser's current literal storage
		t.Errorf("NumberLiteral.TokenLiteral not %s. got=%s", expectedLiteral, numLit.TokenLiteral())
	}
}

// Helper for string literals
func testStringLiteral(t *testing.T, expr ast.ExpressionNode, expectedValue string) {
	t.Helper()
	strLit, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("expr not *ast.StringLiteral. got=%T", expr)
	}

	if strLit.Value != expectedValue {
		t.Errorf("StringLiteral.Value not %q. got=%q", expectedValue, strLit.Value)
	}

	// TokenLiteral should be the raw string *content* as seen by the lexer
	if strLit.TokenLiteral() != expectedValue {
		t.Errorf("StringLiteral.TokenLiteral not %q. got=%q", expectedValue, strLit.TokenLiteral())
	}
}

// Helper for identifier nodes
func testIdentifier(t *testing.T, expr ast.ExpressionNode, expectedValue string) {
	t.Helper()
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("expr not *ast.Identifier. got=%T (%s)", expr, expr.String())
	}

	if ident.Value != expectedValue {
		t.Errorf("Identifier.Value not %s. got=%s", expectedValue, ident.Value)
	}

	if ident.TokenLiteral() != expectedValue {
		t.Errorf("Identifier.TokenLiteral not %s. got=%s", expectedValue, ident.TokenLiteral())
	}
}
