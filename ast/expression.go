package ast

import "compiler/lexer"

type ExpressionStatement struct {
	Token      lexer.LangToken // The first token of the expression
	Expression ExpressionNode
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
