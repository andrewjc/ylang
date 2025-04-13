package lexer

import (
	"fmt"
	"strings"
)

func (l *Lexer) readString() (string, error) {
	var strBuilder strings.Builder
	openingQuoteStyle := l.ch
	startLine, startPos := l.line, l.linePos // Record start for error reporting

	l.readChar() // Consume the opening quote

	for l.ch != openingQuoteStyle {
		if l.ch == 0 { // Check for EOF (unterminated string)
			// Return the partially built string along with the error
			return strBuilder.String(), fmt.Errorf("unterminated string literal starting at line %d, pos %d", startLine+1, startPos)
		}

		if l.ch == '\\' { // Handle escape sequences
			l.readChar() // Consume the backslash
			escapeChar := l.ch
			if escapeChar == 0 { // EOF after backslash
				return strBuilder.String(), fmt.Errorf("unterminated string literal starting at line %d, pos %d (EOF after backslash)", startLine+1, startPos)
			}
			switch escapeChar {
			case 'n':
				strBuilder.WriteRune('\n')
			case 't':
				strBuilder.WriteRune('\t')
			case 'r':
				strBuilder.WriteRune('\r')
			case '\\':
				strBuilder.WriteRune('\\')
			case '"': // Allow escaping double quotes
				strBuilder.WriteRune('"')
			case '\'': // Allow escaping single quotes
				strBuilder.WriteRune('\'')
			case '`': // Allow escaping backticks
				strBuilder.WriteRune('`')
			default:
				// Treat unknown escapes as literal backslash + character
				// This maintains the original behavior, though one might choose
				// to error on invalid escape sequences depending on language spec.
				strBuilder.WriteRune('\\')
				strBuilder.WriteRune(escapeChar)
			}
		} else {
			// Regular character in the string
			strBuilder.WriteRune(l.ch)
		}
		l.readChar() // Consume character inside string (or the second char of escape sequence)
	}

	// If the loop finished, we should be at the closing quote.
	// The check `l.ch != openingQuoteStyle` inside the loop already guarantees this
	// unless EOF was hit (which is handled).

	// Consume the closing quote
	l.readChar()

	return strBuilder.String(), nil // Return the complete string and nil error
}
