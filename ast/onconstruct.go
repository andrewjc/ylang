package ast

import "compiler/lexer"

type OnConstructStatement struct {
	Token  lexer.LangToken // The 'onConstruct' token
	Lambda *LambdaExpression
}

func (oc *OnConstructStatement) expressionNode()      {}
func (oc *OnConstructStatement) TokenLiteral() string { return oc.Token.Literal }
