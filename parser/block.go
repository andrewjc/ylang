package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseBlockStatement() ast.ExpressionNode {
	if p.peekTokenIs(TokenTypeLeftBrace) {
		p.nextToken()
		p.nextToken()
	}

	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}
