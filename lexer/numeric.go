package lexer

import (
	"compiler/common"
	"strings"
)

func (l *Lexer) readNumber() string {
	var numBuilder strings.Builder
	hasDecimal := false

	// Loop while the current character is a digit OR
	// it's the first decimal point encountered AND the *next* character is a digit.
	for common.IsDigit(l.ch) || (l.ch == '.' && !hasDecimal && common.IsDigit(l.peekChar())) {
		if l.ch == '.' {
			hasDecimal = true
		}
		numBuilder.WriteRune(l.ch)
		l.readChar() // Use the corrected lowercase 'r' function name
	}

	// The check for a standalone '.' is less critical now because the loop
	// condition prevents entering the loop for just '.', but it doesn't hurt.
	// A standalone '.' will be handled as TokenTypeDot in the main NextToken switch.
	// if hasDecimal && numBuilder.Len() == 1 {
	//     // This case should ideally not be reached with the corrected loop condition.
	//     // If it were, it would indicate an issue elsewhere.
	// }

	return numBuilder.String()
}
