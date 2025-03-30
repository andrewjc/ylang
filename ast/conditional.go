package ast

import (
	. "compiler/lexer"
	"strings"
)

// IfStatement represents an 'if' control flow statement.
type IfStatement struct {
	Token       LangToken // The 'if' token
	Condition   ExpressionNode
	Consequence ExpressionNode
	Alternative ExpressionNode
}

func (is *IfStatement) Accept(visitor Visitor) error {
	return visitor.VisitIfStatement(is)
}

func (is *IfStatement) expressionNode() {
	//TODO implement me
	panic("implement me")
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	return is.StringIndent(0) // Use StringIndent for consistency
}

func (is *IfStatement) StringIndent(indent int) string {
	indentStr := strings.Repeat("    ", indent)
	var out strings.Builder

	out.WriteString(indentStr + "if ")
	out.WriteString(is.Condition.String())
	out.WriteString(" ")

	// Handle consequence formatting (might be Block or single expression)
	if block, ok := is.Consequence.(*BlockStatement); ok {
		out.WriteString(block.StringIndent(indent)) // Use block's indenting
	} else if is.Consequence != nil {
		out.WriteString(is.Consequence.String()) // Assume single line if not block
	} else {
		out.WriteString("{}") // Empty consequence
	}

	if is.Alternative != nil {
		// Add newline and indent before 'else' if consequence was a block
		_, consIsBlock := is.Consequence.(*BlockStatement)
		if consIsBlock {
			out.WriteString("\n" + indentStr)
		} else {
			out.WriteString(" ") // Space if consequence was single line
		}
		out.WriteString("else ")

		// Handle alternative formatting
		if elseIf, ok := is.Alternative.(*IfStatement); ok {
			// 'else if' - print 'if' part directly without extra indent/braces
			out.WriteString(elseIf.StringIndent(indent)) // Pass current indent
		} else if block, ok := is.Alternative.(*BlockStatement); ok {
			out.WriteString(block.StringIndent(indent)) // Use block's indenting
		} else {
			out.WriteString(is.Alternative.String()) // Assume single line
		}
	}

	return out.String()
}
