package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) isFunctionDefinition() bool {
	/*
		Handle checking if the current token is a function definition

		If the current token is an identifier and the next token is a left parenthesis or lambda arrow, then it is a function definition

		Examples:
		- add(a, b) { return a + b; }
		- (a, b) -> { return a + b; }
		- (a, b) -> printf("a: %d, b: %d", a, b);
		- main() -> { printf("Hello, World!"); }
		- main(argv, argc) -> { printf("Hello, World!"); }
	*/
	if p.currentTokenIs(TokenTypeIdentifier) && (p.peekTokenIs(TokenTypeLeftParenthesis) || p.peekTokenIs(TokenTypeLambdaArrow)) {
		return true
	}

	// handle anonymous methods / lambdas eg () -> { printf("Hello, World!"); } and (a) -> { printf("Hello {a}"); }
	if p.currentTokenIs(TokenTypeLeftParenthesis) && (p.peekTokenIs(TokenTypeIdentifier) || p.peekTokenIs(TokenTypeRightParenthesis)) {
		return true
	}

	return false
}

func (p *Parser) parseFunctionDefinition() ast.ExpressionNode {
	/*
			Handle parsing of function definitions
			- Function definition can be a traditional function definition or a lambda-style function definition
			- Traditional function definition: function name(parameters) { body }
			- Lambda-style function definition: (parameters) -> { body }
			- Parameters are a list of identifiers separated by commas
			- Body is a block statement
		    Examples:
			- func add(a, b) { return a + b; }
			- (a, b) -> { return a + b; }
			- (a, b) -> printf("a: %d, b: %d", a, b);
		    - main() -> { printf("Hello, World!"); }
			- main(argv, argc) -> { printf("Hello, World!"); }

	*/
	fn := &ast.FunctionDefinition{Token: p.currentToken}

	if !p.isFunctionDefinition() {
		line, pos := p.peekToken.Line, p.peekToken.Pos
		snippet := p.lexer.GetCodeFragment(line, pos, DEFAULT_LOGGING_LEAD_LINES, DEFAULT_LOGGING_FOLLOW_LINES) // Get 10 characters around the error location
		parseError := &ParserError{
			Line:         line,
			Pos:          pos,
			Message:      fmt.Sprintf("Syntax error: Expected function definition, got %s", p.peekToken.Type),
			CodeFragment: snippet,
		}
		fmt.Println(parseError)
		return nil
	}

	if p.currentTokenIs(TokenTypeIdentifier) {
		fn.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		p.nextToken() // Skip the left parenthesis
	}

	if p.currentTokenIs(TokenTypeLeftParenthesis) {
		fn.Parameters = p.parseFunctionParameters()
		p.nextToken() // skip the right parenthesis
	}

	if !p.peekTokenIs(TokenTypeLambdaArrow) && !p.peekTokenIs(TokenTypeLeftBrace) {
		line, pos := p.peekToken.Line, p.peekToken.Pos
		snippet := p.lexer.GetCodeFragment(line, pos, DEFAULT_LOGGING_LEAD_LINES, DEFAULT_LOGGING_FOLLOW_LINES)
		parseError := &ParserError{
			Line:         line,
			Pos:          pos,
			Message:      fmt.Sprintf("Syntax error: Expected '->' or '{', got %s", p.peekToken.Type),
			CodeFragment: snippet,
		}
		fmt.Println(parseError)
		return nil
	}

	if p.currentTokenIs(TokenTypeLambdaArrow) {
		p.nextToken() // Skip the lambda arrow
	}

	if !p.currentTokenIs(TokenTypeLeftBrace) {
		line, pos := p.peekToken.Line, p.peekToken.Pos
		snippet := p.lexer.GetCodeFragment(line, pos, DEFAULT_LOGGING_LEAD_LINES, DEFAULT_LOGGING_FOLLOW_LINES)
		parseError := &ParserError{
			Line:         line,
			Pos:          pos,
			Message:      fmt.Sprintf("Syntax error: Expected '{', got %s", p.peekToken.Type),
			CodeFragment: snippet,
		}
		fmt.Println(parseError)
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseFunctionParameters() []string {
	var parameters []string
	for !p.currentTokenIs(TokenTypeRightParenthesis) {
		if p.currentTokenIs(TokenTypeIdentifier) {
			parameters = append(parameters, p.currentToken.Literal)
		}
		p.nextToken()
		if p.currentTokenIs(TokenTypeComma) {
			p.nextToken()
		}
	}

	return parameters
}
