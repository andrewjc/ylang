package ast

import "compiler/lexer"

type LetStatement struct {
	Token lexer.LangToken // the TokenTypeLet token
	Name  *Identifier
	Value ExpressionNode
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       lexer.LangToken // the TokenTypeReturn token
	ReturnValue ExpressionNode
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

type ExpressionStatement struct {
	Token      lexer.LangToken // the first token of the expression
	Expression ExpressionNode
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// DotOperator
type DotOperator struct {
	Token lexer.LangToken // The '.' token
	Left  ExpressionNode
	Right *Identifier
}

func (do *DotOperator) expressionNode()      {}
func (do *DotOperator) TokenLiteral() string { return do.Token.Literal }
