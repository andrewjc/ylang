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
	// Handle named functions
	if p.currentTokenIs(TokenTypeFunction) {
		return true
	}
	// Handle anonymous functions: '(', parameter list, ')', '->'
	if p.currentTokenIs(TokenTypeLeftParenthesis) {
		pos := p.lexer.Position
		line := p.currentToken.Line

		// Attempt to parse parameters
		_ = p.parseFunctionParameters()

		// Check for '->' after parameters
		if p.peekTokenIs(TokenTypeLambdaArrow) {
			// Reset lexer position after lookahead
			p.lexer.Position = pos
			p.currentToken.Line = line
			p.nextToken()
			return true
		} else {
			// Not a function definition
			p.lexer.Position = pos
			p.currentToken.Line = line
			p.nextToken()
			return false
		}
	}
	return false
}

func (p *Parser) parseFunctionDefinition() *ast.FunctionDefinition {
	fn := &ast.FunctionDefinition{Token: p.currentToken}

	// Optional 'function' keyword
	if p.currentTokenIs(TokenTypeFunction) {
		p.nextToken() // Consume 'function' keyword
	}

	// Function name (optional)
	var fnName *ast.Identifier
	if p.currentTokenIs(TokenTypeIdentifier) && p.peekTokenIs(TokenTypeLeftParenthesis) {
		fnName = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		p.nextToken()
	} else if p.currentTokenIs(TokenTypeLeftParenthesis) {
		// Anonymous function, no name
		fnName = nil
	} else {
		// Error
		fmt.Println("Expected function name or '(' for anonymous function")
		return nil
	}

	fn.Name = fnName

	if !p.currentTokenIs(TokenTypeLeftParenthesis) {
		fmt.Println("Expected '(' after function name or 'function' keyword")
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()

	// Optional return type
	if p.peekTokenIs(TokenTypeColon) {
		p.nextToken() // skip ':'
		p.nextToken()
		if !p.currentTokenIs(TokenTypeIdentifier) {
			fmt.Println("Expected return type after ':'")
			return nil
		}
		fn.ReturnType = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	}

	if !p.expectPeek(TokenTypeLambdaArrow) {
		fmt.Println("Expected '->' after function definition")
		return nil
	}

	p.nextToken() // consume '->'

	// Parse function body
	if p.currentTokenIs(TokenTypeLeftBrace) {
		fn.Body = p.parseBlockStatement()
	} else {
		fn.Body = p.parseExpression(LOWEST)
	}

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
