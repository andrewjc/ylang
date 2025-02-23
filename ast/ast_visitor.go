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

	// ac: todo add more visit methods here
}
