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
	// After parseBlockStatement(), current is the token AFTER the closing '}'

	if p.currentTokenIs(TokenTypeElse) {
		if p.peekTokenIs(TokenTypeIf) { // else if
			p.nextToken()                       // advance to 'if'
			ifStmt.Alternative = p.parseIfStatement()
		} else if p.peekTokenIs(TokenTypeLeftBrace) { // else { ... }
			p.nextToken() // advance to '{'
			altNode := p.parseBlockStatement()
			if altNode == nil {
				return nil
			}
			ifStmt.Alternative = altNode
		} else {
			p.errors = append(p.errors, fmt.Sprintf("Expected 'if' or '{' after 'else', got %s at line %d", p.peekToken.Type, p.peekToken.Line))
			return nil
		}
	}

	return ifStmt
}

func (p *Parser) parseTraditionalTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	ternaryExp := &ast.TraditionalTernaryExpression{
		Token:     p.currentToken, // '?' token (already current when called as infix fn)
		Condition: condition,
	}

	p.nextToken() // advance past '?'
	ternaryExp.TrueExpr = p.parseExpression(TERNARY)

	if !p.expectPeek(TokenTypeColon) {
		return nil
	}

	p.nextToken()
	ternaryExp.FalseExpr = p.parseExpression(TERNARY)

	return ternaryExp
}

func (p *Parser) parseLambdaStyleTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	ternaryExp := &ast.LambdaStyleTernaryExpression{
		Token:     p.currentToken, // '->' token (already current when called as infix fn)
		Condition: condition,
	}

	p.nextToken() // advance past '->'
	ternaryExp.TrueExpr = p.parseExpression(TERNARY)

	if !p.expectPeek(TokenTypeColon) {
		fmt.Println("Expected colon")
		return nil
	}

	p.nextToken()
	ternaryExp.FalseExpr = p.parseExpression(TERNARY)

	return ternaryExp
}

func (p *Parser) parseInlineIfElseTernaryExpression(trueExpr ast.ExpressionNode) ast.ExpressionNode {
	ternaryExp := &ast.InlineIfElseTernaryExpression{
		Token:    p.currentToken, // 'if' token (already current when called as infix fn)
		TrueExpr: trueExpr,       // left-hand expr is what's returned when condition is true
	}

	p.nextToken() // advance past 'if' to the condition
	ternaryExp.Condition = p.parseExpression(TERNARY)

	if !p.expectPeek(TokenTypeElse) {
		fmt.Println("Expected 'else'")
		return nil
	}

	p.nextToken()
	ternaryExp.FalseExpr = p.parseExpression(TERNARY)

	return ternaryExp
}
