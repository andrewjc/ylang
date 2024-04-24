package lexer

import "strings"

func (l *Lexer) readString() string {
	var strBuilder strings.Builder

	openingQuoteStyle := l.ch

	l.ReadChar() // Consume the opening double quote
	for l.ch != openingQuoteStyle && l.ch != 0 {
		if l.ch == '\\' {
			l.ReadChar() // Consume the backslash
			switch l.ch {
			case 'n':
				strBuilder.WriteRune('\n')
			case 't':
				strBuilder.WriteRune('\t')
			case 'r':
				strBuilder.WriteRune('\r')
			case '\\':
				strBuilder.WriteRune('\\')
			case '"':
				strBuilder.WriteRune('"')
			default:
				strBuilder.WriteRune('\\')
				strBuilder.WriteRune(l.ch)
			}
		} else {
			strBuilder.WriteRune(l.ch)
		}
		l.ReadChar()
	}
	if l.ch != openingQuoteStyle {
		// TODO Handle the error - unclosed string literal
		return ""
	}

	// Consume the closing quote
	l.ReadChar()

	return strBuilder.String()
}
