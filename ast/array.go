package ast

import (
	"compiler/lexer"
	"strings"
)

type ArrayLiteral struct {
	Token    lexer.LangToken // The first token of the expression
	Elements []ExpressionNode
}

func (es *ArrayLiteral) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (es *ArrayLiteral) statementNode()       {}
func (es *ArrayLiteral) TokenLiteral() string { return es.Token.Literal }
func (al *ArrayLiteral) String() string {
	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}
