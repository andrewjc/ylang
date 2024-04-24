package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseClassDeclaration() *ast.ClassDeclaration {
	classDecl := &ast.ClassDeclaration{Token: p.currentToken}

	if p.currentTokenIs(TokenTypeIdentifier) {
		classDecl.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		classDecl.LambdaStyle = true
		p.nextToken()

		if !p.expectPeek(TokenTypeArrow) {
			fmt.Println("Expected '->' after class name")
			return nil
		}
	} else if p.currentTokenIs(TokenTypeType) {
		p.nextToken()

		if !p.expectPeek(TokenTypeIdentifier) {
			fmt.Println("Expected identifier after 'type'")
			return nil
		}

		classDecl.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		classDecl.LambdaStyle = false
	} else {
		return nil
	}

	if !p.expectPeek(TokenTypeLeftBrace) {
		fmt.Println("Expected '{' after class name")
		return nil
	}

	classDecl.Members = p.parseClassMembers()

	if !p.expectPeek(TokenTypeRightBrace) {
		fmt.Println("Expected '}' after class members")
		return nil
	}

	return classDecl
}

func (p *Parser) parseClassMembers() []*ast.ClassMember {
	var members []*ast.ClassMember

	for !p.currentTokenIs(TokenTypeRightBrace) {
		var member *ast.ClassMember

		if p.currentTokenIs(TokenTypeLet) {
			variableDecl := p.parseVariableDeclaration()
			if variableDecl != nil {
				member = &ast.ClassMember{VariableDeclaration: variableDecl.(*ast.VariableDeclaration)}
			}
		} else {
			methodDecl := p.parseMethodDeclaration()
			if methodDecl != nil {
				member = &ast.ClassMember{MethodDeclaration: methodDecl}
			}
		}

		if member != nil {
			members = append(members, member)
		}

		p.nextToken()
	}

	return members
}

func (p *Parser) parseMethodDeclaration() *ast.MethodDeclaration {
	method := &ast.MethodDeclaration{Token: p.currentToken}

	// Parse return type (optional)
	if p.peekTokenIs(TokenTypeIdentifier) && p.peekTokenAtIndex(2).Type == TokenTypeIdentifier {
		p.nextToken()
		method.ReturnType = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	}

	// Parse method name
	if !p.expectPeek(TokenTypeIdentifier) {
		fmt.Println("Expected method name")
		return nil
	}
	method.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	// Parse parameters
	if !p.expectPeek(TokenTypeLeftParenthesis) {
		fmt.Println("Expected '(' after method name")
		return nil
	}
	method.Parameters = p.parseParameterList()
	if !p.expectPeek(TokenTypeRightParenthesis) {
		fmt.Println("Expected ')' after method parameters")
		return nil
	}

	// Parse lambda arrow
	if !p.expectPeek(TokenTypeLambdaArrow) {
		fmt.Println("Expected '->' after method parameters")
		return nil
	}

	// Parse method body
	method.Body = p.parseBlockStatement()

	return method
}

func (p *Parser) parseVariableDeclaration() ast.ExpressionNode {
	varDecl := &ast.VariableDeclaration{Token: p.currentToken}

	if !p.expectPeek(TokenTypeIdentifier) {
		return nil
	}

	varDecl.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if p.peekTokenIs(TokenTypeAssignment) {
		p.nextToken()

		p.nextToken()
		varDecl.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(TokenTypeSemicolon) {
		p.nextToken()
	}

	return varDecl
}

func (p *Parser) parseParameterList() []*ast.Parameter {
	var parameters []*ast.Parameter

	if p.peekTokenIs(TokenTypeRightParenthesis) {
		return parameters
	}

	p.nextToken()

	param := &ast.Parameter{Token: p.currentToken}
	param.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if p.peekTokenIs(TokenTypeLeftParenthesis) {
		p.nextToken()
		if !p.expectPeek(TokenTypeIdentifier) {
			fmt.Println("Expected type after '('")
			return nil
		}
		param.Type = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		if !p.expectPeek(TokenTypeRightParenthesis) {
			fmt.Println("Expected ')' after type")
			return nil
		}
	}

	parameters = append(parameters, param)

	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken()
		p.nextToken()

		param := &ast.Parameter{Token: p.currentToken}
		param.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

		if p.peekTokenIs(TokenTypeLeftParenthesis) {
			p.nextToken()
			if !p.expectPeek(TokenTypeIdentifier) {
				fmt.Println("Expected type after '('")
				return nil
			}
			param.Type = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
			if !p.expectPeek(TokenTypeRightParenthesis) {
				fmt.Println("Expected ')' after type")
				return nil
			}
		}

		parameters = append(parameters, param)
	}

	return parameters
}

func (p *Parser) parseOnConstructStatement() *ast.OnConstructStatement {
	stmt := &ast.OnConstructStatement{Token: p.currentToken}
	if !p.expectPeek(TokenTypeIdentifier) {
		fmt.Println("Expected identifier after 'on construct'")
		return nil
	}
	stmt.Lambda = p.parseLambdaExpression().(*ast.LambdaExpression)
	return stmt
}

func (p *Parser) parseOnDestructStatement() *ast.OnDestructStatement {
	stmt := &ast.OnDestructStatement{Token: p.currentToken}
	if !p.expectPeek(TokenTypeIdentifier) {
		fmt.Println("Expected identifier after 'on destruct'")
		return nil
	}
	stmt.Lambda = p.parseLambdaExpression().(*ast.LambdaExpression)
	return stmt
}
