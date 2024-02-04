package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	// Initial expression parsing based on the current token
	var leftExp ast.ExpressionNode

	switch p.currentToken.Type {
	case TokenTypeNumber:
		leftExp = p.parseNumberLiteral()
	case TokenTypeString:
		leftExp = p.parseStringLiteral()
	case TokenTypeIdentifier:
		leftExp = p.parseIdentifier()
	case TokenTypeLeftParenthesis:
		leftExp = p.parseParenthesisExpression()
	case TokenTypeLeftBracket:
		leftExp = p.parseArrayLiteral()
	case TokenTypeQuestionMark:
		leftExp = p.parseTraditionalTernaryExpression(leftExp)
	case TokenTypeLambdaArrow:
		if p.isFunctionDefinition() {
			leftExp = p.parseFunctionDefinition()
		} else if p.isTernary() {
			leftExp = p.parseTernaryExpression(leftExp)
		} else {
			leftExp = p.parseLambdaExpression()
		}
	case TokenTypeIf:
		leftExp = p.parseInlineIfElseTernaryExpression(leftExp)
	case TokenTypeLet:
		leftExp = p.parseVariableDeclaration()
	case TokenTypeAssignment:
		leftExp = p.parseAssignmentStatement()
	case TokenTypeLeftBrace:
		leftExp = p.parseBlockStatement()
	}

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
