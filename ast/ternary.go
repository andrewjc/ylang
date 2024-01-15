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

// ArrowStyleTernaryExpression represents the arrow-style ternary expression.
type LambdaStyleTernaryExpression struct {
	Token     LangToken      // The '->' token
	Condition ExpressionNode // The condition expression
	TrueExpr  ExpressionNode // The expression if the condition is true
	FalseExpr ExpressionNode // The expression if the condition is false
}

func (aste *LambdaStyleTernaryExpression) expressionNode()      {}
func (aste *LambdaStyleTernaryExpression) TokenLiteral() string { return aste.Token.Literal }

type InlineIfElseTernaryExpression struct {
	Token     LangToken      // The 'if' token
	Condition ExpressionNode // The condition expression
	TrueExpr  ExpressionNode // The expression if the condition is true
	FalseExpr ExpressionNode // The expression if the condition is false
}

func (iite *InlineIfElseTernaryExpression) expressionNode()      {}
func (iite *InlineIfElseTernaryExpression) TokenLiteral() string { return iite.Token.Literal }
