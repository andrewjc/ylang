package ast

import "compiler/lexer"

type FunctionDefinition struct {
	Token      lexer.LangToken // The first token of the expression
	Name       *Identifier
	Expression ExpressionNode
	Parameters []*Identifier
	Body       ExpressionNode
}

func (es *FunctionDefinition) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (es *FunctionDefinition) statementNode()       {}
func (es *FunctionDefinition) TokenLiteral() string { return es.Token.Literal }
