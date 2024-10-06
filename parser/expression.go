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

	p.nextToken() // consume '('

	if p.peekTokenIs(TokenTypeRightParenthesis) {
		p.nextToken() // consume ')'
		if p.peekTokenIs(TokenTypeLambdaArrow) {
			p.nextToken() // consume '->'
			lambda := &ast.LambdaExpression{Token: p.currentToken}
			lambda.Parameters = []*ast.Identifier{}
			p.nextToken()
			if p.currentTokenIs(TokenTypeLeftBrace) {
				lambda.Body = p.parseBlockStatement()
			} else {
				lambda.Body = p.parseExpression(LOWEST)
			}
			return lambda
		} else {
			// Empty parentheses expression
			return nil
		}
	}

	// Save current state to backtrack if needed
	saveCurrentToken := p.currentToken
	savePeekToken := p.peekToken

	parameters := []*ast.Identifier{}
	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	parameters = append(parameters, ident)

	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken() // consume ','
		p.nextToken()
		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		parameters = append(parameters, ident)
	}

	if !p.expectPeek(TokenTypeRightParenthesis) {
		// Not a lambda, restore tokens and parse as expression
		p.currentToken = saveCurrentToken
		p.peekToken = savePeekToken
		p.nextToken() // consume '('
		expr := p.parseExpression(LOWEST)
		if !p.expectPeek(TokenTypeRightParenthesis) {
			return nil
		}
		return expr
	}

	if p.peekTokenIs(TokenTypeLambdaArrow) {
		p.nextToken() // consume '->'
		lambda := &ast.LambdaExpression{Token: p.currentToken}
		lambda.Parameters = parameters
		p.nextToken()
		if p.currentTokenIs(TokenTypeLeftBrace) {
			lambda.Body = p.parseBlockStatement()
		} else {
			lambda.Body = p.parseExpression(LOWEST)
		}
		return lambda
	} else {
		// Not a lambda, restore tokens and parse as expression
		p.currentToken = saveCurrentToken
		p.peekToken = savePeekToken
		p.nextToken() // consume '('
		expr := p.parseExpression(LOWEST)
		if !p.expectPeek(TokenTypeRightParenthesis) {
			return nil
		}
		return expr
	}
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
	msg := fmt.Sprintf("no prefix parse function for %s found at line %d, position %d", t, p.currentToken.Line, p.currentToken.Pos)
	p.errors = append(p.errors, msg)
}
