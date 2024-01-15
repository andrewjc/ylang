package ast

import (
	. "compiler/lexer"
)

type Node interface {
	TokenLiteral() string
}

type Identifier struct {
	Token LangToken // The token.IDENT token
	Value string
}

type ExpressionNode interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type Statement interface {
	Node
	statementNode()
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type NumberLiteral struct {
	Token LangToken
	Value float64
}

func (nl *NumberLiteral) expressionNode()      {}
func (nl *NumberLiteral) TokenLiteral() string { return nl.Token.Literal }

type StringLiteral struct {
	Token LangToken
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }

type InfixExpression struct {
	Token    LangToken
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
