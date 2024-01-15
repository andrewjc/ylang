package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return stmt
}
