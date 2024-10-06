package ast

import (
	"compiler/lexer"
	"strings"
)

type ClassDeclaration struct {
	Token       lexer.LangToken // The 'type' or identifier token
	Name        *Identifier
	Members     []*ClassMember
	LambdaStyle bool
}

func (cd *ClassDeclaration) expressionNode()      {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }
func (cd *ClassDeclaration) String() string {
	var members []string
	for _, member := range cd.Members {
		members = append(members, member.String())
	}
	return cd.Name.String() + " {" + strings.Join(members, " ") + "}"
}

type CallExpression struct {
	Token     lexer.LangToken // The '(' token
	Function  ExpressionNode  // The function being called
	Arguments []ExpressionNode
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var args []string
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}

type Field struct {
	Token lexer.LangToken // The 'let' token
	Name  *Identifier
}

func (f *Field) expressionNode()      {}
func (f *Field) TokenLiteral() string { return f.Token.Literal }

type ClassMember struct {
	VariableDeclaration *VariableDeclaration
	MethodDeclaration   *MethodDeclaration
}

func (cm *ClassMember) expressionNode() {}
func (cm *ClassMember) TokenLiteral() string {
	if cm.VariableDeclaration != nil {
		return cm.VariableDeclaration.TokenLiteral()
	}
	if cm.MethodDeclaration != nil {
		return cm.MethodDeclaration.TokenLiteral()
	}
	return ""
}
func (cm *ClassMember) String() string {
	if cm.VariableDeclaration != nil {
		return cm.VariableDeclaration.String()
	}
	if cm.MethodDeclaration != nil {
		return cm.MethodDeclaration.String()
	}
	return ""
}

type MethodDeclaration struct {
	Token      lexer.LangToken // The identifier token
	ReturnType *Identifier
	Name       *Identifier
	Parameters []*Parameter
	Body       ExpressionNode
}

func (md *MethodDeclaration) expressionNode()      {}
func (md *MethodDeclaration) TokenLiteral() string { return md.Token.Literal }
func (md *MethodDeclaration) String() string {
	var params []string
	for _, param := range md.Parameters {
		params = append(params, param.String())
	}
	return md.ReturnType.String() + " " + md.Name.String() + "(" + strings.Join(params, ", ") + ") " + md.Body.String()
}

type Parameter struct {
	Token lexer.LangToken // The identifier token
	Name  *Identifier
	Type  *Identifier
}

func (p *Parameter) expressionNode()      {}
func (p *Parameter) TokenLiteral() string { return p.Token.Literal }
func (p *Parameter) String() string {
	return p.Type.String() + " " + p.Name.String()
}
