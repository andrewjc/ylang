package ast

import (
	"compiler/lexer"
	"strings"
)

type FunctionDefinition struct {
	Token      lexer.LangToken // The first token of the expression
	Name       *Identifier
	Expression ExpressionNode
	Parameters []*Identifier
	Body       ExpressionNode
	ReturnType *Identifier
}

func (es *FunctionDefinition) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (es *FunctionDefinition) statementNode()       {}
func (es *FunctionDefinition) TokenLiteral() string { return es.Token.Literal }

func (fd *FunctionDefinition) String() string {
	var params []string
	for _, param := range fd.Parameters {
		params = append(params, param.String())
	}

	var out strings.Builder
	if fd.Name != nil {
		out.WriteString(fd.Name.String())
	} else {
		out.WriteString("anonymous")
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	if fd.ReturnType != nil {
		out.WriteString(": ")
		out.WriteString(fd.ReturnType.String())
	}

	out.WriteString("-> ")

	if fd.Body != nil {
		out.WriteString(fd.Body.String())
	}

	return out.String()
}
