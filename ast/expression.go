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
	return ls.StringIndent(0)
}

func (ls *LetStatement) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	var out strings.Builder
	out.WriteString("let " + ls.Name.String() + " = " + ls.Value.String() + ";")
	return indentStr + out.String()
}

type ReturnStatement struct {
	Token       lexer.LangToken // the TokenTypeReturn token
	ReturnValue ExpressionNode
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	return rs.StringIndent(0)
}

func (rs *ReturnStatement) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	var out strings.Builder
	out.WriteString("return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return indentStr + out.String()
}

type ExpressionStatement struct {
	Token      lexer.LangToken // the first token of the expression
	Expression ExpressionNode
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	return es.StringIndent(0)
}

func (es *ExpressionStatement) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	return indentStr + es.Expression.String()
}

// DotOperator
type DotOperator struct {
	Token lexer.LangToken // The '.' token
	Left  ExpressionNode
	Right *Identifier
}

func (do *DotOperator) Accept(visitor Visitor) error {
	return visitor.VisitDotOperator(do)
}

func (do *DotOperator) expressionNode()      {}
func (do *DotOperator) TokenLiteral() string { return do.Token.Literal }
func (do *DotOperator) String() string {
	return do.Left.String() + "." + do.Right.String()
}
