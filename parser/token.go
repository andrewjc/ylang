package parser

import (
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.peekToken2
	p.peekToken2 = p.peekToken3
	nextTokenFromLexer, lexErr := p.lexer.NextToken()

	p.peekToken3 = nextTokenFromLexer

	if lexErr != nil && p.peekToken3.Type != TokenTypeEOF { // Report errors unless it's just EOF
		errMsg := fmt.Sprintf("Lexer error: %v at line %d, pos %d", lexErr, nextTokenFromLexer.Line+1, nextTokenFromLexer.Pos)
		isDuplicate := false
		for _, existingErr := range p.errors {
			if existingErr == errMsg {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			p.errors = append(p.errors, errMsg)
		}
	}
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekToken2Is(t TokenType) bool {
	return p.peekToken2.Type == t
}

func (p *Parser) peekToken3Is(t TokenType) bool {
	return p.peekToken3.Type == t
}

func (p *Parser) peekTokenAtIndex(index int) LangToken {
	switch index {
	case 0:
		return p.currentToken
	case 1:
		return p.peekToken
	case 2:
		return p.peekToken2
	case 3:
		return p.peekToken3
	default:
		return LangToken{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 0, Length: 0}
	}
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
