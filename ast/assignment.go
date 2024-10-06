package ast

import (
	. "compiler/lexer"
)

// AssignmentExpression represents an assignment operation (e.g., x = 5)
type AssignmentExpression struct {
	Token    LangToken // The '=' token
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

// Ensure AssignmentExpression implements ExpressionNode
func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	return ae.Left.String() + " " + ae.Operator + " " + ae.Right.String()
}

func (ae *AssignmentExpression) Visit(v Visitor) error {
	return v.VisitAssignmentExpression(ae)
}
