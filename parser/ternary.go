package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) isTernary() bool {
	return p.currentTokenIs(TokenTypeQuestionMark)
}

func (p *Parser) parseTernaryExpression(exp ast.ExpressionNode) ast.ExpressionNode {
	ternary := &ast.TraditionalTernaryExpression{
		Token:     LangToken{},
		Condition: nil,
		TrueExpr:  nil,
		FalseExpr: nil,
	}
	return ternary
}

func (p *Parser) parseTraditionalTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	expression := &ast.TraditionalTernaryExpression{
		Token:     p.currentToken,
		Condition: condition,
	}

	p.nextToken() // Skip '?'
	expression.TrueExpr = p.parseExpression(LOWEST)

	if err := p.expectPeek(TokenTypeColon); err != nil {
		fmt.Println(err)
		return nil // Error handling; expected a colon
	}

	p.nextToken()
	expression.FalseExpr = p.parseExpression(LOWEST)

	return expression
}

// Lambda Style Ternary
func (p *Parser) parseLambdaStyleTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	// Handle parsing of lambda-style ternary expressions

	expression := &ast.LambdaStyleTernaryExpression{
		Token:     p.currentToken,
		Condition: condition,
	}

	p.nextToken() // Skip '->'
	expression.TrueExpr = p.parseExpression(LOWEST)

	if err := p.expectPeek(TokenTypeColon); err != nil {
		fmt.Println(err)
		return nil // Error handling; expected a colon
	}

	p.nextToken()
	expression.FalseExpr = p.parseExpression(LOWEST)

	return expression
}

// Inline If-Else Ternary
func (p *Parser) parseInlineIfElseTernaryExpression(condition ast.ExpressionNode) ast.ExpressionNode {
	// Start by expecting 'if', then parse the condition, 'then' expression, and 'else' expression

	expression := &ast.InlineIfElseTernaryExpression{
		Token:     p.currentToken,
		Condition: condition,
	}

	p.nextToken() // Skip 'if'

	//expression.Condition = p.parseExpression(LOWEST)

	if err := p.expectPeek(TokenTypeThen); err != nil {
		fmt.Println(err)
		return nil // Error handling; expected 'then'
	}

	p.nextToken()
	expression.TrueExpr = p.parseExpression(LOWEST)

	if err := p.expectPeek(TokenTypeElse); err != nil {
		fmt.Println(err)
		return nil // Error handling; expected 'else'
	}

	p.nextToken()
	expression.FalseExpr = p.parseExpression(LOWEST)

	return expression
}
