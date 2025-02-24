package ast

import (
	"compiler/lexer"
)

// SyscallExpression represents something like:
//
//	syscall(SYS_write, fd, buf, size, ..., ...)
//
// It's an expression, returning the syscall's result (i64).
// This should be wrapped up in stdlib calls like stdout.write(buf, size) to hide the uglyness...
type SyscallExpression struct {
	Token lexer.LangToken  // e.g. 'syscall'
	Num   ExpressionNode   // The syscall number (e.g. 1 for write)
	Args  []ExpressionNode // Up to 6 arguments
}

func (se *SyscallExpression) expressionNode()      {}
func (se *SyscallExpression) TokenLiteral() string { return se.Token.Literal }

func (se *SyscallExpression) String() string {
	// e.g. "syscall(1, arg0, arg1, ...)"
	// For debugging:
	out := "syscall("
	if se.Num != nil {
		out += se.Num.String()
	}
	for i, arg := range se.Args {
		if i == 0 {
			out += ", "
		} else {
			out += ", "
		}
		out += arg.String()
	}
	out += ")"
	return out
}

func (se *SyscallExpression) Accept(v Visitor) error {
	return v.VisitSyscallExpression(se)
}
