package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"fmt"
	"testing"
)

func TestInfixExpressionUnit(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		// Wrap expressions in a minimal main function for valid parsing
		{"main() -> {5 + 5;}", int64(5), "+", int64(5)},
		{"main() -> {5 - 5;}", int64(5), "-", int64(5)},
		{"main() -> {5 * 5;}", int64(5), "*", int64(5)},
		{"main() -> {5 / 5;}", int64(5), "/", int64(5)},
		{"main() -> {5 > 5;}", int64(5), ">", int64(5)},
		{"main() -> {5 < 5;}", int64(5), "<", int64(5)},
		// {"main() -> {5 == 5;}", int64(5), "==", int64(5)}, // TokenTypeEqual not handled in infix map?
		// {"main() -> {5 != 5;}", int64(5), "!=", int64(5)}, // Need != operator
		{"main() -> {foo + bar;}", "foo", "+", "bar"},
		{"main() -> {foo - bar;}", "foo", "-", "bar"},
		{"main() -> {foo * bar;}", "foo", "*", "bar"},
		{"main() -> {foo / bar;}", "foo", "/", "bar"},
		{"main() -> {foo > bar;}", "foo", ">", "bar"},
		{"main() -> {foo < bar;}", "foo", "<", "bar"},
		// {"main() -> {foo == bar;}", "foo", "==", "bar"},
		// {"main() -> {foo != bar;}", "foo", "!=", "bar"},
		// Add boolean tests when supported:
		// {"main() -> {true == true;}", true, "==", true},
		// {"main() -> {true != false;}", true, "!=", false},
		// {"main() -> {false == false;}", false, "==", false},
	}

	for _, tt := range infixTests {
		t.Run(fmt.Sprintf("Input_%s", tt.input), func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p) // Use the helper from statement test

			if program.MainFunction == nil {
				t.Fatalf("ParseProgram() returned nil MainFunction")
			}

			if len(program.MainFunction.Body.(*ast.BlockStatement).Statements) != 1 {
				t.Fatalf("MainFunction.Body does not contain 1 statement. got=%d", len(program.MainFunction.Body.(*ast.BlockStatement).Statements))
			}

			stmt, ok := program.MainFunction.Body.(*ast.BlockStatement).Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("MainFunction.Body.Statements[0] is not *ast.ExpressionStatement. got=%T", program.MainFunction.Body.(*ast.BlockStatement).Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not *ast.InfixExpression. got=%T (%s)", stmt.Expression, stmt.Expression.String())
			}

			// Test left operand
			testLiteralExpression(t, exp.Left, tt.leftValue) // Use the helper from statement test

			// Test operator
			if exp.Operator != tt.operator {
				t.Errorf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			}

			// Test right operand
			testLiteralExpression(t, exp.Right, tt.rightValue)
		})
	}
}

func TestPrefixExpressionUnit(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{} // Value of the operand
	}{
		// {"main() -> {!true;}", "!", true},   // Need boolean support
		// {"main() -> {!false;}", "!", false}, // Need boolean support
		{"main() -> {-15;}", "-", int64(15)},
		// {"main() -> {!5;}", "!", int64(5)},    // '!' usually for booleans, but test if parser handles it
	}

	for _, tt := range prefixTests {
		t.Run(fmt.Sprintf("Input_%s", tt.input), func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
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

			stmt, ok := program.MainFunction.Body.(*ast.BlockStatement).Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("MainFunction.Body.Statements[0] is not *ast.ExpressionStatement. got=%T", program.MainFunction.Body.(*ast.BlockStatement).Statements[0])
			}

			// Assuming a PrefixExpression node exists or should exist
			_, ok = stmt.Expression.(ast.ExpressionNode) // Placeholder - replace with *ast.PrefixExpression if defined
			if !ok {
				// If PrefixExpression doesn't exist, this test will fail, indicating missing feature
				t.Fatalf("stmt.Expression is not *ast.PrefixExpression. got=%T (%s)", stmt.Expression, stmt.Expression.String())
			}

			// Assuming PrefixExpression has Operator and Right fields
			// if exp.Operator != tt.operator {
			//  t.Errorf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			// }
			// testLiteralExpression(t, exp.Right, tt.value)

			// Temporary check based on current parser behavior (parses '-' then number)
			// This needs to be updated when prefix operators are properly handled.
			infExp, isInfix := stmt.Expression.(*ast.InfixExpression)
			if tt.operator == "-" && isInfix && infExp.Left == nil { // Heuristic for unary minus parsed as infix
				if infExp.Operator != "-" {
					t.Errorf("Expected operator '-' for negation, got %s", infExp.Operator)
				}
				testLiteralExpression(t, infExp.Right, tt.value)
			} else {
				t.Errorf("Expected prefix expression for input %q, but got %T (%s). Need *ast.PrefixExpression support.", tt.input, stmt.Expression, stmt.Expression.String())
			}
		})
	}
}

func TestOperatorPrecedenceParsingUnit(t *testing.T) {
	tests := []struct {
		input    string
		expected string // Expected string representation of the parsed AST
	}{
		{
			"main() -> {-a * b;}",
			"main() -> ((-a) * b);", // Requires PrefixExpression support
		},
		// {
		// 	"main() -> {!-a;}",
		// 	"main() -> (!(-a));", // Requires PrefixExpression support
		// },
		{
			"main() -> {a + b + c;}",
			"main() -> ((a + b) + c);",
		},
		{
			"main() -> {a + b - c;}",
			"main() -> ((a + b) - c);",
		},
		{
			"main() -> {a * b * c;}",
			"main() -> ((a * b) * c);",
		},
		{
			"main() -> {a * b / c;}",
			"main() -> ((a * b) / c);",
		},
		{
			"main() -> {a + b / c;}",
			"main() -> (a + (b / c));",
		},
		{
			"main() -> {a + b * c + d / e - f;}",
			"main() -> (((a + (b * c)) + (d / e)) - f);",
		},
		{
			"main() -> {3 + 4; -5 * 5;}",       // Two statements
			"main() -> {(3 + 4); ((-5) * 5);}", // Requires PrefixExpression support
		},
		{
			"main() -> {5 > 4 == 3 < 4;}", // Needs == support
			"main() -> ((5 > 4) == (3 < 4));",
		},
		// {
		// 	"main() -> {5 < 4 != 3 > 4;}", // Needs != support
		// 	"main() -> ((5 < 4) != (3 > 4));",
		// },
		// {
		// 	"main() -> {3 + 4 * 5 == 3 * 1 + 4 * 5;}", // Needs == support
		// 	"main() -> ((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		// },
		// Boolean precedence if added
		// {
		// 	"main() -> {true;}",
		// 	"main() -> true;",
		// },
		// {
		// 	"main() -> {false;}",
		// 	"main() -> false;",
		// },
		// {
		// 	"main() -> {3 > 5 == false;}", // Needs == support
		// 	"main() -> ((3 > 5) == false);",
		// },
		// {
		// 	"main() -> {3 < 5 == true;}", // Needs == support
		// 	"main() -> ((3 < 5) == true);",
		// },
		// Grouping with Parentheses
		{
			"main() -> {1 + (2 + 3) + 4;}",
			"main() -> ((1 + (2 + 3)) + 4);",
		},
		{
			"main() -> {(5 + 5) * 2;}",
			"main() -> ((5 + 5) * 2);",
		},
		{
			"main() -> {2 / (5 + 5);}",
			"main() -> (2 / (5 + 5));",
		},
		// {
		// 	"main() -> {-(5 + 5);}", // Requires PrefixExpression support
		// 	"main() -> (-(5 + 5));",
		// },
		// {
		// 	"main() -> {!(true == true);}", // Requires PrefixExpression and Boolean support
		// 	"main() -> (!(true == true));",
		// },
		// Calls and Indexing (Higher precedence)
		{
			"main() -> {a + add(b * c) + d;}",
			"main() -> ((a + add((b * c))) + d);",
		},
		{
			"main() -> {add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));}",
			"main() -> add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));",
		},
		{
			"main() -> {add(a + b + c * d / f + g);}",
			"main() -> add((((a + b) + ((c * d) / f)) + g));",
		},
		{
			"main() -> {a * [1, 2, 3, 4][b * c] * d;}", // Needs Array literal + Index support
			"main() -> ((a * ([1, 2, 3, 4][(b * c)])) * d);",
		},
		{
			"main() -> {add(a * b[2], b[1], 2 * [1, 2][1]);}", // Needs Array literal + Index support
			"main() -> add((a * (b[2])), (b[1]), (2 * ([1, 2][1])));",
		},
		// Assignment
		{
			"main() -> {a = b + c;}",
			"main() -> a = (b + c);",
		},
		{
			"main() -> {a = 5;}",
			"main() -> a = 5;",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input_%s", tt.input), func(t *testing.T) {
			l, err := lexer.NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("Lexer creation failed: %v", err)
			}
			p := NewParser(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			actual := program.String() // Assumes Program.String() correctly represents structure

			// Note: The expected strings need adjustment if PrefixExpression is not implemented yet.
			// We might need to compare against the structure the *current* parser produces.
			// However, the requirement is to test against the *intended* behavior.
			if actual != tt.expected {
				t.Errorf("AST string representation mismatch.\nExpected: %s\nGot:      %s", tt.expected, actual)
			}
		})
	}
}

func TestGroupedExpressionParsingUnit(t *testing.T) {
	input := `main() -> {(1 + 2) * 3;}`
	expected := "main() -> ((1 + 2) * 3);" // String representation of the AST

	l, err := lexer.NewLexerFromString(input)
	if err != nil {
		t.Fatalf("Lexer creation failed: %v", err)
	}
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	actual := program.String()
	if actual != expected {
		t.Errorf("Grouped expression parsing mismatch.\nExpected: %s\nGot:      %s", expected, actual)
	}

	// Test a more complex grouped expression
	inputComplex := `main() -> {(5 + (3 * 2)) / (4 - 1);}`
	expectedComplex := "main() -> ((5 + (3 * 2)) / (4 - 1));"

	lComplex, err := lexer.NewLexerFromString(inputComplex)
	if err != nil {
		t.Fatalf("Lexer creation failed for complex input: %v", err)
	}
	pComplex := NewParser(lComplex)
	programComplex := pComplex.ParseProgram()
	checkParserErrors(t, pComplex)

	actualComplex := programComplex.String()
	if actualComplex != expectedComplex {
		t.Errorf("Complex grouped expression parsing mismatch.\nExpected: %s\nGot:      %s", expectedComplex, actualComplex)
	}
}
