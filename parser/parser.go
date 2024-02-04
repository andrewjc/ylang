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
	eof          bool

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
	if err != nil {
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

func (p *Parser) parseArrayLiteral() ast.ExpressionNode {
	arrayLit := &ast.ArrayLiteral{Token: p.currentToken}
	arrayLit.Elements = []ast.ExpressionNode{}

	// Skip the opening bracket
	p.nextToken()

	// Parse elements until we reach a closing bracket
	for !p.currentTokenIs(TokenTypeRightBracket) {
		elem := p.parseExpression(LOWEST)
		if elem != nil {
			arrayLit.Elements = append(arrayLit.Elements, elem)
		}

		// Move to the next token, which should be a comma or a closing bracket
		p.nextToken()
		if p.currentTokenIs(TokenTypeComma) {
			p.nextToken() // Skip comma
		}
	}

	// Ensure we have a closing bracket
	if err := p.expectPeek(TokenTypeRightBracket); err != nil {
		fmt.Println(err)
		return nil
	}
	return arrayLit
}

func (p *Parser) parseNumberLiteral() ast.ExpressionNode {
	lit := &ast.NumberLiteral{Token: p.currentToken}

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		// Handle error; could log or set an error on the parser
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.ExpressionNode {
	lit := &ast.StringLiteral{Token: p.currentToken}

	lit.Value = p.currentToken.Literal
	return lit
}

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseParenthesisExpression() ast.ExpressionNode {
	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if err := p.expectPeek(TokenTypeRightParenthesis); err != nil {
		fmt.Println(err)
		return nil
	}

	return exp
}

func (p *Parser) expectPeek(t TokenType) *ParserError {
	if p.peekTokenIs(t) {
		p.nextToken()
		return nil
	} else {
		return &ParserError{
			Message: fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type),
			Line:    p.peekToken.Line,
			Pos:     p.peekToken.Pos,
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

func (p *Parser) isFunctionDefinition() bool {
	return p.peekTokenIs(TokenTypeLeftParenthesis)
}

func (p *Parser) isTernary() bool {
	return p.peekTokenIs(TokenTypeQuestionMark)
}

func (p *Parser) parseTernaryExpression(exp ast.ExpressionNode) ast.ExpressionNode {
	ternary := &ast.TraditionalTernaryExpression{
		Token:     LangToken{},
		Condition: nil,
		TrueExpr:  nil,
		FalseExpr: nil,
	}
	return ternary
}
