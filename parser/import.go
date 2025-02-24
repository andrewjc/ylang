package parser

import "compiler/ast"
import . "compiler/lexer"

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.currentToken}

	if !p.expectPeek(TokenTypeString) {
		p.errors = append(p.errors, "expected string after 'import'")
		return nil
	}
	stmt.Path = p.currentToken.Literal

	// optionally skip semicolon
	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}
	return stmt
}
