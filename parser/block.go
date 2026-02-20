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

	for !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
		lastTokenLine := p.currentToken.Line
		lastTokenPos := p.currentToken.Pos

		errorsBefore := len(p.errors)
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		if len(p.errors) > errorsBefore {
			// An error occurred; recover and continue parsing the block
			p.advanceToRecoveryPoint()
			if p.currentTokenIs(TokenTypeRightBrace) || p.currentTokenIs(TokenTypeEOF) {
				break
			}
			// Consume a semicolon recovery point so the next statement can start cleanly
			if p.currentTokenIs(TokenTypeSemicolon) {
				p.nextToken()
			}
			if p.currentTokenIs(TokenTypeRightBrace) || p.currentTokenIs(TokenTypeEOF) {
				break
			}
		} else if p.currentToken.Line == lastTokenLine && p.currentToken.Pos == lastTokenPos {
			// No progress made by parseStatement - force advance
			p.errors = append(p.errors, fmt.Sprintf("Parser stuck on token %s ('%s') while parsing block at line %d, pos %d. Attempting recovery.",
				p.currentToken.Type, p.currentToken.Literal, p.currentToken.Line+1, p.currentToken.Pos))
			p.nextToken()
		}
	}

	if !p.currentTokenIs(TokenTypeRightBrace) {
		p.errors = append(p.errors, fmt.Sprintf("Expected '}' to close block starting line %d, but reached %s", block.Token.Line+1, p.currentToken.Type))
	} else {
		p.nextToken()
	}

	return block
}
