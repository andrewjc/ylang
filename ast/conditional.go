package ast

import . "compiler/lexer"

// IfStatement represents an 'if' control flow statement.
type IfStatement struct {
	Token       LangToken // The 'if' token
	Condition   ExpressionNode
	Consequence ExpressionNode
	Alternative ExpressionNode
}

func (is *IfStatement) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
