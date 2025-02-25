package lexer

import (
	"bufio"
	"compiler/common"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Lexer struct {
	reader      *bufio.Reader
	Position    int
	line        int
	linePos     int
	ch          rune
	peekBuffer  []rune
	eof         bool
	peekToken   LangToken
	lines       []string
	currentLine strings.Builder
	start       int
	end         int
}

const DEFAULT_ROLLING_BUFFER = 10

func NewLexer(inputFile string) (*Lexer, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lexer := &Lexer{
		reader:   reader,
		Position: 0,
		ch:       0,
		lines:    make([]string, DEFAULT_ROLLING_BUFFER),
	}

	lexer.ReadChar()

	return lexer, nil
}

func NewLexerFromString(input string) (*Lexer, error) {
	reader := bufio.NewReader(strings.NewReader(input))

	lexer := &Lexer{
		reader:   reader,
		Position: 0,
		ch:       0,
		lines:    make([]string, DEFAULT_ROLLING_BUFFER),
	}

	lexer.ReadChar()

	return lexer, nil
}

func (l *Lexer) ReadChar() {
	prevCh := l.ch
	if len(l.peekBuffer) > 0 {
		l.ch = l.peekBuffer[0]
		l.peekBuffer = l.peekBuffer[1:]
	} else {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			l.ch = 0 // End of file or error
			if prevCh != '\n' {
				l.lines = append(l.lines, l.currentLine.String())
				l.currentLine.Reset()
			}
		} else {
			l.ch = ch
		}
	}
	l.Position++
	if prevCh == '\n' {
		l.line++
		l.linePos = 0
		l.lines = append(l.lines, l.currentLine.String())
		l.currentLine.Reset()
	} else {
		l.linePos++
		l.currentLine.WriteRune(l.ch)
	}
	l.updatePeekTokenPosition()
}

func (l *Lexer) getCurrentLineNumber() int {
	return l.line
}

func (l *Lexer) getCurrentLinePosition() int {
	return l.linePos
}

func (l *Lexer) updatePeekTokenPosition() {
	line := l.getCurrentLineNumber()
	pos := l.getCurrentLinePosition()

	l.peekToken.Line = line
	l.peekToken.Pos = pos
}

func (l *Lexer) NextToken() (LangToken, error) {
	var tok LangToken

	if l.eof {
		return LangToken{Type: TokenTypeEOF}, errors.New("no more tokens, reached end of file")
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
			l.ReadChar() // Consume '>'
		} else {
			tok = newToken(TokenTypeAssignment, l.ch)
		}
	case '?':
		if l.peekChar() == ':' {
			tok = newToken(TokenTypeColon, l.ch)
			l.ReadChar() // Consume ':'
		} else {
			tok = newToken(TokenTypeQuestionMark, l.ch)
		}
	case '+':
		tok = newToken(TokenTypePlus, l.ch)
	case '-':
		if l.peekChar() == '>' {
			tok = LangToken{Type: TokenTypeLambdaArrow, Literal: "->"}
			l.ReadChar() // Consume '>'
		} else {
			tok = newToken(TokenTypeMinus, l.ch)
		}
	case '*':
		tok = newToken(TokenTypeMultiply, l.ch)
	case '/':
		if l.peekChar() == '/' {
			// Consume the comment
			for l.ch != '\n' && l.ch != 0 {
				l.ReadChar()
			}
			return l.NextToken()
		} else if l.peekChar() == '*' {
			// Consume the comment
			for {
				l.ReadChar()
				if l.ch == '*' && l.peekChar() == '/' {
					l.ReadChar() // Consume '/'
					break
				}
			}
			return l.NextToken()
		} else {
			tok = newToken(TokenTypeDivide, l.ch)
		}
	case '(':
		tok = newToken(TokenTypeLeftParenthesis, l.ch)
	case ')':
		tok = newToken(TokenTypeRightParenthesis, l.ch)
	case '{':
		tok = newToken(TokenTypeLeftBrace, l.ch)
	case '}':
		tok = newToken(TokenTypeRightBrace, l.ch)
	case '[':
		tok = newToken(TokenTypeLeftBracket, l.ch)
	case ']':
		tok = newToken(TokenTypeRightBracket, l.ch)
	case ',':
		tok = newToken(TokenTypeComma, l.ch)
	case ';':
		tok = newToken(TokenTypeSemicolon, l.ch)
	case '<':
		if l.peekChar() == '=' {
			tok = newToken(TokenTypeLessThanEqual, l.ch)
			l.ReadChar() // Consume '='
		} else {
			tok = newToken(TokenTypeLessThan, l.ch)
		}
	case '>':
		tok = newToken(TokenTypeGreaterThan, l.ch)
	case '"', '`', '\'':
		tok.Type = TokenTypeString
		tok.Literal = l.readString()
	case 'a':
		// Check for the 'asm' keyword
		if l.peekChar() == 's' && l.peekCharAtIndex(2) == 'm' && !common.IsLetter(l.peekCharAtIndex(3)) {
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
			l.peekCharAtIndex(7) == 'n' && !common.IsLetter(l.peekCharAtIndex(8)) {
			tok.Literal = l.readFunction()
			tok.Type = TokenTypeFunction
			return tok, nil
		}
		// If not 'function', treat it as an identifier
		tokType, literal := l.readIdentifier()
		return LangToken{Type: tokType, Literal: literal}, nil
	case 'i':
		// Check for the 'if' keyword
		if l.peekChar() == 'f' && !common.IsLetter(l.peekCharAtIndex(2)) {
			tokType, literal := l.readIdentifier()
			tok.Literal = literal
			tok.Type = tokType
			return tok, nil
		}
		if l.peekChar() == 'n' && !common.IsLetter(l.peekCharAtIndex(2)) {
			tokType, literal := l.readIdentifier()
			tok.Literal = literal
			tok.Type = tokType
			return tok, nil
		}
		// If not 'if' or 'in', treat it as an identifier
		tokType, literal := l.readIdentifier()
		return LangToken{Type: tokType, Literal: literal}, nil
	case 'r':
		// Check for the 'return' keyword
		if l.peekChar() == 'e' && l.peekCharAtIndex(2) == 't' &&
			l.peekCharAtIndex(3) == 'u' && l.peekCharAtIndex(4) == 'r' &&
			l.peekCharAtIndex(5) == 'n' && !common.IsLetter(l.peekCharAtIndex(6)) {
			tok.Literal = l.readReturn()
			tok.Type = TokenTypeReturn
			return tok, nil
		}
		// If not 'return', treat it as an identifier
		tokType, literal := l.readIdentifier()
		return LangToken{Type: tokType, Literal: literal}, nil
	case '.':
		tok = newToken(TokenTypeDot, l.ch)
	default:
		if common.IsLetter(l.ch) {
			tokType, literal := l.readIdentifier()
			return LangToken{Type: tokType, Literal: literal}, nil
		} else if common.IsDigit(l.ch) {
			tok.Type = TokenTypeNumber
			tok.Literal = l.readNumber()
			return tok, nil
		} else {
			return LangToken{}, fmt.Errorf("unexpected character: %s", string(l.ch))
		}
	}

	l.ReadChar()

	tok.Line = l.getCurrentLineNumber()
	tok.Pos = l.getCurrentLinePosition()
	return tok, nil
}

func (l *Lexer) readFunction() string {
	var funcBuilder strings.Builder
	for !unicode.IsSpace(l.ch) && l.ch != 0 {
		funcBuilder.WriteRune(l.ch)
		l.ReadChar()
	}
	return funcBuilder.String()
}

func (l *Lexer) readReturn() string {
	var returnBuilder strings.Builder
	for !unicode.IsSpace(l.ch) && l.ch != 0 {
		returnBuilder.WriteRune(l.ch)
		l.ReadChar()
	}
	return returnBuilder.String()
}

func (l *Lexer) peekChar() rune {
	if len(l.peekBuffer) == 0 {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			return 0
		}
		l.peekBuffer = append(l.peekBuffer, ch)
		if ch == '\n' {
			l.line++
			l.linePos = 0
		} else {
			l.linePos++
		}
		l.updatePeekTokenPosition()
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
		if ch == '\n' {
			l.line++
			l.linePos = 0
		} else {
			l.linePos++
		}
		l.updatePeekTokenPosition()
	}
	return l.peekBuffer[index]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.ReadChar()
	}
}

func (l *Lexer) readIdentifier() (TokenType, string) {
	var identBuilder strings.Builder
	for common.IsLetter(l.ch) || common.IsDigit(l.ch) {
		identBuilder.WriteRune(l.ch)
		l.ReadChar()
	}
	ident := identBuilder.String()
	if tokType, isKeyword := Keywords[ident]; isKeyword {
		return tokType, ident
	}
	return TokenTypeIdentifier, ident
}

func (l *Lexer) readAssembly() string {
	l.skipWhitespace() // Skip leading whitespace before reading assembly instruction

	var asmBuilder strings.Builder
	for !unicode.IsSpace(l.ch) && l.ch != 0 {
		asmBuilder.WriteRune(l.ch)
		l.ReadChar()
	}
	return asmBuilder.String()
}

func (l *Lexer) GetCodeFragment(line int, pos int, leadingLines int, followingLines int) string {
	start := (l.end - leadingLines + DEFAULT_ROLLING_BUFFER) % DEFAULT_ROLLING_BUFFER
	end := (l.end + followingLines) % DEFAULT_ROLLING_BUFFER

	var lines []string
	for i := start; i != end; i = (i + 1) % DEFAULT_ROLLING_BUFFER {
		lines = append(lines, l.lines[i])
	}

	// If there's content in the currentLine buffer, append it to the lines array
	if l.currentLine.Len() > 0 {
		lines = append(lines, l.currentLine.String())
	}

	return strings.Join(lines, "\n")
}
