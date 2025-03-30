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
	//p.registerPrefix(TokenTypeComment, p.parseCommentExpression)
	//p.registerPrefix(TokenTypeImport, p.parseImportStatement)

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
	program.Functions = []*ast.FunctionDefinition{}
	program.ClassDeclarations = []*ast.ClassDeclaration{}
	program.DataStructures = []*ast.DataStructure{}
	program.ImportStatements = []*ast.ImportStatement{}

	for !p.currentTokenIs(TokenTypeEOF) {
		parseStartPos := p.lexer.Position
		parseStartToken := p.currentToken

		var parsedItem bool = false // if current token was consumed by a parser

		switch p.currentToken.Type {
		case TokenTypeImport:
			stmtNode := p.parseImportStatement()
			if stmtNode != nil {
				program.ImportStatements = append(program.ImportStatements, stmtNode)
				parsedItem = true
			}

		case TokenTypeFunction, TokenTypeIdentifier:
			looksLikeFunc := (p.currentTokenIs(TokenTypeFunction) && p.peekTokenIs(TokenTypeIdentifier)) ||
				(p.currentTokenIs(TokenTypeIdentifier) && p.peekTokenIs(TokenTypeLeftParenthesis))

			if looksLikeFunc {
				funcNode := p.parseFunctionDefinition()
				if funcNode != nil {
					if funcNode.Name != nil && funcNode.Name.Value == "main" {
						if program.MainFunction != nil {
							p.errors = append(p.errors, fmt.Sprintf("Redefinition of main function at line %d", funcNode.Token.Line))
						}
						program.MainFunction = funcNode
					} else {
						program.Functions = append(program.Functions, funcNode)
					}
					parsedItem = true
				}
			} else if p.currentTokenIs(TokenTypeIdentifier) {
				isDecl := false
				if p.peekTokenIs(TokenTypeArrow) || p.peekTokenIs(TokenTypeAssignment) || p.peekTokenIs(TokenTypeColon) {
					if p.peekTokenIs(TokenTypeArrow) {
						classNode := p.parseClassDeclaration()
						if classNode != nil {
							program.ClassDeclarations = append(program.ClassDeclarations, classNode)
							isDecl = true
						}
					} else {
						dataNode := p.parseDataStructure()
						if dataNode != nil {
							program.DataStructures = append(program.DataStructures, dataNode)
							isDecl = true
						}
					}
				}
				if isDecl {
					parsedItem = true
				} else {
					p.errors = append(p.errors, fmt.Sprintf("Syntax error: Unexpected identifier '%s' at top level near line %d.", p.currentToken.Literal, p.currentToken.Line))
				}
			} else {
				p.errors = append(p.errors, fmt.Sprintf("Syntax error: Expected identifier after 'function' near line %d.", p.currentToken.Line))
			}

		case TokenTypeType, TokenTypeData:
			if p.currentToken.Type == TokenTypeType {
				classNode := p.parseClassDeclaration()
				if classNode != nil {
					program.ClassDeclarations = append(program.ClassDeclarations, classNode)
					parsedItem = true
				}
			} else {
				dataNode := p.parseDataStructure()
				if dataNode != nil {
					program.DataStructures = append(program.DataStructures, dataNode)
					parsedItem = true
				}
			}
		case TokenTypeSemicolon:
			p.nextToken()
			parsedItem = true

		default:
			p.errors = append(p.errors, fmt.Sprintf("Syntax error: Unexpected token '%s' (%s) at top level near line %d.", p.currentToken.Literal, p.currentToken.Type, p.currentToken.Line))
		}

		if !p.currentTokenIs(TokenTypeEOF) && !parsedItem && p.lexer.Position == parseStartPos && p.currentToken.Type == parseStartToken.Type {
			p.errors = append(p.errors, fmt.Sprintf("[CRITICAL] Parser loop stuck on token '%s' (%s) near line %d. Forcing advance.", p.currentToken.Literal, p.currentToken.Type, p.currentToken.Line))
			p.nextToken()
		}
	}

	if program.MainFunction == nil {
		fmt.Println("[INFO] No main function defined in the program.")
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
	if _, isKeyword := Keywords[p.currentToken.Literal]; isKeyword {
		if p.currentToken.Type != TokenTypeIf {
			p.errors = append(p.errors, fmt.Sprintf("unexpected keyword '%s' used as expression at line %d", p.currentToken.Literal, p.currentToken.Line))
			return nil
		}
	}
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseCallExpression(function ast.ExpressionNode) ast.ExpressionNode {
	exp := &ast.CallExpression{Token: p.currentToken, Function: function}
	exp.Arguments = p.parseExpressionList(TokenTypeRightParenthesis)
	return exp
}

func (p *Parser) parseExpressionList(end TokenType) []ast.ExpressionNode {
	var list []ast.ExpressionNode

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()

	firstExpr := p.parseExpression(LOWEST)
	if firstExpr == nil {
		for !p.currentTokenIs(end) && !p.currentTokenIs(TokenTypeEOF) {
			p.nextToken()
		}
		if p.currentTokenIs(end) {
			// p.nextToken() // Consume end token?
		}
		return nil
	}
	list = append(list, firstExpr)

	for p.peekTokenIs(TokenTypeComma) {
		p.nextToken()
		p.nextToken()

		if p.currentTokenIs(end) {
			p.errors = append(p.errors, fmt.Sprintf("Unexpected '%s' after comma in list at line %d", p.currentToken.Literal, p.currentToken.Line))
			break
		}

		nextExpr := p.parseExpression(LOWEST)
		if nextExpr == nil {
			for !p.currentTokenIs(end) && !p.currentTokenIs(TokenTypeEOF) {
				p.nextToken()
			}
			if p.currentTokenIs(end) {
				// p.nextToken() // Consume end token?
			}
			return nil
		}
		list = append(list, nextExpr)
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}
