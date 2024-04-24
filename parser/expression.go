package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(TokenTypeSemicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseInfixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
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

	if !p.expectPeek(TokenTypeRightParenthesis) {
		fmt.Println("Expected ')' after expression")
		return nil
	}

	return exp
}

func (p *Parser) parseDotOperator(left ast.ExpressionNode) ast.ExpressionNode {
	dotOperator := &ast.DotOperator{
		Token: p.currentToken,
		Left:  left,
	}

	p.nextToken()

	if !p.currentTokenIs(TokenTypeIdentifier) {
		return nil
	}

	dotOperator.Right = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	return dotOperator
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
