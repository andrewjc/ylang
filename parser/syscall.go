package parser

import (
	"compiler/ast"
	. "compiler/lexer"
)

func (p *Parser) parseSysCallExpression() ast.ExpressionNode {
	// Current token is "syscall"
	expr := &ast.SyscallExpression{Token: p.currentToken}

	// Next token should be '('
	if !p.expectPeek(TokenTypeLeftParenthesis) {
		return nil
	}

	// parse the expression list, e.g. ( 1, str, len, 0, ... )
	exprList := p.parseExpressionList(TokenTypeRightParenthesis)
	// The first item is the syscall number, the next items are the 6 arguments
	if len(exprList) == 0 {
		// no syscall number -> error
		return nil
	}

	expr.Num = exprList[0]
	if len(exprList) > 1 {
		expr.Args = exprList[1:]
	} else {
		expr.Args = []ast.ExpressionNode{}
	}

	return expr
}
