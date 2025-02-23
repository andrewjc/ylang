package ast

import (
	"compiler/lexer"
	"strings"
)

type DataStructure struct {
	Token  lexer.LangToken // The 'data' token
	Name   *Identifier
	Fields []*Field
	Style  DataStructureStyle
}

type DataStructureStyle int

const (
	DataStructureStyleBraces DataStructureStyle = iota
	DataStructureStyleEquals
	DataStructureStyleColon
	DataStructureStyleTupleLike
)

func (ds *DataStructure) expressionNode()      {}
func (ds *DataStructure) TokenLiteral() string { return ds.Token.Literal }
func (ds *DataStructure) String() string {
	var out strings.Builder
	out.WriteString("DataStructure: ")
	out.WriteString(ds.Name.String()) // Assuming Name is a field that has its own String method
	// Add more fields or details as needed, e.g., ds.Fields
	return out.String()
}
