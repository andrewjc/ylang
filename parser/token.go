package parser

import (
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	// Get the next token from the lexer, handling potential lexer errors
	p.peekToken, p.peekTokenErr = p.lexer.NextToken()
	// If lexer returned an error, store it
	if p.peekTokenErr != nil && p.peekToken.Type != TokenTypeEOF { // Don't treat EOF signal as error here
		// Avoid adding duplicate lexer errors if already reported
		isDuplicate := false
		errMsg := fmt.Sprintf("Lexer error: %v at line %d, pos %d", p.peekTokenErr, p.peekToken.Line, p.peekToken.Pos)
		for _, err := range p.errors {
			if err == errMsg {
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
