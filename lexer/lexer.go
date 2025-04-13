// lexer/lexer.go
package lexer

import (
	"bufio"
	"compiler/common"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	reader     *bufio.Reader
	Position   int  // Overall position in the input stream (rune count)
	line       int  // Current line number (0-based)
	linePos    int  // Current column position on the line (1-based)
	ch         rune // Current character under examination
	peekBuffer []rune
	eof        bool
}

// Note: DEFAULT_ROLLING_BUFFER and related fields/logic are commented out
// as they complicate position tracking and weren't fully utilized for errors.
// Simpler line/pos tracking is used for now.

func NewLexer(inputFile string) (*Lexer, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	// File closing is handled by the caller or a wrapper usually.
	// If the lexer is meant to own the file handle, defer file.Close() here.

	reader := bufio.NewReader(file)
	lexer := &Lexer{
		reader:   reader,
		Position: 0,
		line:     0,
		linePos:  1, // Start at column 1
		ch:       0,
		// lines:    make([]string, DEFAULT_ROLLING_BUFFER), // If needed later
	}
	lexer.readChar() // Load the first character
	return lexer, nil
}

func NewLexerFromString(input string) (*Lexer, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	lexer := &Lexer{
		reader:   reader,
		Position: 0,
		line:     0,
		linePos:  0, // Start at column 1
		ch:       0,
		// lines:    make([]string, DEFAULT_ROLLING_BUFFER), // If needed later
	}
	if len(lexer.peekBuffer) > 0 {
		lexer.ch = lexer.peekBuffer[0]
		lexer.peekBuffer = lexer.peekBuffer[1:]
	} else {
		ch, _, err := lexer.reader.ReadRune()
		if err != nil {
			// Handle EOF or read error
			// Check if already EOF to prevent repeated messages/state changes
			if !lexer.eof {
				lexer.ch = 0 // Use 0 to represent EOF
				lexer.eof = true
			} else {
				lexer.ch = 0 // Ensure ch remains 0 if called again after EOF
			}
		}
		lexer.ch = ch
	}

	return lexer, nil
}

// readChar consumes the next rune from the input, updating the lexer's position.
// This is the ONLY function that should advance the main position counters (line, linePos, Position).
func (l *Lexer) readChar() {
	if len(l.peekBuffer) > 0 {
		l.ch = l.peekBuffer[0]
		l.peekBuffer = l.peekBuffer[1:]
	} else {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			// Handle EOF or read error
			// Check if already EOF to prevent repeated messages/state changes
			if !l.eof {
				l.ch = 0 // Use 0 to represent EOF
				l.eof = true
			} else {
				l.ch = 0 // Ensure ch remains 0 if called again after EOF
			}
			return // Return early on error/EOF
		}
		l.ch = ch
	}

	// Update position tracking based on the character *just consumed* (l.ch)
	l.Position++
	if l.ch == '\n' {
		l.line++
		l.linePos = 0 // Reset to column 1 for the new line
	} else {
		l.linePos++ // Increment column position
	}
}

// Helper function to create a token with a single character literal
func newTokenSingle(tokenType TokenType, ch rune) LangToken {
	return LangToken{Type: tokenType, Literal: string(ch)}
}

// Helper function to create a token with a specific literal string
func newTokenLiteral(tokenType TokenType, literal string) LangToken {
	return LangToken{Type: tokenType, Literal: literal}
}

// peekChar reads ahead one character without advancing the main lexer position state.
func (l *Lexer) peekChar() rune {
	if len(l.peekBuffer) == 0 {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			return 0 // Indicate EOF or error during peek
		}
		l.peekBuffer = append(l.peekBuffer, ch)
	}
	// Check again in case read failed or buffer was initially empty
	if len(l.peekBuffer) > 0 {
		return l.peekBuffer[0]
	}
	return 0
}

// peekCharAtIndex reads ahead 'index' characters without advancing the main lexer position state.
// Note: index=0 peeks the next char, index=1 peeks the char after that, etc.
func (l *Lexer) peekCharAtIndex(index int) rune {
	// Ensure peek buffer is filled up to the required index + 1 length
	for len(l.peekBuffer) <= index {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			// Buffer might end up shorter than index if EOF is hit
			break // Stop filling buffer
		}
		l.peekBuffer = append(l.peekBuffer, ch)
	}
	// Check if the requested index is within the bounds of the filled buffer
	if index < len(l.peekBuffer) {
		return l.peekBuffer[index]
	}
	return 0 // Index out of bounds (EOF reached during peeking)
}

// skipWhitespace consumes all subsequent whitespace characters.
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// readIdentifier reads a full identifier or keyword.
func (l *Lexer) readIdentifier() string {
	startPosition := l.Position // Internal tracking if needed, less relevant now
	var identBuilder strings.Builder
	for common.IsLetter(l.ch) || common.IsDigit(l.ch) {
		identBuilder.WriteRune(l.ch)
		l.readChar() // Consume character
	}
	_ = startPosition // Avoid unused variable error if needed later
	// Note: We don't return the final character read, because the loop condition
	// fails *after* reading the first non-identifier character.
	// The main NextToken loop will handle that non-identifier character next.
	return identBuilder.String()
}

// readSingleLineComment consumes characters until newline or EOF.
func (l *Lexer) readSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	// The newline or EOF will be handled by the next NextToken call
}

// readMultiLineComment consumes characters until '*/' or EOF.
func (l *Lexer) readMultiLineComment() error {
	startLine, startPos := l.line, l.linePos
	l.readChar() // Consume initial '*' after '/'

	for {
		if l.ch == 0 { // Check for EOF
			return fmt.Errorf("unterminated multi-line comment starting at line %d, pos %d", startLine+1, startPos-1) // Adjust position slightly
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // Consume '*'
			l.readChar() // Consume '/'
			return nil   // Successfully consumed comment
		}
		l.readChar() // Consume character inside comment
	}
}

// readSpecificKeyword consumes the exact keyword string and returns it.
// Assumes the lexer is positioned at the start of the keyword.
// It's generally better to use readIdentifier and check the map, but this
// can be used if complex lookahead identified the keyword.
func (l *Lexer) readSpecificKeyword(keyword string) string {
	for _, expectedChar := range keyword {
		if l.ch == expectedChar {
			l.readChar()
		} else {
			// This should not happen if called correctly
			break
		}
	}
	return keyword
}

// NextToken reads the next token from the input stream.
func (l *Lexer) NextToken() (LangToken, error) {
	var tok LangToken
	var err error // To capture errors from helpers

	l.skipWhitespace()

	// ---> CAPTURE START POSITION <---
	startLine := l.line
	startPos := l.linePos // Record position AFTER whitespace skip

	// Flag to indicate if the main loop needs to call readChar after processing the token
	advanceChar := true

	switch l.ch {
	case 0: // Handling EOF
		tok.Type = TokenTypeEOF
		tok.Literal = ""
		advanceChar = false // No more characters to read
		// Don't set l.eof here, let the caller check token type

	// --- Operators and Delimiters ---
	case '=':
		if l.peekChar() == '>' { // => (Currently TokenTypeArrow, maybe change?)
			l.readChar() // consume '>'
			tok = newTokenLiteral(TokenTypeArrow, "=>")
			advanceChar = true // Need to read char after '>'
		} else if l.peekChar() == '=' { // == (Equality)
			l.readChar() // consume '='
			tok = newTokenLiteral(TokenTypeEqual, "==")
			advanceChar = true
		} else { // = (Assignment)
			tok = newTokenSingle(TokenTypeAssignment, l.ch)
		}
	case '+':
		tok = newTokenSingle(TokenTypePlus, l.ch)
	case '-':
		if l.peekChar() == '>' { // -> (Lambda Arrow)
			l.readChar() // consume '>'
			tok = newTokenLiteral(TokenTypeLambdaArrow, "->")
			advanceChar = true // Need to read char after '>'
		} else { // - (Minus)
			tok = newTokenSingle(TokenTypeMinus, l.ch)
		}
	case '*':
		tok = newTokenSingle(TokenTypeMultiply, l.ch)
	case '/':
		if l.peekChar() == '/' { // Single-line comment
			l.readChar() // consume second '/'
			l.readSingleLineComment()
			return l.NextToken() // Tail recursion: get the *next* actual token
		} else if l.peekChar() == '*' { // Multi-line comment
			l.readChar() // consume '*'
			commentErr := l.readMultiLineComment()
			if commentErr != nil { // Handle unterminated comment error
				tok.Type = TokenTypeUndefined
				tok.Literal = "/*..." // Indicate unterminated comment
				tok.Line = startLine
				tok.Pos = startPos
				tok.Length = 2      // Length of the starting '/*'
				err = commentErr    // Assign the error
				advanceChar = false // Already advanced past error point
				// No return here, finish processing token below
			} else {
				return l.NextToken() // Tail recursion if comment closed successfully
			}
		} else { // / (Divide)
			tok = newTokenSingle(TokenTypeDivide, l.ch)
		}
	case '<':
		if l.peekChar() == '=' { // <= (Less Than Equal)
			l.readChar() // consume '='
			tok = newTokenLiteral(TokenTypeLessThanEqual, "<=")
			advanceChar = true
		} else { // < (Less Than)
			tok = newTokenSingle(TokenTypeLessThan, l.ch)
		}
	case '>':
		// TODO: Check for >=
		tok = newTokenSingle(TokenTypeGreaterThan, l.ch)
	case '(':
		tok = newTokenSingle(TokenTypeLeftParenthesis, l.ch)
	case ')':
		tok = newTokenSingle(TokenTypeRightParenthesis, l.ch)
	case '{':
		tok = newTokenSingle(TokenTypeLeftBrace, l.ch)
	case '}':
		tok = newTokenSingle(TokenTypeRightBrace, l.ch)
	case '[':
		tok = newTokenSingle(TokenTypeLeftBracket, l.ch)
	case ']':
		tok = newTokenSingle(TokenTypeRightBracket, l.ch)
	case ',':
		tok = newTokenSingle(TokenTypeComma, l.ch)
	case ';':
		tok = newTokenSingle(TokenTypeSemicolon, l.ch)
	case ':':
		tok = newTokenSingle(TokenTypeColon, l.ch)
	case '?':
		tok = newTokenSingle(TokenTypeQuestionMark, l.ch)
	case '.':
		tok = newTokenSingle(TokenTypeDot, l.ch)

	// --- Literals ---
	case '"', '`', '\'': // String literals
		literal, strErr := l.readString()
		// Assign even if error occurred, to get partial literal/position
		tok.Type = TokenTypeString
		tok.Literal = literal
		advanceChar = false // readString already advanced past the end quote or error point
		if strErr != nil {
			tok.Type = TokenTypeUndefined // Mark as undefined on error
			err = strErr                  // Capture the error
			// Length will be calculated below based on partial literal
		}

	default:
		if common.IsLetter(l.ch) { // Identifier or Keyword
			literal := l.readIdentifier()
			tok.Type = LookupIdent(literal) // Check if it's a keyword
			tok.Literal = literal
			advanceChar = false // readIdentifier advanced past the token
		} else if common.IsDigit(l.ch) { // Number literal
			literal := l.readNumber()
			tok.Type = TokenTypeNumber
			tok.Literal = literal
			advanceChar = false // readNumber advanced past the token
		} else {
			// Illegal character
			err = fmt.Errorf("unexpected character: %q at line %d, pos %d", l.ch, startLine+1, startPos)
			tok = newTokenSingle(TokenTypeUndefined, l.ch)
			// Keep advanceChar=true to move past the illegal char (handled below)
		}
	} // end switch

	// Only advance if needed and not already at EOF caused by read errors inside switch
	if advanceChar && l.ch != 0 {
		l.readChar()
	}

	// Handle explicit EOF check *after* potential advance
	if l.ch == 0 && tok.Type == TokenTypeUndefined && err == nil {
		// If we ended up here after the switch with no token type assigned
		// and we are at EOF, emit EOF. Avoids spurious Undefined tokens.
		// Check if we already processed EOF in the switch
		if !l.eof {
			tok.Type = TokenTypeEOF
			tok.Literal = ""
			l.eof = true
		}
	} else if tok.Type == TokenTypeEOF {
		l.eof = true // Ensure EOF flag is set if type is EOF
	}

	// ---> SET TOKEN POSITION AND LENGTH <---
	tok.Line = startLine
	tok.Pos = startPos
	tok.Length = utf8.RuneCountInString(tok.Literal) // Calculate rune length

	// Check if we should return EOF error
	if l.eof && tok.Type != TokenTypeEOF {
		// If we generated a token but *then* hit EOF, the token is valid,
		// but the *next* call should signal the EOF error state.
		// However, if we are returning EOF itself, no error.
	} else if l.eof && tok.Type == TokenTypeEOF && err == nil {
		// Return EOF token without error, but signal state for next call.
		// The error return is primarily for *lexing* errors (bad char, unterminated).
		// Let caller handle repeated EOF token checking if needed.
		// But standard practice is often to return an error on subsequent calls.
		// Let's add the error return on subsequent calls for clarity.
		// We need state to track if EOF was *already* returned.
		// Simplification: Let's allow multiple EOF returns without error for now.
		// Caller can check type.
	}

	return tok, err
}
