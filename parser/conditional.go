package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseIfStatement() ast.ExpressionNode {
	ifStmt := &ast.IfStatement{Token: p.currentToken}

	if !p.expectPeek(TokenTypeLeftParenthesis) {
		fmt.Println("Expected left parenthesis")
		return nil
	}

	p.nextToken()
	ifStmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(TokenTypeRightParenthesis) {
		fmt.Println("Expected right parenthesis")
		return nil
	}

	if !p.expectPeek(TokenTypeLeftBrace) {
		fmt.Println("Expected left brace")
		return nil
	}

	ifStmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(TokenTypeElse) {
		p.nextToken()

		if p.peekTokenIs(TokenTypeIf) {
			p.nextToken()
			ifStmt.Alternative = p.parseIfStatement().(*ast.IfStatement)
		} else if p.peekTokenIs(TokenTypeLeftBrace) {
			p.nextToken()
			ifStmt.Alternative = p.parseBlockStatement().(*ast.BlockStatement)
		} else {
			return nil
		}
	}

	return ifStmt
}

func (p *Parser) parseTraditionalTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	ternaryExp := &ast.TraditionalTernaryExpression{
		Token:     p.currentToken,
		Condition: condition,
	}

	if !p.expectPeek(TokenTypeQuestionMark) {
		return nil
	}

	p.nextToken()
	ternaryExp.TrueExpr = p.parseExpression(TERNARY)

	if !p.expectPeek(TokenTypeSemicolon) {
		return nil
	}

	p.nextToken()
	ternaryExp.FalseExpr = p.parseExpression(TERNARY)

	return ternaryExp
}

func (p *Parser) parseLambdaStyleTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	ternaryExp := &ast.LambdaStyleTernaryExpression{
		Token:     p.currentToken,
		Condition: condition,
	}

	if !p.expectPeek(TokenTypeLambdaArrow) {
		fmt.Println("Expected lambda arrow")
		return nil
	}

	p.nextToken()
	ternaryExp.TrueExpr = p.parseExpression(TERNARY)

	if !p.expectPeek(TokenTypeColon) {
		fmt.Println("Expected colon")
		return nil
	}

	p.nextToken()
	ternaryExp.FalseExpr = p.parseExpression(TERNARY)

	return ternaryExp
}

func (p *Parser) parseInlineIfElseTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	ternaryExp := &ast.InlineIfElseTernaryExpression{
		Token:     p.currentToken,
		Condition: condition,
	}

	if !p.expectPeek(TokenTypeIf) {
		fmt.Println("Expected 'if'")
		return nil
	}

	p.nextToken()
	ternaryExp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(TokenTypeThen) {
		fmt.Println("Expected 'then'")
		return nil
	}

	p.nextToken()
	ternaryExp.TrueExpr = p.parseExpression(TERNARY)

	if !p.expectPeek(TokenTypeElse) {
		fmt.Println("Expected 'else'")
		return nil
	}

	p.nextToken()
	ternaryExp.FalseExpr = p.parseExpression(TERNARY)

	return ternaryExp
}
