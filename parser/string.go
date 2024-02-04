package parser

import "compiler/ast"

func (p *Parser) parseStringLiteral() ast.ExpressionNode {
	lit := &ast.StringLiteral{Token: p.currentToken}

	lit.Value = p.currentToken.Literal
	return lit
}
