package ast

import "compiler/lexer"

type AssemblyStatement struct {
	Token lexer.LangToken // The 'asm' token
	Code  *StringLiteral
}

func (as *AssemblyStatement) expressionNode()      {}
func (as *AssemblyStatement) TokenLiteral() string { return as.Token.Literal }
