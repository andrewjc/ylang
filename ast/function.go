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
func (fd *FunctionDefinition) Accept(v Visitor) error {
	return v.VisitFunctionDefinition(fd)
}
func (fd *FunctionDefinition) String() string {
	return fd.StringIndent(0)
}

func (fd *FunctionDefinition) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	var params []string
	for _, param := range fd.Parameters {
		params = append(params, param.String())
	}

	var out strings.Builder
	// function header
	if fd.Name != nil {
		out.WriteString(fd.Name.String())
	} else {
		out.WriteString("anonymous")
	}
	out.WriteString("(" + strings.Join(params, ", ") + ") -> ")

	// if the body is a block, then insert a newline and indent it
	if block, ok := fd.Body.(*BlockStatement); ok {
		out.WriteString("\n")
		out.WriteString(block.StringIndent(indent))
	} else if fd.Body != nil {
		out.WriteString(fd.Body.String())
	}
	return indentStr + out.String()
}
