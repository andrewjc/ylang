package ast

import "compiler/lexer"

type AssemblyExpression struct {
	Token lexer.LangToken
	Code  *StringLiteral
	Args  []ExpressionNode
}

func (ae *AssemblyExpression) expressionNode()      {}
func (ae *AssemblyExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssemblyExpression) String() string {
	s := "asm(" + ae.Code.String()
	for _, arg := range ae.Args {
		s += ", " + arg.String()
	}
	s += ")"
	return s
}

func (ae *AssemblyExpression) Accept(v Visitor) error {
	visitor, ok := v.(interface {
		VisitAssemblyExpression(*AssemblyExpression) error
	})
	if !ok {
		panic("Visitor does not implement VisitAssemblyExpression")
	}
	return visitor.VisitAssemblyExpression(ae)
}
