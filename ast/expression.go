package ast

import (
	"compiler/lexer"
	"strings"
)

type LetStatement struct {
	Token lexer.LangToken // the TokenTypeLet token
	Name  *Identifier
	Value ExpressionNode
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out strings.Builder

	out.WriteString("let ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token       lexer.LangToken // the TokenTypeReturn token
	ReturnValue ExpressionNode
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString("return ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      lexer.LangToken // the first token of the expression
	Expression ExpressionNode
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// DotOperator
type DotOperator struct {
	Token lexer.LangToken // The '.' token
	Left  ExpressionNode
	Right *Identifier
}

func (do *DotOperator) expressionNode()      {}
func (do *DotOperator) TokenLiteral() string { return do.Token.Literal }
func (do *DotOperator) String() string {
	return do.Left.String() + "." + do.Right.String()
}
