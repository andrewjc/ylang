package ast

import (
	"compiler/lexer"
	"strings"
)

// VariableDeclaration represents a variable declaration statement.
type VariableDeclaration struct {
	Token lexer.LangToken // The 'let' token
	Name  *Identifier
	Type  *Identifier // Optional type name
	Value ExpressionNode
}

func (vd *VariableDeclaration) expressionNode() {
	panic("implement me")
}
func (vd *VariableDeclaration) Accept(v Visitor) error {
	return v.VisitVariableDeclaration(vd)
}
func (vd *VariableDeclaration) statementNode()       {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Literal }
func (vd *VariableDeclaration) String() string {
	var out strings.Builder

	out.WriteString("let ")
	out.WriteString(vd.Name.String())

	// Include type annotation if present
	if vd.Type != nil {
		out.WriteString(": ")
		out.WriteString(vd.Type.String())
	}

	if vd.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vd.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// AssignmentStatement represents an assignment statement.
type AssignmentStatement struct {
	Token lexer.LangToken // The '=' token
	Name  *Identifier
	Value ExpressionNode
}

func (as *AssignmentStatement) expressionNode() {
	panic("implement me")
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
