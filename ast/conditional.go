package ast

import (
	. "compiler/lexer"
	"strings"
)

// IfStatement represents an 'if' control flow statement.
type IfStatement struct {
	Token       LangToken // The 'if' token
	Condition   ExpressionNode
	Consequence ExpressionNode
	Alternative ExpressionNode
}

func (is *IfStatement) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var out strings.Builder

	out.WriteString("if ")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")
	out.WriteString(is.Consequence.String())

	if is.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(is.Alternative.String())
	}

	return out.String()
}
