package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseBlockStatement() ast.ExpressionNode {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	lastPosition := -1

	for !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
		if p.lexer.Position == lastPosition {
			p.errors = append(p.errors, fmt.Sprintf("Parser stuck on token %s ('%s') while parsing block at line %d, pos %d. Attempting recovery.",
				p.currentToken.Type, p.currentToken.Literal, p.currentToken.Line+1, p.currentToken.Pos))
			p.nextToken()
			lastPosition = p.lexer.Position
			continue
		}
		lastPosition = p.lexer.Position

		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		if p.errors != nil && len(p.errors) > 0 {
			p.errors = append(p.errors, fmt.Sprintf("Error parsing statement in block at line %d, pos %d", p.currentToken.Line+1, p.currentToken.Pos))
			p.advanceToRecoveryPoint()
			break
		}
	}

	if !p.currentTokenIs(TokenTypeRightBrace) {
		p.errors = append(p.errors, fmt.Sprintf("Expected '}' to close block starting line %d, but reached %s", block.Token.Line+1, p.currentToken.Type))
	} else {
		p.nextToken()
	}

	return block
}
