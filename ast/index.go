package ast

import "compiler/lexer"

type IndexExpression struct {
	Token lexer.LangToken // The '[' token
	Left  ExpressionNode
	Index ExpressionNode
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}
