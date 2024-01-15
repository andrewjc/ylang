package ast

import "compiler/lexer"

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

func (vd *VariableDeclaration) statementNode()       {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Literal }

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
