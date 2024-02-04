package parser

import (
	"compiler/ast"
	. "compiler/lexer"
	"fmt"
)

// Parser represents a parser with an associated lexer, current token and a peek token.
type Parser struct {
	lexer        *Lexer    // The lexer from which tokens are read
	currentToken LangToken // The current token being processed
	peekToken    LangToken // The next token to be processed
	eof          bool

	// Map of infix parse functions, indexed by token type.
	infixParseFns map[TokenType]infixParseFn
}

var DEFAULT_LOGGING_LEAD_LINES = 10
var DEFAULT_LOGGING_FOLLOW_LINES = 10

// infixParseFn defines the type for functions used for infix parsing.
type infixParseFn func(ast.ExpressionNode) ast.ExpressionNode

const (
	LOWEST  int = 1
	SUM     int = 2 // +
	PRODUCT int = 3 // *
	TERNARY int = 4 // Ternary expressions
	ARRAY   int = 5 // Array literals
	DOT     int = 6 // Dot operator
	// Define other precedence levels
)

// precedences defines the order of operation for infix operators
var precedences = map[TokenType]int{
	TokenTypePlus:         SUM,
	TokenTypeMinus:        SUM,
	TokenTypeMultiply:     PRODUCT,
	TokenTypeDivide:       PRODUCT,
	TokenTypeQuestionMark: TERNARY,
	TokenTypeLeftBracket:  ARRAY,
	TokenTypeDot:          DOT,
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
	parser.registerInfix(TokenTypeDot, parser.parseDotOperator)

	// Add other infix operators here

	return parser
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	p.eof = false

	for !p.eof {
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
	if err != nil && p.peekToken.Type != TokenTypeEOF {
		p.eof = true
		fmt.Println("Error reading next token:", err)
	} else if p.peekToken.Type == TokenTypeEOF {
		p.eof = true
	} else {
		p.eof = false
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

func (p *Parser) parseDotOperator(left ast.ExpressionNode) ast.ExpressionNode {
	exp := &ast.MemberAccessExpression{
		Token: p.currentToken,
		Left:  left,
	}

	p.nextToken()
	if !p.currentTokenIs(TokenTypeIdentifier) {
		// Error handling: Expecting an identifier after '.'
		return nil
	}

	exp.Member = p.currentToken.Literal
	p.nextToken()
	return exp
}

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) expectPeek(t TokenType) *ParserError {
	if p.peekTokenIs(t) {
		p.nextToken()
		return nil
	} else {
		line, pos := p.peekToken.Line, p.peekToken.Pos
		snippet := p.lexer.GetCodeFragment(line, pos, DEFAULT_LOGGING_LEAD_LINES, DEFAULT_LOGGING_FOLLOW_LINES) // Get 10 characters around the error location
		return &ParserError{
			Line:         line,
			Pos:          pos,
			Message:      fmt.Sprintf("Expected token %s, got %s", t, p.peekToken.Type),
			CodeFragment: snippet,
		}
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
