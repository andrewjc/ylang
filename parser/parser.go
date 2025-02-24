package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	ASSIGN      // =
	TERNARY     // ?:
)

var precedences = map[TokenType]int{
	TokenTypeEqual:           EQUALS,
	TokenTypeLessThan:        LESSGREATER,
	TokenTypeGreaterThan:     LESSGREATER,
	TokenTypePlus:            SUM,
	TokenTypeMinus:           SUM,
	TokenTypeMultiply:        PRODUCT,
	TokenTypeDivide:          PRODUCT,
	TokenTypeLeftParenthesis: CALL,
	TokenTypeLeftBracket:     INDEX,
	TokenTypeDot:             INDEX,
	TokenTypeQuestionMark:    TERNARY,
	TokenTypeLambdaArrow:     TERNARY,
	TokenTypeIf:              TERNARY,
	TokenTypeAssignment:      ASSIGN,
}

var DEFAULT_LOGGING_LEAD_LINES = 10
var DEFAULT_LOGGING_FOLLOW_LINES = 10

type (
	prefixParseFn func() ast.ExpressionNode
	infixParseFn  func(ast.ExpressionNode) ast.ExpressionNode
)

type Parser struct {
	lexer  *Lexer
	errors []string

	currentToken LangToken
	peekToken    LangToken
	peekTokenErr error

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.registerPrefix(TokenTypeIdentifier, p.parseIdentifier)
	p.registerPrefix(TokenTypeNumber, p.parseNumberLiteral)
	p.registerPrefix(TokenTypeString, p.parseStringLiteral)
	p.registerPrefix(TokenTypeLeftParenthesis, p.parseParenthesisExpression)
	p.registerPrefix(TokenTypeLeftBracket, p.parseArrayLiteral)
	p.registerPrefix(TokenTypeLambdaArrow, p.parseLambdaExpression)
	p.registerPrefix(TokenTypeIf, p.parseIfStatement)
	p.registerPrefix(TokenTypeLet, p.parseVariableDeclaration)
	p.registerPrefix(TokenTypeLeftBrace, p.parseBlockStatement)
	p.registerPrefix(TokenTypeSyscall, p.parseSysCallExpression)

	p.infixParseFns = make(map[TokenType]infixParseFn)
	p.registerInfix(TokenTypePlus, p.parseInfixExpression)
	p.registerInfix(TokenTypeMinus, p.parseInfixExpression)
	p.registerInfix(TokenTypeMultiply, p.parseInfixExpression)
	p.registerInfix(TokenTypeDivide, p.parseInfixExpression)
	p.registerInfix(TokenTypeEqual, p.parseInfixExpression)
	p.registerInfix(TokenTypeLessThan, p.parseInfixExpression)
	p.registerInfix(TokenTypeGreaterThan, p.parseInfixExpression)
	p.registerInfix(TokenTypeLeftParenthesis, p.parseCallExpression)
	p.registerInfix(TokenTypeLeftBracket, p.parseIndexExpression)
	p.registerInfix(TokenTypeDot, p.parseDotOperator)
	p.registerInfix(TokenTypeQuestionMark, p.parseTraditionalTernaryExpression)
	p.registerInfix(TokenTypeLambdaArrow, p.parseLambdaStyleTernaryExpression)
	p.registerInfix(TokenTypeIf, p.parseInlineIfElseTernaryExpression)
	p.registerInfix(TokenTypeAssignment, p.parseAssignmentExpression)

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	// Parse class declarations, functions, and data structures
	for !p.currentTokenIs(TokenTypeEOF) {
		switch p.currentToken.Type {
		case TokenTypeIdentifier, TokenTypeType:
			// If next token is '(' => parse a function definition
			if p.peekTokenIs(TokenTypeLeftParenthesis) {
				function := p.parseFunctionDefinition()
				if function != nil {
					if function.Name != nil && function.Name.Value == "main" {
						program.MainFunction = function
					} else {
						program.Functions = append(program.Functions, function)
					}
				}
			} else if p.peekTokenIs(TokenTypeArrow) {
				classDecl := p.parseClassDeclaration()
				if classDecl != nil {
					program.ClassDeclarations = append(program.ClassDeclarations, classDecl)
				}
			} else {
				dataStruct := p.parseDataStructure()
				if dataStruct != nil {
					program.DataStructures = append(program.DataStructures, dataStruct)
				}
			}

		default:
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead at line %d, position %d",
		t, p.peekToken.Type, p.currentToken.Line, p.currentToken.Pos)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseCallExpression(function ast.ExpressionNode) ast.ExpressionNode {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.parseExpressionList(TokenTypeRightParenthesis)
	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []ast.ExpressionNode {
	var list []ast.ExpressionNode

	// If the very next token is our 'end', then it's an empty list.
	if p.peekTokenIs(end) {
		p.nextToken() // consume the ']' (or whatever 'end' is)
		return list
	}

	// Otherwise, parse the first expression.
	p.nextToken() // advance so p.currentToken is now the first expression
	first := p.parseExpression(LOWEST)
	if first != nil {
		list = append(list, first)
	}

	// As long as the next token is a comma, consume it and parse another expression.
	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken() // consume the comma
		p.nextToken() // move to the next expression
		expr := p.parseExpression(LOWEST)
		if expr != nil {
			list = append(list, expr)
		}
	}

	// Finally, we expect the 'end' token. E.g. a closing bracket.
	if !p.expectPeek(end) {
		// If it's missing, we can return nil or report an error
		return nil
	}

	return list
}
