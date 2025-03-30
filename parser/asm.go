package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseAssemblyStatement() ast.ExpressionNode {
	// eg:
	/*
		asm('
			mov src, dst
		');
	*/

	expr := &ast.AssemblyExpression{Token: p.currentToken}

	if !p.expectPeek(TokenTypeLeftParenthesis) {
		p.errors = append(p.errors, fmt.Sprintf("expected '(' after 'asm', got %s", p.peekToken.Type))
		return nil
	}

	if !p.expectPeek(TokenTypeString) {
		p.errors = append(p.errors, fmt.Sprintf("expected string literal for assembly code, got %s", p.peekToken.Type))
		return nil
	}
	expr.Code = p.parseStringLiteral().(*ast.StringLiteral)

	// Check for optional arguments
	if p.peekTokenIs(TokenTypeComma) {
		p.nextToken() // Consume ','
		expr.Args = p.parseExpressionList(TokenTypeRightParenthesis)
	} else if p.expectPeek(TokenTypeRightParenthesis) {
		expr.Args = []ast.ExpressionNode{}
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected ',' or ')' after assembly code string, got %s", p.peekToken.Type))
		return nil
	}

	return expr
}
