package ast

import "compiler/lexer"

type ClassDeclaration struct {
	Token       lexer.LangToken // The 'type' or identifier token
	Name        *Identifier
	Members     []*ClassMember
	LambdaStyle bool
}

func (cd *ClassDeclaration) expressionNode()      {}
func (cd *ClassDeclaration) TokenLiteral() string { return cd.Token.Literal }

type CallExpression struct {
	Token     lexer.LangToken // The '(' token
	Function  ExpressionNode  // The function being called
	Arguments []ExpressionNode
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }

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

type MethodDeclaration struct {
	Token      lexer.LangToken // The identifier token
	ReturnType *Identifier
	Name       *Identifier
	Parameters []*Parameter
	Body       ExpressionNode
}

func (md *MethodDeclaration) expressionNode()      {}
func (md *MethodDeclaration) TokenLiteral() string { return md.Token.Literal }

type Parameter struct {
	Token lexer.LangToken // The identifier token
	Name  *Identifier
	Type  *Identifier
}

func (p *Parameter) expressionNode()      {}
func (p *Parameter) TokenLiteral() string { return p.Token.Literal }
