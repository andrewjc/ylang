package ast

import "compiler/lexer"

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
