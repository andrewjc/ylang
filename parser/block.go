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
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		} else {
			continue
		}
		// No nextToken() needed here inside the loop
	}

	if !p.currentTokenIs(TokenTypeRightBrace) {
		p.errors = append(p.errors, fmt.Sprintf("Expected '}' to close block starting line %d, got %s", block.Token.Line+1, p.currentToken.Type))
	} else {
		p.nextToken() // Consume '}'
	}

	return block
}
