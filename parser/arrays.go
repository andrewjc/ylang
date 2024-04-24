package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseArrayLiteral() ast.ExpressionNode {
	arrayLit := &ast.ArrayLiteral{Token: p.currentToken}
	arrayLit.Elements = []ast.ExpressionNode{}

	// Skip the opening bracket
	p.nextToken()

	// Parse elements until we reach a closing bracket
	for !p.peekTokenIs(TokenTypeRightBracket) {
		elem := p.parseExpression(LOWEST)
		if elem != nil {
			arrayLit.Elements = append(arrayLit.Elements, elem)
		}

		// Move to the next token, which should be a comma or a closing bracket
		p.nextToken()
		if p.currentTokenIs(TokenTypeComma) {
			p.nextToken() // Skip comma
		}
	}

	// Ensure we have a closing bracket
	if !p.expectPeek(TokenTypeRightBracket) {
		fmt.Println("Expected closing bracket")
		return nil
	}
	return arrayLit
}

func (p *Parser) parseIndexExpression(left ast.ExpressionNode) ast.ExpressionNode {
	exp := &ast.IndexExpression{Token: p.currentToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(TokenTypeRightBracket) {
		return nil
	}

	return exp
}
