package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
	"strconv"
)

// Parser represents a parser with an associated lexer, current token and a peek token.
type Parser struct {
	lexer        *Lexer    // The lexer from which tokens are read
	currentToken LangToken // The current token being processed
	peekToken    LangToken // The next token to be processed

	// Map of infix parse functions, indexed by token type.
	infixParseFns map[TokenType]infixParseFn
}

// infixParseFn defines the type for functions used for infix parsing.
type infixParseFn func(ast.ExpressionNode) ast.ExpressionNode

const (
	LOWEST  int = 1
	SUM     int = 2 // +
	PRODUCT int = 3 // *
	TERNARY int = 4 // Ternary expressions
	// Define other precedence levels
)

// precedences defines the order of operation for infix operators
var precedences = map[TokenType]int{
	TokenTypePlus:         SUM,
	TokenTypeMinus:        SUM,
	TokenTypeMultiply:     PRODUCT,
	TokenTypeDivide:       PRODUCT,
	TokenTypeQuestionMark: TERNARY,
	// Add other operators and their precedences
}

// NewParser creates and initializes a new Parser.
func NewParser(lexer *Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	// Initialize both currentToken and peekToken
	parser.nextToken()
	parser.nextToken()

	// Register infix parse functions
	parser.infixParseFns = make(map[TokenType]infixParseFn)
	parser.registerInfix(TokenTypePlus, parser.parseInfixExpression)
	parser.registerInfix(TokenTypeMinus, parser.parseInfixExpression)
	parser.registerInfix(TokenTypeMultiply, parser.parseInfixExpression)
	parser.registerInfix(TokenTypeDivide, parser.parseInfixExpression)
	parser.registerInfix(TokenTypeQuestionMark, parser.parseTraditionalTernaryExpression)
	parser.registerInfix(TokenTypeLambdaArrow, parser.parseLambdaStyleTernaryExpression)
	parser.registerInfix(TokenTypeIf, parser.parseInlineIfElseTernaryExpression)

	// Add other infix operators here

	return parser
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(TokenTypeEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// nextToken advances both the currentToken and the peekToken.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	var err error
	p.peekToken, err = p.lexer.NextToken()
	if p.peekToken.Type != TokenTypeEOF && err != nil {
		// Error handling (could be logging or panic depending on your design)
		fmt.Println("Error reading next token:", err)
	}
}

// peekTokenType returns the type of the peek token.
func (p *Parser) peekTokenType() TokenType {
	return p.peekToken.Type
}

// peekTokenIs checks if the peek token is of a certain type.
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekTokenType() == t
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	// Initial expression parsing based on the current token
	var leftExp ast.ExpressionNode

	switch p.currentToken.Type {
	case TokenTypeNumber:
		return p.parseNumberLiteral()
	case TokenTypeString:
		return p.parseStringLiteral()
	case TokenTypeIdentifier:
		return p.parseIdentifier()
	case TokenTypeLeftParenthesis:
		return p.parseParenExpression()
	case TokenTypeQuestionMark:
		return p.parseTraditionalTernaryExpression(leftExp)
	case TokenTypeLambdaArrow:
		return p.parseLambdaStyleTernaryExpression(leftExp)
	case TokenTypeIf:
		return p.parseInlineIfElseTernaryExpression(leftExp)
	case TokenTypeLet:
		return p.parseVariableDeclaration()
	case TokenTypeAssignment:
		return p.parseAssignmentStatement()
	case TokenTypeLeftBrace:
		return p.parseBlockStatement()
	}

	for !p.peekTokenIs(TokenTypeSemicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)

	}
	return leftExp
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case TokenTypeLet:
		return p.parseVariableDeclaration()
	case TokenTypeIf:
		if p.peekTokenIs(TokenTypeLeftParenthesis) {
			return p.parseIfStatement()
		} else {
			return p.parseLambdaIfStatement()
		}
	// Include other cases for different statement types
	default:
		return p.parseExpressionStatement() // Default to expression statement
	}
}

func (p *Parser) parseNumberLiteral() ast.ExpressionNode {
	lit := &ast.NumberLiteral{Token: p.currentToken}

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		// Handle error; could log or set an error on the parser
		return nil
	}

	lit.Value = value
	p.nextToken()
	return lit
}

func (p *Parser) parseStringLiteral() ast.ExpressionNode {
	lit := &ast.StringLiteral{Token: p.currentToken}

	lit.Value = p.currentToken.Literal
	p.nextToken()
	return lit
}

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseParenExpression() ast.ExpressionNode {
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(TokenTypeRightParenthesis) {
		return nil // Error handling; expected a closing parenthesis
	}

	return exp
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		// Add error handling here
		return false
	}
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {

	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseInfixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// Utility functions to determine the precedence of tokens
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
