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

func (p *Parser) parseFunctionDefinition() *ast.FunctionDefinition {
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

	fn.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(TokenTypeLeftParenthesis) {
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	if p.peekTokenIs(TokenTypeLambdaArrow) {
		p.nextToken()
	} else {
		if !p.expectPeek(TokenTypeLeftBrace) {
			return nil
		}
	}

	fnBody := p.parseBlockStatement()
	fn.Body = fnBody

	return fn
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(TokenTypeRightParenthesis) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(TokenTypeRightParenthesis) {
		return nil
	}

	return identifiers
}
