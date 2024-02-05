package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	if p.currentTokenIs(TokenTypeLeftBrace) {
		p.nextToken()
	}

	for !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseParenthesisExpression() ast.ExpressionNode {
	/*
		Handle special case where a lambda expression is provided as the object being assigned to a variable
	*/

	if p.isFunctionDefinition() {
		return p.parseFunctionDefinition()
	}

	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if err := p.expectPeek(TokenTypeRightParenthesis); err != nil {
		fmt.Println(err)
		return nil
	}

	return exp
}
