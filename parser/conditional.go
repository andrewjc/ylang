package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}

	if !p.expectPeek(TokenTypeLeftParenthesis) {
		return nil // Expected '(' after 'if'
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(TokenTypeRightParenthesis) {
		return nil // Expected ')' after condition
	}

	if !p.expectPeek(TokenTypeLeftBrace) {
		return nil // Expected '{' after ')'
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(TokenTypeElse) {
		p.nextToken()

		if !p.expectPeek(TokenTypeLeftBrace) {
			return nil // Expected '{' after 'else'
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseLambdaIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{Token: p.currentToken}

	stmt.Condition = p.parseLambdaExpression()

	if !p.expectPeek(TokenTypeLeftBrace) {
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
