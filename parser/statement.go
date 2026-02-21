package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseStatement() ast.Statement {
	if p.currentTokenIs(TokenTypeSemicolon) {
		p.nextToken()
		return nil
	}

	switch p.currentToken.Type {
	case TokenTypeLet:
		ls := p.parseLetStatement()
		if ls == nil {
			return nil // Propagate nil on failure
		}
		return ls
	case TokenTypeReturn:
		rsNode := p.parseReturnStatement()
		if rsNode == nil {
			return nil // Propagate nil on failure
		}
		if stmt, ok := rsNode.(ast.Statement); ok {
			return stmt
		}
		p.errors = append(p.errors, fmt.Sprintf("INTERNAL ERROR: *ast.ReturnStatement does not satisfy ast.Statement interface near line %d", p.currentToken.Line+1))
		return nil
	case TokenTypeImport:
		return p.parseImportStatement()
	case TokenTypeIf:
		ifNode := p.parseIfStatement()
		if ifNode == nil {
			return nil
		}
		if stmt, ok := ifNode.(ast.Statement); ok {
			return stmt
		}
		p.errors = append(p.errors, fmt.Sprintf("INTERNAL ERROR: *ast.IfStatement does not satisfy ast.Statement interface near line %d", p.currentToken.Line+1))
		return nil
	default:
		es := p.parseExpressionStatement()
		if es == nil {
			return nil // Propagate nil on failure or empty statement
		}
		return es
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
	_ = precedence // assignment is right-associative; parse full right-side expression
	p.nextToken()
	expr.Right = p.parseExpression(LOWEST)
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

	errorsBeforeExpr := len(p.errors)
	stmt.Value = p.parseExpression(LOWEST)
	if stmt.Value == nil {
		if !p.errorsEncounteredSince(errorsBeforeExpr) {
			p.errors = append(p.errors, fmt.Sprintf("Failed to parse expression for let statement '%s' at line %d", stmt.Name.Value, p.currentToken.Line))
		}
		p.advanceToRecoveryPoint()
		return nil
	}

	// Advance past the expression value to avoid leaving cursor on the last expression token
	if !p.currentTokenIs(TokenTypeSemicolon) && !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
		p.nextToken()
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
	errorsBeforeRetExpr := len(p.errors)
	stmt.ReturnValue = p.parseExpression(LOWEST)
	if stmt.ReturnValue == nil {
		if !p.errorsEncounteredSince(errorsBeforeRetExpr) {
			p.errors = append(p.errors, fmt.Sprintf("Failed to parse return value expression at line %d", p.currentToken.Line))
		}
		p.advanceToRecoveryPoint()
		return nil
	}

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken() // advance to ';'
	}
	if p.currentTokenIs(TokenTypeSemicolon) {
		p.nextToken() // advance past ';'
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if p.currentTokenIs(TokenTypeSemicolon) {
		p.nextToken() // Consume the semicolon
		return nil    // Return a true nil, indicating no statement parsed
	}

	stmt := &ast.ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if stmt.Expression == nil {
		return nil
	}
	// BlockStatement prefix fn already advances past '}'; don't double-advance
	if _, isBlock := stmt.Expression.(*ast.BlockStatement); !isBlock {
		if !p.currentTokenIs(TokenTypeSemicolon) && !p.currentTokenIs(TokenTypeRightBrace) && !p.currentTokenIs(TokenTypeEOF) {
			p.nextToken()
		}
		if p.currentTokenIs(TokenTypeSemicolon) {
			p.nextToken() // advance past ';'
		}
	}
	return stmt
}
