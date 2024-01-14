package compiler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// TokenType represents a token type.
type TokenType string

// LangToken represents a token in the source code.
type LangToken struct {
	Type    TokenType // LangToken type
	Literal string    // LangToken literal
}

// Lexer represents the lexer for MyLang.
type Lexer struct {
	reader   *bufio.Reader // Buffered reader for source code
	position int           // Current position in the input
	ch       rune          // Current character under examination
}

const (
	TokenTypeUndefined        TokenType = "Undefined"
	TokenTypeEOF              TokenType = "EOF"
	TokenTypeIdentifier       TokenType = "Identifier"
	TokenTypeNumber           TokenType = "Number"
	TokenTypeString           TokenType = "String"
	TokenTypeAssignment       TokenType = "Assignment"
	TokenTypePlus             TokenType = "Plus"
	TokenTypeMinus            TokenType = "Minus"
	TokenTypeMultiply         TokenType = "Multiply"
	TokenTypeDivide           TokenType = "Divide"
	TokenTypeLeftParenthesis  TokenType = "LeftParenthesis"
	TokenTypeRightParenthesis TokenType = "RightParenthesis"
	TokenTypeLeftBrace        TokenType = "LeftBrace"
	TokenTypeRightBrace       TokenType = "RightBrace"
	TokenTypeComma            TokenType = "Comma"
	TokenTypeSemicolon        TokenType = "Semicolon"
	TokenTypeColon            TokenType = "Colon"
	TokenTypeQuestionMark     TokenType = "QuestionMark"
	TokenTypeArrow            TokenType = "Arrow"
	TokenTypeIf               TokenType = "If"
	TokenTypeElse             TokenType = "Else"
	TokenTypeFor              TokenType = "For"
	TokenTypeWhile            TokenType = "While"
	TokenTypeDo               TokenType = "Do"
	TokenTypeSwitch           TokenType = "Switch"
	TokenTypeCase             TokenType = "Case"
	TokenTypeDefault          TokenType = "Default"
	TokenTypeData             TokenType = "Data"
	TokenTypeType             TokenType = "Type"
	TokenTypeAssembly         TokenType = "Assembly"
	TokenTypeMain             TokenType = "Main"
	TokenTypeComment          TokenType = "Comment"
)

// TokenTypeFunction is the token type for the "function" keyword.
const TokenTypeFunction TokenType = "Function"

// TokenTypeLet is the token type for the "let" keyword.
const TokenTypeLet TokenType = "Let"

// Keywords is a map of reserved keywords to their corresponding token types.
var Keywords = map[string]TokenType{
	"function": TokenTypeFunction,
	"let":      TokenTypeLet,
	"if":       TokenTypeIf,
	"else":     TokenTypeElse,
	"for":      TokenTypeFor,
	"while":    TokenTypeWhile,
	"do":       TokenTypeDo,
	"switch":   TokenTypeSwitch,
	"case":     TokenTypeCase,
	"default":  TokenTypeDefault,
	"data":     TokenTypeData,
	"type":     TokenTypeType,
	"asm":      TokenTypeAssembly, // New assembly keyword
	// Add more keywords here
}

// NewLexer creates a new lexer instance with the provided source code file.
func NewLexer(inputFile string) (*Lexer, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lexer := &Lexer{
		reader:   reader,
		position: 0,
		ch:       0,
	}

	lexer.readChar()

	return lexer, nil
}

// NewLexerFromString creates a new lexer instance with the provided source code string.
func NewLexerFromString(input string) (*Lexer, error) {
	reader := bufio.NewReader(strings.NewReader(input))

	lexer := &Lexer{
		reader:   reader,
		position: 0,
		ch:       0,
	}

	lexer.readChar()

	return lexer, nil
}

// readChar reads the next character in the input and advances the position.
func (l *Lexer) readChar() {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		l.ch = 0 // End of file or error
	} else {
		l.ch = ch
	}

	l.position++
}

// NextToken scans and returns the next token in the source code.
func (l *Lexer) NextToken() (LangToken, error) {
	var tok LangToken

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '>' {
			tok = newToken(TokenTypeArrow, l.ch)
			l.readChar() // Consume '>'
		} else {
			tok = newToken(TokenTypeAssignment, l.ch)
		}
	case '?':
		if l.peekChar() == ':' {
			tok = newToken(TokenTypeColon, l.ch)
			l.readChar() // Consume ':'
		} else {
			tok = newToken(TokenTypeQuestionMark, l.ch)
		}
	case '+':
		tok = newToken(TokenTypePlus, l.ch)
	case '-':
		tok = newToken(TokenTypeMinus, l.ch)
	case '*':
		tok = newToken(TokenTypeMultiply, l.ch)
	case '/':
		tok = newToken(TokenTypeDivide, l.ch)
	case '(':
		tok = newToken(TokenTypeLeftParenthesis, l.ch)
	case ')':
		tok = newToken(TokenTypeRightParenthesis, l.ch)
	case '{':
		tok = newToken(TokenTypeLeftBrace, l.ch)
	case '}':
		tok = newToken(TokenTypeRightBrace, l.ch)
	case ',':
		tok = newToken(TokenTypeComma, l.ch)
	case ';':
		tok = newToken(TokenTypeSemicolon, l.ch)
	case '"':
		tok.Type = TokenTypeString
		tok.Literal = l.readString()
	case 'a':
		// Check for the 'asm' keyword
		if l.peekChar() == 's' && l.peekCharAtIndex(2) == 'm' && !isLetter(l.peekCharAtIndex(3)) {
			tok.Literal = l.readAssembly()
			tok.Type = TokenTypeAssembly
			return tok, nil
		}
		// If not 'asm', treat it as an identifier
		tok.Literal = string(l.readIdentifier())
		tok.Type = LookupIdent(tok.Literal)
		return tok, nil
	default:
		if isLetter(l.ch) {
			tok.Literal = string(l.readIdentifier())
			tok.Type = LookupIdent(tok.Literal)
			return tok, nil
		} else if isDigit(l.ch) {
			tok.Type = TokenTypeNumber
			tok.Literal = l.readNumber()
			return tok, nil
		} else {
			return LangToken{}, fmt.Errorf("unexpected character: %s", string(l.ch))
		}
	}

	l.readChar()
	return tok, nil
}

func (l *Lexer) peekChar() rune {
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		return 0
	}
	l.reader.UnreadRune() // Revert the read operation
	return ch
}

// Helper function to skip whitespace characters.
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func LookupIdent(ident string) TokenType {
	if tokType, ok := Keywords[ident]; ok {
		return tokType
	}
	return TokenTypeIdentifier
}

// Helper function to read an identifier or reserved keyword.
func (l *Lexer) readIdentifier() TokenType {
	var identBuilder strings.Builder
	for isLetter(l.ch) || isDigit(l.ch) {
		identBuilder.WriteRune(l.ch)
		l.readChar()
	}
	ident := identBuilder.String()
	if tokType, isKeyword := Keywords[ident]; isKeyword {
		return tokType
	}
	return TokenTypeIdentifier
}

// Helper function to read a number.
func (l *Lexer) readNumber() string {
	var numBuilder strings.Builder
	for isDigit(l.ch) {
		numBuilder.WriteRune(l.ch)
		l.readChar()
	}
	return numBuilder.String()
}

// Helper function to read a string.
func (l *Lexer) readString() string {
	var strBuilder strings.Builder
	l.readChar() // Consume the opening double quote
	for l.ch != '"' && l.ch != 0 {
		strBuilder.WriteRune(l.ch)
		l.readChar()
	}
	return strBuilder.String()
}

// Helper function to read an assembly block.
func (l *Lexer) readAssembly() string {
	var asmBuilder strings.Builder
	for !unicode.IsSpace(l.ch) && l.ch != 0 {
		asmBuilder.WriteRune(l.ch)
		l.readChar()
	}
	return asmBuilder.String()
}

// Helper function to peek the character at a specific index without consuming it.
func (l *Lexer) peekCharAtIndex(index int) rune {
	currentPos := l.position
	var peekRune rune

	for i := 0; i <= index; i++ {
		peekRune, _, _ = l.reader.ReadRune()
	}

	// Reset reader to original position
	l.reader.Reset(l.reader)
	for i := 0; i < currentPos; i++ {
		l.reader.ReadRune()
	}

	return peekRune
}

// Helper function to check if a character is a letter.
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

// Helper function to check if a character is a digit.
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// newToken creates a new LangToken with the given type and literal value.
func newToken(tokenType TokenType, ch rune) LangToken {
	return LangToken{Type: tokenType, Literal: string(ch)}
}
