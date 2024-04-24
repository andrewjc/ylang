package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

func (p *Parser) parseDataStructure() *ast.DataStructure {
	dataStruct := &ast.DataStructure{Token: p.currentToken}

	if p.currentTokenIs(TokenTypeData) {
		p.nextToken()

		if p.currentTokenIs(TokenTypeIdentifier) {
			dataStruct.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

			switch p.peekToken.Type {
			case TokenTypeLeftBrace:
				dataStruct.Style = ast.DataStructureStyleBraces
			case TokenTypeColon:
				dataStruct.Style = ast.DataStructureStyleColon
			default:
				return nil
			}

			p.nextToken()
		} else {
			return nil
		}
	} else if p.currentTokenIs(TokenTypeIdentifier) {
		dataStruct.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

		if p.peekTokenIs(TokenTypeAssignment) {
			p.nextToken()
			if p.peekTokenIs(TokenTypeData) {
				dataStruct.Style = ast.DataStructureStyleEquals
				p.nextToken()
			} else {
				dataStruct.Style = ast.DataStructureStyleTupleLike
			}
		} else {
			return nil
		}

		if !p.expectPeek(TokenTypeLeftBrace) {
			fmt.Println("Expected '{' after data structure name")
			return nil
		}
	} else {
		return nil
	}

	dataStruct.Fields = p.parseFieldList()

	if !p.expectPeek(TokenTypeRightBrace) {
		fmt.Println("Expected '}' after data structure fields")
		return nil
	}

	return dataStruct
}

func (p *Parser) parseFieldList() []*ast.Field {
	var fields []*ast.Field

	if !p.expectPeek(TokenTypeLet) {
		fmt.Println("Expected 'let' after '{'")
		return nil
	}

	field := &ast.Field{Token: p.currentToken}
	if !p.expectPeek(TokenTypeIdentifier) {
		fmt.Println("Expected field name after 'let'")
		return nil
	}
	field.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	fields = append(fields, field)

	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken()
		p.nextToken()

		field := &ast.Field{Token: p.currentToken}
		if !p.expectPeek(TokenTypeIdentifier) {
			fmt.Println("Expected field name after ','")
			return nil
		}
		field.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		fields = append(fields, field)
	}

	return fields
}
