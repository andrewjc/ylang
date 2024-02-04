package parser

import "fmt"

type ParserError struct {
	Message string
	Line    int
	Pos     int
}

func (pe *ParserError) Error() string {
	return fmt.Sprintf("Parser error at line %d, position %d: %s", pe.Line, pe.Pos, pe.Message)
}
