package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseArrayLiteral() ast.ExpressionNode {
	arrayLit := &ast.ArrayLiteral{Token: p.currentToken}
	arrayLit.Elements = []ast.ExpressionNode{}

	// Check for empty array: '[]'
	if p.peekTokenIs(TokenTypeRightBracket) {
		p.nextToken() // Consume ']'
		return arrayLit
	}

	p.nextToken()

	firstElem := p.parseExpression(LOWEST)
	if firstElem == nil {
		return nil
	}
	arrayLit.Elements = append(arrayLit.Elements, firstElem)

	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken() // Consume ','
		p.nextToken() // Move to the start of the next expression

		if p.currentTokenIs(TokenTypeRightBracket) {
			p.errors = append(p.errors, fmt.Sprintf("unexpected trailing comma in array literal at line %d, pos %d", p.currentToken.Line, p.currentToken.Pos))
			goto endLoop
		}

		elem := p.parseExpression(LOWEST)
		if elem == nil {
			return arrayLit
		}
		arrayLit.Elements = append(arrayLit.Elements, elem)
	}

endLoop:

	// After parsing elements and commas (or just the first element), expect ']'
	if !p.expectPeek(TokenTypeRightBracket) {
		if !p.peekTokenIs(TokenTypeEOF) {
			p.advanceToRecoveryPoint()
		}
		return arrayLit
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
