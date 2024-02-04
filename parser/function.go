package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseFunctionDefinition() ast.ExpressionNode {
	fn := &ast.FunctionDefinition{Token: p.currentToken}

	// Skip the left parenthesis
	p.nextToken()

	// Parse the parameters
	fn.Parameters = p.parseFunctionParameters()

	// Skip the right parenthesis
	p.nextToken()

	// Parse the function body
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
