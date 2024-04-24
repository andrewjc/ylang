package parser

import (
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken, p.peekTokenErr = p.lexer.NextToken()
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekTokenAtIndex(index int) LangToken {
	// Save the current position
	currentPosition := p.lexer.Position

	// Move to the desired token index
	for i := 0; i < index; i++ {
		_, err := p.lexer.NextToken()
		if err != nil {
			fmt.Println("Error getting token at index", index)
			return LangToken{}
		}
	}

	// Get the token at the desired index
	token, err := p.lexer.NextToken()
	if err != nil {
		return LangToken{}
	}

	// Restore the lexer position
	p.lexer.Position = currentPosition
	p.lexer.ReadChar()

	return token
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
