package ast

import "compiler/lexer"

type SwitchStatement struct {
	Token       lexer.LangToken // The 'switch' token
	Expression  ExpressionNode
	Cases       []*SwitchCase
	DefaultCase *SwitchCase
}

func (ss *SwitchStatement) expressionNode()      {}
func (ss *SwitchStatement) TokenLiteral() string { return ss.Token.Literal }

type SwitchCase struct {
	Token      lexer.LangToken // The 'case' or 'default' token
	Expression ExpressionNode
	Block      ExpressionNode
}

func (sc *SwitchCase) expressionNode()      {}
func (sc *SwitchCase) TokenLiteral() string { return sc.Token.Literal }
