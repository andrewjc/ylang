package ast

import "compiler/lexer"

type BlockStatement struct {
	Token      lexer.LangToken // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
