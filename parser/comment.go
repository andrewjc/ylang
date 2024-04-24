package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseComment() *ast.Comment {
	comment := &ast.Comment{Token: p.currentToken}
	if p.currentTokenIs(TokenTypeComment) {
		comment.Text = p.currentToken.Literal
	} else if p.currentTokenIs(TokenTypeMultiLineComment) {
		comment.Text = p.currentToken.Literal
	}
	p.nextToken()
	return comment
}
