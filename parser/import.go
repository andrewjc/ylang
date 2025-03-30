package parser

import (
	"compiler/ast"
	"fmt"
)
import . "compiler/lexer"

func (p *Parser) parseImportStatement() *ast.ImportStatement {
	stmt := &ast.ImportStatement{Token: p.currentToken}
	if !p.expectPeek(TokenTypeString) {
		p.advanceToRecoveryPoint()
		return nil
	}
	stmt.Path = p.currentToken.Literal

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
		p.nextToken()
	} else {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) advanceToRecoveryPoint() {
	fmt.Printf("[DEBUG] Attempting recovery from error near line %d...\n", p.currentToken.Line)
	recoveryTokens := map[TokenType]bool{
		TokenTypeSemicolon:  true,
		TokenTypeRightBrace: true,
		TokenTypeEOF:        true,
		// Potentially add keywords that start new top-level/block-level items
		TokenTypeLet:      true,
		TokenTypeIf:       true,
		TokenTypeReturn:   true,
		TokenTypeFunction: true,
		TokenTypeType:     true,
		TokenTypeData:     true,
		TokenTypeImport:   true,
	}
	for !recoveryTokens[p.currentToken.Type] {
		p.nextToken()
		if p.currentToken.Type == TokenTypeEOF {
			break
		} // Prevent infinite loop at EOF
	}
	fmt.Printf("[DEBUG] Recovery advanced to token '%s' (%s) near line %d\n", p.currentToken.Literal, p.currentToken.Type, p.currentToken.Line)
}

func (p *Parser) errorsEncounteredSince(countBefore int) bool {
	return len(p.errors) > countBefore
}
