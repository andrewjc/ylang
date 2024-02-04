package ast

import "compiler/lexer"

type FunctionDefinition struct {
	Token      lexer.LangToken // The first token of the expression
	Expression ExpressionNode
	Parameters []string
	Body       *BlockStatement
}

func (es *FunctionDefinition) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (es *FunctionDefinition) statementNode()       {}
func (es *FunctionDefinition) TokenLiteral() string { return es.Token.Literal }
