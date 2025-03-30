package ast

import (
	"compiler/lexer"
	"strings"
)

type BlockStatement struct {
	Token      lexer.LangToken // The '{' token
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) Accept(v Visitor) error {
	visitor, ok := v.(interface {
		VisitBlockStatement(*BlockStatement) error
	})
	if !ok {
		panic("Visitor does not implement VisitBlockStatement") // Or return error
	}
	return visitor.VisitBlockStatement(bs)
}
func (bs *BlockStatement) String() string {
	return bs.StringIndent(0)
}

func (bs *BlockStatement) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	var out strings.Builder
	out.WriteString("{\n")
	for _, stmt := range bs.Statements {
		stmtStr := ""
		if s, ok := stmt.(interface{ StringIndent(int) string }); ok {
			stmtStr = s.StringIndent(indent + 1)
		} else {
			stmtStr = strings.Repeat("    ", indent+1) + stmt.String()
		}
		out.WriteString(stmtStr)
		switch stmt.(type) {
		case *LetStatement, *ReturnStatement, *ExpressionStatement:
			if !strings.HasSuffix(stmtStr, ";") && !strings.HasSuffix(stmtStr, "}") {
				out.WriteString(";")
			}
		}
		out.WriteString("\n")
	}
	out.WriteString(indentStr + "}")
	return out.String()
}
