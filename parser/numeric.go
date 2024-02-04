package parser

import (
	"compiler/ast"
	"strconv"
)

func (p *Parser) parseNumberLiteral() ast.ExpressionNode {
	lit := &ast.NumberLiteral{Token: p.currentToken}

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		// Handle error; could log or set an error on the parser
		return nil
	}

	lit.Value = value
	return lit
}
