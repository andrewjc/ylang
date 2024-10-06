package ast

import (
	"compiler/lexer"
	"strings"
)

type LambdaExpression struct {
	Token      lexer.LangToken // the TokenTypeLeftParenthesis token
	Parameters []*Identifier
	Body       ExpressionNode
}

func (le *LambdaExpression) expressionNode()      {}
func (le *LambdaExpression) TokenLiteral() string { return le.Token.Literal }
func (le *LambdaExpression) String() string {
	var params []string
	for _, param := range le.Parameters {
		params = append(params, param.String())
	}

	var out strings.Builder
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") -> ")

	if le.Body != nil {
		out.WriteString(le.Body.String())
	}

	return out.String()
}
