package ast

import "compiler/lexer"

type OnDestructStatement struct {
	Token  lexer.LangToken // The 'onDestruct' token
	Lambda *LambdaExpression
}

func (od *OnDestructStatement) expressionNode()      {}
func (od *OnDestructStatement) TokenLiteral() string { return od.Token.Literal }
