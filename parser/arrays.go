package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseArrayLiteral() ast.ExpressionNode {
	arrayLit := &ast.ArrayLiteral{Token: p.currentToken} // Current token is '['
	arrayLit.Elements = []ast.ExpressionNode{}

	// Check for empty array: '[]'
	if p.peekTokenIs(TokenTypeRightBracket) {
		p.nextToken() // Consume ']'
		return arrayLit
	}

	// Consume '[' and move to the first element (or comma/']')
	p.nextToken()

	firstElem := p.parseExpression(LOWEST)
	if firstElem != nil { // Check if parsing succeeded
		arrayLit.Elements = append(arrayLit.Elements, firstElem)
	}

	// Parse subsequent elements (preceded by a comma)
	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken() // Consume ','
		p.nextToken() // Move to the start of the next expression

		// Check for trailing comma before closing bracket, e.g., [1, 2,]
		if p.currentTokenIs(TokenTypeRightBracket) {
			p.errors = append(p.errors, fmt.Sprintf("unexpected trailing comma in array literal at line %d, pos %d", p.currentToken.Line, p.currentToken.Pos))
			break // Exit loop, let expectPeek handle ']'
		}

		elem := p.parseExpression(LOWEST)
		if elem != nil {
			arrayLit.Elements = append(arrayLit.Elements, elem)
		}
	}

	if p.expectPeek(TokenTypeRightBracket) {
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("expected ']' at line %d, got %s", p.currentToken.Line, p.currentToken.Type))
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
