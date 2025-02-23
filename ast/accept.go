package ast

// Accept allows a Visitor to visit the node.
func (p *Program) Accept(v Visitor) error {
	return v.VisitProgram(p)
}

func (fd *FunctionDefinition) Visit(v Visitor) error {
	return v.VisitFunctionDefinition(fd)
}

func (ls *LetStatement) Accept(v Visitor) error {
	return v.VisitLetStatement(ls)
}

func (rs *ReturnStatement) Accept(v Visitor) error {
	return v.VisitReturnStatement(rs)
}

func (es *ExpressionStatement) Accept(v Visitor) error {
	return v.VisitExpressionStatement(es)
}

func (nl *NumberLiteral) Accept(v Visitor) error {
	return v.VisitNumberLiteral(nl)
}

func (sl *StringLiteral) Accept(v Visitor) error {
	return v.VisitStringLiteral(sl)
}

func (id *Identifier) Accept(v Visitor) error {
	return v.VisitIdentifier(id)
}

func (ie *InfixExpression) Accept(v Visitor) error {
	return v.VisitInfixExpression(ie)
}

func (ce *CallExpression) Accept(v Visitor) error {
	return v.VisitCallExpression(ce)
}

func (al *ArrayLiteral) Accept(v Visitor) error {
	return v.VisitArrayLiteral(al)
}

func (le *LambdaExpression) Accept(v Visitor) error {
	return v.VisitLambdaExpression(le)
}

func (mae *MemberAccessExpression) Accept(v Visitor) error {
	return v.VisitMemberAccessExpression(mae)
}

func (as *AssignmentExpression) Accept(v Visitor) error {
	return v.VisitAssignmentExpression(as)
}
