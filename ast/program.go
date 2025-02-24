package ast

import "strings"

type Program struct {
	MainFunction      *FunctionDefinition
	ClassDeclarations []*ClassDeclaration
	Functions         []*FunctionDefinition
	DataStructures    []*DataStructure
	ImportStatements  []*ImportStatement
}

func (p *Program) String() string {
	var out strings.Builder
	// Print main function if available
	if p.MainFunction != nil {
		out.WriteString(p.MainFunction.String())
		out.WriteString("\n")
	}
	return out.String()
}
