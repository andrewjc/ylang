package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) isFunctionDefinition() bool {
	return p.currentTokenIs(TokenTypeLeftParenthesis) || p.currentTokenIs(TokenTypeLambdaArrow)
}

func (p *Parser) parseFunctionDefinition() ast.ExpressionNode {
	fn := &ast.FunctionDefinition{Token: p.currentToken}

	if !p.currentTokenIs(TokenTypeRightParenthesis) && !p.currentTokenIs(TokenTypeLambdaArrow) {
		line, pos := p.peekToken.Line, p.peekToken.Pos
		snippet := p.lexer.GetCodeFragment(line, pos, DEFAULT_LOGGING_LEAD_LINES, DEFAULT_LOGGING_FOLLOW_LINES) // Get 10 characters around the error location
		parseError := &ParserError{
			Line:         line,
			Pos:          pos,
			Message:      fmt.Sprintf("Syntax error: Expected ')' or '->' after function parameters, got %s", p.peekToken.Type),
			CodeFragment: snippet,
		}
		fmt.Println(parseError)
		return nil
	}

	if p.currentTokenIs(TokenTypeRightParenthesis) {
		p.nextToken()
		fn.Parameters = p.parseFunctionParameters()

		p.nextToken()

	} else if p.currentTokenIs(TokenTypeLambdaArrow) {
		p.nextToken()
		fn.Parameters = p.parseFunctionParameters()

		p.nextToken()
	}

	// Parse the function body
	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parseFunctionParameters() []string {
	var parameters []string
	for !p.currentTokenIs(TokenTypeRightParenthesis) {
		if p.currentTokenIs(TokenTypeIdentifier) {
			parameters = append(parameters, p.currentToken.Literal)
		}
		p.nextToken()
		if p.currentTokenIs(TokenTypeComma) {
			p.nextToken()
		}
	}

	return parameters
}
