package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) isFunctionDefinition() bool {
	if p.currentTokenIs(TokenTypeFunction) && p.peekTokenIs(TokenTypeIdentifier) && p.peekTokenAtIndex(1).Type == TokenTypeLeftParenthesis {
		return true
	}
	if p.currentTokenIs(TokenTypeIdentifier) && p.peekTokenIs(TokenTypeLeftParenthesis) {
		return true
	}
	if p.currentTokenIs(TokenTypeLeftParenthesis) {
		return false
	}
	return false
}

func (p *Parser) parseFunctionDefinition() *ast.FunctionDefinition {
	startTokenLine := p.currentToken.Line
	startTokenType := p.currentToken.Type

	fn := &ast.FunctionDefinition{Token: p.currentToken}

	if p.currentTokenIs(TokenTypeFunction) {
		fn.Token = p.currentToken
		p.nextToken()
	}

	if p.currentTokenIs(TokenTypeIdentifier) {
		if p.peekTokenIs(TokenTypeLeftParenthesis) {
			fn.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
			p.nextToken()
		} else {
			if startTokenType == TokenTypeFunction {
				p.errors = append(p.errors, fmt.Sprintf("Expected identifier followed by '(' after 'function' keyword, got '%s' at line %d", p.currentToken.Literal, p.currentToken.Line))
				p.advanceToRecoveryPoint()
				return nil
			}

			p.errors = append(p.errors, fmt.Sprintf("Internal parser error: parseFunctionDefinition called incorrectly for identifier '%s' at line %d", p.currentToken.Literal, p.currentToken.Line))
			return nil
		}
	} else if p.currentTokenIs(TokenTypeLeftParenthesis) {
		if startTokenType == TokenTypeFunction {
			p.errors = append(p.errors, fmt.Sprintf("Cannot use 'function' keyword with anonymous function definition starting at line %d", startTokenLine))
			p.advanceToRecoveryPoint()
			return nil
		}
		fn.Name = nil
	} else {
		p.errors = append(p.errors, fmt.Sprintf("Expected function name identifier or '(' after '%s' keyword, got %s at line %d", startTokenType, p.currentToken.Type, p.currentToken.Line))
		p.advanceToRecoveryPoint()
		return nil
	}

	if !p.currentTokenIs(TokenTypeLeftParenthesis) {
		p.errors = append(p.errors, fmt.Sprintf("Expected '(' for function parameters, got %s at line %d", p.currentToken.Type, p.currentToken.Line))
		p.advanceToRecoveryPoint()
		return nil
	}

	fn.Parameters = p.parseFunctionParameters()
	if fn.Parameters == nil {
		p.advanceToRecoveryPoint()
		return nil
	}

	if p.peekTokenIs(TokenTypeColon) {
		p.nextToken() // Consume ')'
		p.nextToken() // Consume ':'

		if !p.currentTokenIs(TokenTypeIdentifier) {
			p.errors = append(p.errors, fmt.Sprintf("Expected return type identifier after ':', got %s at line %d", p.currentToken.Type, p.currentToken.Line))
			p.advanceToRecoveryPoint()
			return nil
		}
		fn.ReturnType = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		p.nextToken()
	} else {
		p.nextToken() // Consume ')'
	}

	if !p.currentTokenIs(TokenTypeLambdaArrow) {
		p.errors = append(p.errors, fmt.Sprintf("Expected '->' after function signature, got '%s' instead at line %d", p.currentToken.Literal, p.currentToken.Line))
		p.advanceToRecoveryPoint()
		return nil
	}
	fn.Token = p.currentToken
	p.nextToken()

	if p.currentTokenIs(TokenTypeLeftBrace) {
		bodyNode := p.parseBlockStatement()
		if bodyNode == nil {
			p.advanceToRecoveryPoint()
			return nil
		}
		fn.Body = bodyNode
	} else {
		fn.Body = p.parseExpression(LOWEST)
		if fn.Body == nil {
			if !p.errorsEncounteredSince(len(p.errors)) {
				p.errors = append(p.errors, fmt.Sprintf("Failed to parse function body expression at line %d", p.currentToken.Line))
			}
			p.advanceToRecoveryPoint()
			return nil
		}
	}

	return fn
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if !p.currentTokenIs(TokenTypeLeftParenthesis) {
		p.errors = append(p.errors, fmt.Sprintf("Internal Error: parseFunctionParameters called without '(' token at line %d", p.currentToken.Line))
		return nil
	}

	if p.peekTokenIs(TokenTypeRightParenthesis) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	if p.currentTokenIs(TokenTypeRightParenthesis) {
		return identifiers
	}

	if !p.currentTokenIs(TokenTypeIdentifier) {
		p.errors = append(p.errors, fmt.Sprintf("Expected identifier or ')' in parameter list, got %s at line %d", p.currentToken.Type, p.currentToken.Line))
		p.advanceToRecoveryPoint()
		return nil
	}
	ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, ident)
	p.nextToken()

	for p.currentTokenIs(TokenTypeComma) {
		p.nextToken()
		if !p.currentTokenIs(TokenTypeIdentifier) {
			p.errors = append(p.errors, fmt.Sprintf("Expected identifier after comma in parameter list, got %s at line %d", p.currentToken.Type, p.currentToken.Line))
			p.advanceToRecoveryPoint()
			return nil
		}
		ident := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, ident)
		p.nextToken()
	}

	if !p.currentTokenIs(TokenTypeRightParenthesis) {
		p.errors = append(p.errors, fmt.Sprintf("Expected ')' to end parameter list, got %s '%s' instead at line %d", p.currentToken.Type, p.currentToken.Literal, p.currentToken.Line))
		p.advanceToRecoveryPoint()
		return nil
	}

	return identifiers
}
