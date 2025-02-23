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
	return le.StringIndent(0)
}

func (le *LambdaExpression) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	var params []string
	for _, param := range le.Parameters {
		params = append(params, param.String())
	}
	var out strings.Builder
	out.WriteString("(" + strings.Join(params, ", ") + ") -> ")

	// if the body is a block, print it on a new line with an extra indent
	if block, ok := le.Body.(*BlockStatement); ok {
		out.WriteString("\n")
		out.WriteString(block.StringIndent(indent + 1))
	} else if le.Body != nil {
		out.WriteString(le.Body.String())
	}
	return indentStr + out.String()
}
