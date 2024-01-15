package lexer

import (
	"compiler/common"
	"strings"
)

func (l *Lexer) readNumber() string {
	var numBuilder strings.Builder
	hasDecimal := false

	for common.IsDigit(l.ch) || (l.ch == '.' && !hasDecimal) {
		if l.ch == '.' {
			hasDecimal = true
		}
		numBuilder.WriteRune(l.ch)
		l.readChar()
	}

	if hasDecimal && numBuilder.Len() == 1 { // Handle single '.' as an error or different token
		return "." // Or handle appropriately, e.g., return an error or a different token type
	}

	return numBuilder.String()
}
