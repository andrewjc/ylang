package ast

type Visitor interface {
	VisitProgram(program *Program) error
	VisitFunctionDefinition(fn *FunctionDefinition) error
	VisitLetStatement(ls *LetStatement) error
	VisitReturnStatement(rs *ReturnStatement) error
	VisitExpressionStatement(es *ExpressionStatement) error
	VisitNumberLiteral(nl *NumberLiteral) error
	VisitStringLiteral(sl *StringLiteral) error
	VisitIdentifier(id *Identifier) error
	VisitInfixExpression(ie *InfixExpression) error
	VisitCallExpression(ce *CallExpression) error
	VisitArrayLiteral(al *ArrayLiteral) error
	VisitLambdaExpression(le *LambdaExpression) error
	VisitMemberAccessExpression(mae *MemberAccessExpression) error
	VisitAssignmentExpression(as *AssignmentExpression) error
	VisitBlockStatement(bs *BlockStatement) error
	VisitIndexExpression(ie *IndexExpression) error
	VisitVariableDeclaration(vd *VariableDeclaration) error
	VisitIfStatement(is *IfStatement) error
	VisitTraditionalTernaryExpression(te *TraditionalTernaryExpression) error
	VisitLambdaStyleTernaryExpression(aste *LambdaStyleTernaryExpression) error
	VisitInlineIfElseTernaryExpression(iite *InlineIfElseTernaryExpression) error
	VisitDotOperator(do *DotOperator) error
	VisitSyscallExpression(se *SyscallExpression) error
	VisitImportStatement(is *ImportStatement) error

	// ac: todo add more visit methods here
}
