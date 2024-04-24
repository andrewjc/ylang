package ast

import "compiler/lexer"

type Comment struct {
	Token lexer.LangToken // The '//' or '/*' token
	Text  string
}

func (c *Comment) expressionNode()      {}
func (c *Comment) TokenLiteral() string { return c.Token.Literal }
