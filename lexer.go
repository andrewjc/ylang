package compiler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Lexer struct {
	reader     *bufio.Reader
	position   int
	ch         rune
	peekBuffer []rune
	eof        bool
}

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

func (l *Lexer) readChar() {
	if len(l.peekBuffer) > 0 {
		l.ch = l.peekBuffer[0]
		l.peekBuffer = l.peekBuffer[1:]
	} else {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			l.ch = 0 // End of file or error
		} else {
			l.ch = ch
		}
	}
	l.position++
}

func (l *Lexer) NextToken() (LangToken, error) {
	var tok LangToken

	if l.eof {
		return LangToken{}, errors.New("no more tokens, reached end of file")
	}

	l.skipWhitespace()

	switch l.ch {
	case 0: // Handling EOF
		tok.Literal = ""
		tok.Type = TokenTypeEOF
		l.eof = true
		return tok, nil
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
		if l.peekChar() == '>' {
			tok = LangToken{Type: TokenTypeLambdaArrow, Literal: "->"}
			l.readChar() // Consume '>'
		} else {
			tok = newToken(TokenTypeMinus, l.ch)
		}
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
	case '<':
		if l.peekChar() == '=' {
			tok = newToken(TokenTypeLessThanEqual, l.ch)
			l.readChar() // Consume '='
		} else {
			tok = newToken(TokenTypeLessThan, l.ch)
		}
	case '>':
		tok = newToken(TokenTypeGreaterThan, l.ch)
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
		tokType, literal := l.readIdentifier()
		return LangToken{Type: tokType, Literal: literal}, nil
	case 'f':
		// Check for the 'function' keyword
		if l.peekChar() == 'u' && l.peekCharAtIndex(2) == 'n' &&
			l.peekCharAtIndex(3) == 'c' && l.peekCharAtIndex(4) == 't' &&
			l.peekCharAtIndex(5) == 'i' && l.peekCharAtIndex(6) == 'o' &&
			l.peekCharAtIndex(7) == 'n' && !isLetter(l.peekCharAtIndex(8)) {
			tok.Literal = l.readFunction()
			tok.Type = TokenTypeFunction
			return tok, nil
		}
		// If not 'function', treat it as an identifier
		tokType, literal := l.readIdentifier()
		return LangToken{Type: tokType, Literal: literal}, nil
	case 'i':
		// Check for the 'if' keyword
		if l.peekChar() == 'f' && !isLetter(l.peekCharAtIndex(2)) {
			tokType, literal := l.readIdentifier()
			tok.Literal = literal
			tok.Type = tokType
			return tok, nil
		}
		if l.peekChar() == 'n' && !isLetter(l.peekCharAtIndex(2)) {
			tokType, literal := l.readIdentifier()
			tok.Literal = literal
			tok.Type = tokType
			return tok, nil
		}
		// If not 'if' or 'in', treat it as an identifier
		tokType, literal := l.readIdentifier()
		return LangToken{Type: tokType, Literal: literal}, nil
	default:
		if isLetter(l.ch) {
			tokType, literal := l.readIdentifier()
			return LangToken{Type: tokType, Literal: literal}, nil
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

func (l *Lexer) readFunction() string {
	var funcBuilder strings.Builder
	for !unicode.IsSpace(l.ch) && l.ch != 0 {
		funcBuilder.WriteRune(l.ch)
		l.readChar()
	}
	return funcBuilder.String()
}

func (l *Lexer) peekChar() rune {
	if len(l.peekBuffer) == 0 {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			return 0
		}
		l.peekBuffer = append(l.peekBuffer, ch)
	}
	return l.peekBuffer[0]
}

func (l *Lexer) peekCharAtIndex(index int) rune {
	for len(l.peekBuffer) <= index {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			return 0
		}
		l.peekBuffer = append(l.peekBuffer, ch)
	}
	return l.peekBuffer[index]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() (TokenType, string) {
	var identBuilder strings.Builder
	for isLetter(l.ch) || isDigit(l.ch) {
		identBuilder.WriteRune(l.ch)
		l.readChar()
	}
	ident := identBuilder.String()
	if tokType, isKeyword := Keywords[ident]; isKeyword {
		return tokType, ident
	}
	return TokenTypeIdentifier, ident
}

func (l *Lexer) readNumber() string {
	var numBuilder strings.Builder
	for isDigit(l.ch) {
		numBuilder.WriteRune(l.ch)
		l.readChar()
	}
	return numBuilder.String()
}

func (l *Lexer) readString() string {
	var strBuilder strings.Builder
	l.readChar() // Consume the opening double quote
	for l.ch != '"' && l.ch != 0 {
		strBuilder.WriteRune(l.ch)
		l.readChar()
	}
	return strBuilder.String()
}

func (l *Lexer) readAssembly() string {
	var asmBuilder strings.Builder
	for !unicode.IsSpace(l.ch) && l.ch != 0 {
		asmBuilder.WriteRune(l.ch)
		l.readChar()
	}
	return asmBuilder.String()
}
