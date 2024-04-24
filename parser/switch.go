package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseSwitchStatement() *ast.SwitchStatement {
	stmt := &ast.SwitchStatement{Token: p.currentToken}

	if !p.expectPeek(TokenTypeLeftParenthesis) {
		fmt.Println("Expected '(' after 'switch'")
		return nil
	}

	p.nextToken()
	stmt.Expression = p.parseExpression(LOWEST)

	if !p.expectPeek(TokenTypeRightParenthesis) {
		fmt.Println("Expected ')' after switch expression")
		return nil
	}

	if !p.expectPeek(TokenTypeLeftBrace) {
		fmt.Println("Expected '{' after switch expression")
		return nil
	}

	for !p.currentTokenIs(TokenTypeRightBrace) {
		if p.currentTokenIs(TokenTypeDefault) || p.currentTokenIs(TokenTypeUndefined) {
			stmt.DefaultCase = p.parseSwitchCase()
		} else {
			switchCase := p.parseSwitchCase()
			if switchCase != nil {
				stmt.Cases = append(stmt.Cases, switchCase)
			}
		}
		p.nextToken()
	}
	if !p.expectPeek(TokenTypeRightBrace) {
		fmt.Println("Expected '}' after switch cases")
		return nil
	}
	return stmt
}

func (p *Parser) parseSwitchCase() *ast.SwitchCase {
	switchCase := &ast.SwitchCase{Token: p.currentToken}
	if p.currentTokenIs(TokenTypeCase) {
		p.nextToken()
	}
	switchCase.Expression = p.parseExpression(LOWEST)
	if !p.expectPeek(TokenTypeColon) {
		fmt.Println("Expected ':' after case expression")
		return nil
	}
	switchCase.Block = p.parseBlockStatement()
	return switchCase
}
