package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseVariableDeclaration() *ast.VariableDeclaration {
	// Expect the current token to be 'let'
	if !p.currentTokenIs(TokenTypeLet) {
		return nil
	}

	// Create a new variable declaration node
	varDecl := &ast.VariableDeclaration{
		Token: p.currentToken,
	}

	// Expect the next token to be an identifier (the variable name)
	if err := p.expectPeek(TokenTypeIdentifier); err != nil {
		fmt.Println(err)
		return nil // expectedidentifier
	}

	varDecl.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	// Optional type annotation
	if p.peekTokenIs(TokenTypeLeftParenthesis) {
		p.nextToken() // Consume '('

		if err := p.expectPeek(TokenTypeIdentifier); err != nil {
			fmt.Println(err)
			return nil // expectedidentifier
		}

		varDecl.Type = &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}

		if err := p.expectPeek(TokenTypeRightParenthesis); err != nil {
			fmt.Println(err)
			return nil // Expected ')' after type annotation
		}
	}

	// Expect the next token to be '='
	if err := p.expectPeek(TokenTypeAssignment); err != nil {
		fmt.Println(err)
		return nil // Expected assignment operator
	}

	p.nextToken() // Consume '='

	// Parse the expression that represents the variable's value
	varDecl.Value = p.parseExpression(LOWEST)

	// Handle the end of the statement (e.g., a semicolon in some languages)
	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return varDecl
}

func (p *Parser) parseAssignmentStatement() ast.ExpressionNode {
	stmt := &ast.AssignmentStatement{Token: p.currentToken}

	if err := p.expectPeek(TokenTypeIdentifier); err != nil {
		fmt.Println(err)
		return nil // Expected identifier
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if err := p.expectPeek(TokenTypeAssignment); err != nil {
		fmt.Println(err)
		return nil // Expected assignment operator
	}

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	// Handle semicolon if your language requires it...
	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return stmt
}
