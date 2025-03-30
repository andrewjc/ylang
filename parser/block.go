package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseBlockStatement() ast.ExpressionNode {
	if !p.currentTokenIs(TokenTypeLeftBrace) {
		p.errors = append(p.errors, fmt.Sprintf("Internal parser error: parseBlockStatement called without '{' token at line %d", p.currentToken.Line))
		return nil
	}

	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	if p.currentTokenIs(TokenTypeLeftBrace) {
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected '{' opening block statement near line %d, got %s", p.currentToken.Line, p.currentToken.Type))
	}

	for !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
		parseStartPos := p.lexer.Position
		parseStartToken := p.currentToken
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)

			if p.lexer.Position == parseStartPos && p.currentToken.Type == parseStartToken.Type && stmt.TokenLiteral() != "" {
				p.errors = append(p.errors, fmt.Sprintf("[DEBUG] parseStatement for %s didn't advance parser near line %d. Forcing advance.", stmt.TokenLiteral(), p.currentToken.Line))
				p.nextToken()
			}
		} else {
			if p.lexer.Position == parseStartPos && p.currentToken.Type == parseStartToken.Type {
				if !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
					p.errors = append(p.errors, fmt.Sprintf("[DEBUG] parseStatement returned nil and didn't advance near line %d ('%s'). Forcing advance.", p.currentToken.Line, p.currentToken.Literal))
					p.nextToken()
				}
			}
		}
	}

	if p.currentTokenIs(TokenTypeRightBrace) {
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected '}' after let statement near line %d, got %s", p.currentToken.Line, p.currentToken.Type))
	}

	return block
}
