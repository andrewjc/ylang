package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	if p.currentTokenIs(TokenTypeEOF) {
		p.errors = append(p.errors, fmt.Sprintf("Unexpected EOF while parsing expression at line %d", p.currentToken.Line))
		return nil
	}

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
	startToken := p.currentToken // Keep track of the opening '(' line/pos for errors

	// Check for empty parameters lambda: () ->
	if p.peekTokenIs(TokenTypeRightParenthesis) && p.peekToken2Is(TokenTypeLambdaArrow) {
		p.nextToken()                   // Consume '('
		p.nextToken()                   // Consume ')' - currentToken is now ')'
		lambdaArrowToken := p.peekToken // The '->' token
		p.nextToken()                   // Consume '->' - currentToken is now '->'

		lambda := &ast.LambdaExpression{Token: lambdaArrowToken} // Use '->' as the token for the node
		lambda.Parameters = []*ast.Identifier{}                  // Empty params

		// Parse Body after '->'
		p.nextToken() // Move to the start of the body expression/block
		lambda.Body = p.parseLambdaBody()
		if lambda.Body == nil {
			return nil // Error handled in parseLambdaBody
		}
		return lambda
	}

	// Check for non-empty parameters lambda: (ident, ...) ->
	isLambda, params := p.probeIsLambdaParameters()
	if isLambda {
		p.nextToken()                                          // Consume '('
		lambda := &ast.LambdaExpression{Token: p.currentToken} // Use '(' as token for now
		lambda.Parameters = params

		// Consume tokens up to and including ')'
		for !p.currentTokenIs(TokenTypeRightParenthesis) && !p.currentTokenIs(TokenTypeEOF) {
			p.nextToken()
		}
		if !p.currentTokenIs(TokenTypeRightParenthesis) {
			p.errors = append(p.errors, fmt.Sprintf("Expected ')' after lambda parameters starting line %d, got %s", startToken.Line+1, p.currentToken.Type))
			return nil // Or try recovery?
		}

		// Expect '->'
		if !p.expectPeek(TokenTypeLambdaArrow) { // Consumes ')' and checks/consumes '->'
			p.errors = append(p.errors, fmt.Sprintf("Expected '->' after lambda parameter list ending line %d, got %s", p.currentToken.Line+1, p.peekToken.Type))
			return nil // Or try recovery?
		}
		lambda.Token = p.currentToken // Update token to be '->'

		// Parse Body after '->'
		p.nextToken() // Move to the start of the body
		lambda.Body = p.parseLambdaBody()
		if lambda.Body == nil {
			return nil // Error handled in parseLambdaBody
		}
		return lambda
	}

	// If it doesn't look like a lambda, parse as a grouped expression
	p.nextToken() // Consume '('
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		// Error should have been recorded by parseExpression
		// If we are stuck at ')', consume it to potentially recover.
		if p.currentTokenIs(TokenTypeRightParenthesis) {
			p.nextToken()
		} else {
			// p.advanceToRecoveryPoint() // Maybe needed if parseExpression failed badly
		}
		return nil
	}

	if !p.expectPeek(TokenTypeRightParenthesis) {
		return nil
	}

	return expr
}

func (p *Parser) probeIsLambdaParameters() (bool, []*ast.Identifier) {
	if !(p.currentTokenIs(TokenTypeLeftParenthesis)) {
		return false, nil // Should be called when current is '('
	}

	// Handle empty params () -> checked by caller already
	if p.peekTokenIs(TokenTypeRightParenthesis) && p.peekToken2Is(TokenTypeLambdaArrow) {
		return false, nil // Handled by caller
	}

	// Handle single param (ident) ->
	if p.peekTokenIs(TokenTypeIdentifier) && p.peekToken2Is(TokenTypeRightParenthesis) && p.peekToken3Is(TokenTypeLambdaArrow) {
		params := []*ast.Identifier{{Token: p.peekToken, Value: p.peekToken.Literal}}
		return true, params
	}

	// Handle multiple params (ident, ident, ...) ) ->
	if !p.peekTokenIs(TokenTypeIdentifier) {
		return false, nil // First item after ( must be identifier for non-empty list
	}

	var params []*ast.Identifier
	params = append(params, &ast.Identifier{Token: p.peekToken, Value: p.peekToken.Literal})
	idx := 2 // Start peeking at index 2 (potential comma or ')')

	for {
		pk := p.peekTokenAtIndex(idx)
		if pk.Type == TokenTypeComma {
			idx++ // Move past comma
			pkNext := p.peekTokenAtIndex(idx)
			if pkNext.Type == TokenTypeIdentifier {
				params = append(params, &ast.Identifier{Token: pkNext, Value: pkNext.Literal})
				idx++    // Move past identifier
				continue // Look for next comma or ')'
			} else {
				return false, nil // Expected identifier after comma
			}
		} else if pk.Type == TokenTypeRightParenthesis {
			// Found closing parenthesis. Now check if the *next* token is '->'
			pkArrow := p.peekTokenAtIndex(idx + 1)
			if pkArrow.Type == TokenTypeLambdaArrow {
				return true, params // Looks like a lambda!
			} else {
				return false, nil // Found ')' but not followed by '->'
			}
		} else {
			// Unexpected token in parameter list probe
			return false, nil
		}
	}
}

func (p *Parser) parseLambdaBody() ast.ExpressionNode {
	if p.currentTokenIs(TokenTypeLeftBrace) {
		// Parse block statement. parseBlockStatement should consume '{' and '}'
		bodyNode := p.parseBlockStatement() // Re-use block parser
		if bodyNode == nil {
			p.errors = append(p.errors, fmt.Sprintf("Failed to parse block body for lambda at line %d", p.currentToken.Line+1))
			return nil
		}
		return bodyNode
	} else {
		// Parse single expression body
		bodyExpr := p.parseExpression(LOWEST)
		if bodyExpr == nil {
			p.errors = append(p.errors, fmt.Sprintf("Failed to parse expression body for lambda at line %d", p.currentToken.Line+1))
			return nil
		}
		return bodyExpr
	}
}

func (p *Parser) expressionToParameters(expr ast.ExpressionNode) ([]*ast.Identifier, bool) {

	var params []*ast.Identifier

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

func (p *Parser) parseMemberAccessExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expr := &ast.MemberAccessExpression{
		Token: p.currentToken,
		Left:  left,
	}

	if !p.expectPeek(TokenTypeIdentifier) {
		p.errors = append(p.errors, fmt.Sprintf("Expected identifier after '.', got %s at line %d", p.peekToken.Type, p.peekToken.Line))
		return nil
	}

	expr.Member = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	return expr
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
