package parser

import (
	"fmt"
	"strings"
)

type ParserError struct {
	Line         int
	Pos          int
	Message      string
	CodeFragment string
}

func (e *ParserError) Error() string {
	strippedCodeFragment := strings.TrimSpace(e.CodeFragment)
	return fmt.Sprintf("Parse error at line %d, position %d: %s\n\t%s", e.Line, e.Pos, e.Message, strippedCodeFragment)
}
