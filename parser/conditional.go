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
		fmt.Println("Expected '{' for 'if' consequence")
		return nil
	}

	consequenceNode := p.parseBlockStatement()
	if consequenceNode == nil {
		return nil
	}
	ifStmt.Consequence = consequenceNode
	p.nextToken()

	if p.peekTokenIs(TokenTypeElse) {
		p.nextToken()                   // consume 'else'
		if p.peekTokenIs(TokenTypeIf) { // else if
			p.nextToken() // consume 'if'
			ifStmt.Alternative = p.parseIfStatement()
		} else if p.peekTokenIs(TokenTypeLeftBrace) { // else { ... }
			p.nextToken() // consume '{'
			altNode := p.parseBlockStatement()
			if altNode == nil {
				return nil
			}
			ifStmt.Alternative = altNode
			p.nextToken() // Consume '}'
		} else {
			p.errors = append(p.errors, fmt.Sprintf("Expected 'if' or '{' after 'else', got %s at line %d", p.peekToken.Type, p.peekToken.Line))
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
