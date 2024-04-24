package ast

import "compiler/lexer"

type LambdaExpression struct {
	Token      lexer.LangToken // the TokenTypeLeftParenthesis token
	Parameters []*Identifier
	Body       ExpressionNode
}

func (le *LambdaExpression) expressionNode()      {}
func (le *LambdaExpression) TokenLiteral() string { return le.Token.Literal }
