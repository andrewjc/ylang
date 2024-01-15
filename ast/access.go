package ast

import "compiler/lexer"

type MemberAccessExpression struct {
	Token  lexer.LangToken // The '.' token
	Left   ExpressionNode  // The expression on the left of the dot
	Member string          // The member being accessed
}

func (m *MemberAccessExpression) TokenLiteral() string {
	return m.Token.Literal
}

func (m *MemberAccessExpression) expressionNode() {
	//TODO implement me
	panic("implement me")
}
