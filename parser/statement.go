package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case TokenTypeLet:
		return p.parseLetStatement()
	case TokenTypeReturn:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(TokenTypeIdentifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(TokenTypeAssignment) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return stmt
}
