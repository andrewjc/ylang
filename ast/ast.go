package ast

import (
	. "compiler/lexer"
	"fmt"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Identifier struct {
	Token LangToken // The token.IDENT token
	Value string
}

type ExpressionNode interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

type NumberLiteral struct {
	Token LangToken
	Value float64
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }
func (nl *NumberLiteral) String() string {
	return fmt.Sprintf("%v", nl.Value)
}

type StringLiteral struct {
	Token LangToken
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Value)
}

type InfixExpression struct {
	Token    LangToken
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}
