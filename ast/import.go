package ast

import (
	"compiler/lexer"
)

type ImportStatement struct {
	Token lexer.LangToken
	Path  string
}

func (is *ImportStatement) statementNode()       {}
func (is *ImportStatement) TokenLiteral() string { return is.Token.Literal }
func (is *ImportStatement) String() string       { return "import \"" + is.Path + "\"" }

func (is *ImportStatement) Accept(v Visitor) error {
	return v.VisitImportStatement(is)
}
