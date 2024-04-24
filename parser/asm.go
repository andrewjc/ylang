package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseAssemblyStatement() *ast.AssemblyStatement {
	stmt := &ast.AssemblyStatement{Token: p.currentToken}
	if !p.expectPeek(TokenTypeString) {
		fmt.Println("Expected string after 'assembly'")
		return nil
	}
	stmt.Code = p.parseStringLiteral().(*ast.StringLiteral)
	return stmt
}
