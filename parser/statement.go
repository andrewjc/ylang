package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case TokenTypeLet:
		return p.parseLetStatement()
	case TokenTypeReturn:
		node := p.parseReturnStatement()
		if node == nil {
			return nil
		}
		if stmt, ok := node.(ast.Statement); ok {
			// Return is a statement, but its also an expression
			return stmt
		}
		p.errors = append(p.errors, fmt.Sprintf("Internal Error: *ast.ReturnStatement does not implement ast.Statement near line %d", node.TokenLiteral()))
		return nil
	case TokenTypeImport:
		return p.parseImportStatement()
	case TokenTypeIf:
		node := p.parseIfStatement()
		if node == nil {
			return nil
		}
		if stmt, ok := node.(ast.Statement); ok {
			return stmt
		}
		p.errors = append(p.errors, fmt.Sprintf("Internal Error: *ast.IfStatement does not implement ast.Statement near line %d", node.TokenLiteral()))
		return nil
	case TokenTypeLeftBrace:
		node := p.parseBlockStatement()
		if node == nil {
			return nil
		}
		return node.(ast.Statement)
	// TODO: Add For, While, Switch cases, or they'll be treated as expressions???
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseAssignmentExpression(left ast.ExpressionNode) ast.ExpressionNode {
	switch left.(type) {
	case *ast.Identifier, *ast.IndexExpression, *ast.MemberAccessExpression, *ast.DotOperator:
		// Assignable
	default:
		p.errors = append(p.errors, fmt.Sprintf("Invalid left-hand side in assignment near line %d: %s", p.currentToken.Line, left.String()))
		return nil
	}
	expr := &ast.AssignmentExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal, // "="
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence - 1)
	if expr.Right == nil {
		return nil
	}
	return expr
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}
	if !p.expectPeek(TokenTypeIdentifier) {
		p.advanceToRecoveryPoint()
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	if !p.expectPeek(TokenTypeAssignment) {
		p.advanceToRecoveryPoint()
		return nil
	}

	if p.currentTokenIs(TokenTypeAssignment) {
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected '=' operator after let statement identifier near line %d, got %s", p.currentToken.Line, p.currentToken.Type))
	}

	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil {
		if !p.errorsEncounteredSince(len(p.errors)) {
			p.errors = append(p.errors, fmt.Sprintf("Failed to parse expression for let statement '%s' at line %d", stmt.Name.Value, p.currentToken.Line))
		}
		p.advanceToRecoveryPoint()
		return nil
	}

	if p.currentTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.ExpressionNode {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	p.nextToken()
	if p.currentTokenIs(TokenTypeSemicolon) || p.currentTokenIs(TokenTypeRightBrace) || p.currentTokenIs(TokenTypeEOF) {
		stmt.ReturnValue = nil
		if p.currentTokenIs(TokenTypeSemicolon) {
			p.nextToken()
		}
		return stmt
	}
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if stmt.ReturnValue == nil {
		if !p.errorsEncounteredSince(len(p.errors)) {
			p.errors = append(p.errors, fmt.Sprintf("Failed to parse return value expression at line %d", p.currentToken.Line))
		}
		p.advanceToRecoveryPoint()
		return nil
	}

	if p.currentTokenIs(TokenTypeSemicolon) {
		p.nextToken()

	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression == nil {
		if p.currentTokenIs(TokenTypeSemicolon) {
			p.nextToken()
			return nil
		}
		return nil
	}
	if p.currentTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}
	return stmt
}
