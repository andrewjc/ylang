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
	if leftExp == nil {
		return nil
	}

	for !p.peekTokenIs(TokenTypeSemicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
		if leftExp == nil {
			return nil
		}
		p.nextToken()
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
	if expression.Right == nil {
		return nil
	}

	return expression
}

func (p *Parser) parseParenthesisExpression() ast.ExpressionNode {
	/*
	   Handle special case where a lambda expression is provided as the object being assigned to a variable
	*/

	startToken := p.currentToken // '('
	p.nextToken()                // consume '('

	if p.currentTokenIs(TokenTypeRightParenthesis) && p.peekTokenIs(TokenTypeLambdaArrow) {
		p.nextToken() // consume ')'
		lambdaArrowToken := p.peekToken
		p.nextToken() // consume '->'

		lambda := &ast.LambdaExpression{Token: lambdaArrowToken}
		lambda.Parameters = []*ast.Identifier{}

		// Parse Body (Block or Expression)
		p.nextToken() // Move to start of body
		if p.currentTokenIs(TokenTypeLeftBrace) {
			bodyNode := p.parseBlockStatement()
			if bodyNode == nil {
				return nil
			}
			lambda.Body = bodyNode
			p.nextToken() // Consume '}'
		} else {
			lambda.Body = p.parseExpression(LOWEST)
			if lambda.Body == nil {
				return nil
			}
			if p.currentTokenIs(TokenTypeSemicolon) {
				p.nextToken()
			}
		}
		return lambda
	}

	initialPos := p.lexer.Position
	initialCurrent := p.currentToken
	initialPeek := p.peekToken
	initialErrorCount := len(p.errors)

	var params []*ast.Identifier
	isPotentialParamList := true
	paramParseEndedAtRightParen := false

	// Try parsing identifier list
	if p.currentTokenIs(TokenTypeIdentifier) {
		params = append(params, &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal})
		p.nextToken()
		for p.currentTokenIs(TokenTypeComma) {
			p.nextToken() // consume ','
			if !p.currentTokenIs(TokenTypeIdentifier) {
				isPotentialParamList = false
				break
			}
			params = append(params, &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal})
			p.nextToken() // consume identifier
		}

		if isPotentialParamList && p.currentTokenIs(TokenTypeRightParenthesis) {
			paramParseEndedAtRightParen = true
		} else {
			isPotentialParamList = false
		}
	} else {
		isPotentialParamList = false
	}

	if isPotentialParamList && paramParseEndedAtRightParen && p.peekTokenIs(TokenTypeLambdaArrow) {
		p.nextToken() // consume ')'
		lambdaArrowToken := p.peekToken
		p.nextToken() // consume '->'

		lambda := &ast.LambdaExpression{Token: lambdaArrowToken}
		lambda.Parameters = params

		if p.currentTokenIs(TokenTypeLeftBrace) {
			bodyNode := p.parseBlockStatement()
			if bodyNode == nil {
				return nil
			}
			lambda.Body = bodyNode
		} else {
			lambda.Body = p.parseExpression(LOWEST)
			if lambda.Body == nil {
				return nil
			}
			if p.currentTokenIs(TokenTypeSemicolon) {
				p.nextToken()
			}
		}
		return lambda
	}

	p.errors = p.errors[:initialErrorCount]
	p.currentToken = initialCurrent
	p.peekToken = initialPeek

	p.lexer.Position = initialPos

	if initialPos > 0 {
		p.lexer.ReadChar()
	}
	p.nextToken() // Re-read current token after '('
	p.nextToken() // Re-read peek token

	if p.currentToken.Type != initialCurrent.Type || p.currentToken.Literal != initialCurrent.Literal {
		p.errors = append(p.errors, fmt.Sprintf("[CRITICAL] Parser state reset failed! Cannot parse grouped expression reliably near line %d", startToken.Line))
		for !p.currentTokenIs(TokenTypeRightParenthesis) && !p.currentTokenIs(TokenTypeEOF) {
			p.nextToken()
		}
		if p.currentTokenIs(TokenTypeRightParenthesis) {
			p.nextToken()
		}
		return nil
	}

	// Parse the inner expression
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		p.advanceToRecoveryPoint()
		return nil
	}

	if !p.currentTokenIs(TokenTypeRightParenthesis) {
		if !p.expectPeek(TokenTypeRightParenthesis) { // Checks peek and consumes ')' if correct
			p.errors = append(p.errors, fmt.Sprintf("Expected ')' after grouped expression starting line %d, got %s", startToken.Line, p.peekToken.Type)) // Adjusted error message
			p.advanceToRecoveryPoint()
			return nil
		}
	} else {
		p.nextToken()
	}

	return expr
}

func (p *Parser) expressionToParameters(expr ast.ExpressionNode) ([]*ast.Identifier, bool) {
	params := []*ast.Identifier{}

	if ident, ok := expr.(*ast.Identifier); ok {
		params = append(params, ident)
		return params, true
	}

	current := expr
	for {
		infix, ok := current.(*ast.InfixExpression)
		if ok && infix.Operator == "," {
			rightIdent, rightOk := infix.Right.(*ast.Identifier)
			if !rightOk {
				return nil, false
			}
			params = append([]*ast.Identifier{rightIdent}, params...)
			current = infix.Left
		} else {
			lastIdent, lastOk := current.(*ast.Identifier)
			if !lastOk {
				return nil, false
			}
			params = append([]*ast.Identifier{lastIdent}, params...)
			break
		}
	}
	return params, true
}

func (p *Parser) parseDotOperator(left ast.ExpressionNode) ast.ExpressionNode {
	dotOperator := &ast.DotOperator{
		Token: p.currentToken,
		Left:  left,
	}

	if !p.expectPeek(TokenTypeIdentifier) {
		p.errors = append(p.errors, fmt.Sprintf("Expected identifier after '.', got %s at line %d", p.peekToken.Type, p.peekToken.Line))
		return nil
	}

	dotOperator.Right = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	return dotOperator
}

func (p *Parser) noPrefixParseFnError(t TokenType) {
	// Check for common errors like misplaced operators
	switch t {
	case TokenTypePlus, TokenTypeMinus, TokenTypeMultiply, TokenTypeDivide, TokenTypeAssignment, TokenTypeEqual, TokenTypeLessThan, TokenTypeGreaterThan, TokenTypeComma, TokenTypeColon, TokenTypeSemicolon, TokenTypeRightParenthesis, TokenTypeRightBrace, TokenTypeRightBracket:
		msg := fmt.Sprintf("Operator '%s' cannot start an expression at line %d, position %d", p.currentToken.Literal, p.currentToken.Line, p.currentToken.Pos)
		p.errors = append(p.errors, msg)
	default:
		msg := fmt.Sprintf("Syntax error: Unexpected token '%s' (%s) cannot start an expression at line %d, position %d", p.currentToken.Literal, t, p.currentToken.Line, p.currentToken.Pos)
		p.errors = append(p.errors, msg)
	}
}
