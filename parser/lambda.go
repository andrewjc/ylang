package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseLambdaExpression() ast.ExpressionNode {
	lambda := &ast.LambdaExpression{Token: p.currentToken}

	if !p.expectPeek(TokenTypeLeftParenthesis) {
		return nil
	}

	lambda.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(TokenTypeArrow) {
		return nil
	}

	p.nextToken()

	if p.currentTokenIs(TokenTypeLeftBrace) {
		fnBody := p.parseBlockStatement()
		lambda.Body = fnBody
	} else {
		lambda.Body = p.parseExpression(LOWEST)
	}

	return lambda
}
