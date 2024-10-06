package ast

import (
	"compiler/lexer"
	"strings"
)

type BlockStatement struct {
	Token      lexer.LangToken // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}
