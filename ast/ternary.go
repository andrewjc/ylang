package ast

import . "compiler/lexer"

type TraditionalTernaryExpression struct {
	Token     LangToken      // The '?' token
	Condition ExpressionNode // The condition expression
	TrueExpr  ExpressionNode // The expression if the condition is true
	FalseExpr ExpressionNode // The expression if the condition is false
}

func (tte *TraditionalTernaryExpression) expressionNode()      {}
func (tte *TraditionalTernaryExpression) TokenLiteral() string { return tte.Token.Literal }
func (tte *TraditionalTernaryExpression) String() string {
	return tte.Condition.String() + " ? " + tte.TrueExpr.String() + " : " + tte.FalseExpr.String()
}

// ArrowStyleTernaryExpression represents the arrow-style ternary expression.
type LambdaStyleTernaryExpression struct {
	Token     LangToken      // The '->' token
	Condition ExpressionNode // The condition expression
	TrueExpr  ExpressionNode // The expression if the condition is true
	FalseExpr ExpressionNode // The expression if the condition is false
}

func (aste *LambdaStyleTernaryExpression) expressionNode()      {}
func (aste *LambdaStyleTernaryExpression) TokenLiteral() string { return aste.Token.Literal }
func (aste *LambdaStyleTernaryExpression) String() string {
	return aste.Condition.String() + " -> " + aste.TrueExpr.String() + " : " + aste.FalseExpr.String()
}

type InlineIfElseTernaryExpression struct {
	Token     LangToken      // The 'if' token
	Condition ExpressionNode // The condition expression
	TrueExpr  ExpressionNode // The expression if the condition is true
	FalseExpr ExpressionNode // The expression if the condition is false
}

func (iite *InlineIfElseTernaryExpression) expressionNode()      {}
func (iite *InlineIfElseTernaryExpression) TokenLiteral() string { return iite.Token.Literal }
func (iite *InlineIfElseTernaryExpression) String() string {
	return iite.Condition.String() + " if " + iite.TrueExpr.String() + " else " + iite.FalseExpr.String()
}
