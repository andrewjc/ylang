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
func (bs *BlockStatement) Accept(v Visitor) error {
	return v.VisitBlockStatement(bs)
}
func (bs *BlockStatement) String() string {
	return bs.StringIndent(0)
}

func (bs *BlockStatement) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	var out strings.Builder

	// opening brace on its own line
	out.WriteString(indentStr + "{\n")
	// Each statement printed with one more indent level.
	for _, stmt := range bs.Statements {
		if s, ok := stmt.(interface{ StringIndent(int) string }); ok {
			out.WriteString(s.StringIndent(indent + 1))
		} else {
			out.WriteString(strings.Repeat("    ", indent+1) + stmt.String())
		}
		out.WriteString("\n")
	}
	out.WriteString(indentStr + "}")
	return out.String()
}
