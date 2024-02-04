package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}

	if err := p.expectPeek(TokenTypeLeftParenthesis); err != nil {
		fmt.Println(err)
		return nil // Expected '(' after 'if'
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if err := p.expectPeek(TokenTypeRightParenthesis); err != nil {
		fmt.Println(err)
		return nil // Expected '(' after 'if'
	}

	if err := p.expectPeek(TokenTypeLeftBrace); err != nil {
		fmt.Println(err)
		return nil // Expected '{' after ')'
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(TokenTypeElse) {
		p.nextToken()

		if err := p.expectPeek(TokenTypeLeftBrace); err != nil {
			fmt.Println(err)
			return nil // Expected '{' after ')'
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseLambdaIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}

	stmt.Condition = p.parseLambdaExpression()

	if err := p.expectPeek(TokenTypeLeftBrace); err != nil {
		fmt.Println(err)
		return nil // Expected '{' after lambda expression
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(TokenTypeElse) {
		p.nextToken()

		stmt.Alternative = p.parseLambdaElseBlock()
	}

	return stmt
}
func (p *Parser) parseLambdaExpression() ast.ExpressionNode {
	// Parse the lambda expression, which is an expression node
	// This will include parsing the parameters and body of the lambda

	return nil
}

func (p *Parser) parseLambdaElseBlock() *ast.BlockStatement {
	// Similar to parsing the main lambda block but for the 'else' part

	return nil

}
