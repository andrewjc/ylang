package generator

import (
	"compiler/ast"
)

func (cg *CodeGenerator) VisitExpressionStatement(es *ast.ExpressionStatement) error {
	// Generate code for the expression, discard its value.
	if es.Expression == nil {
		return nil
	}
	return es.Expression.Accept(cg)
}

func (cg *CodeGenerator) VisitAssignmentExpression(ae *ast.AssignmentExpression) error {
	// Evaluate the RHS
	if err := ae.Right.Accept(cg); err != nil {
		return err
	}
	rhs := cg.lastValue

	// Evaluate the LHS to get its address
	if err := ae.Left.Accept(cg); err != nil {
		return err
	}
	addr := cg.lastValue

	if addr != nil && rhs != nil {
		cg.Block.NewStore(rhs, addr)
	}

	// The assignment expression's value is the RHS
	cg.lastValue = rhs
	return nil
}

func (cg *CodeGenerator) VisitIndexExpression(ie *ast.IndexExpression) error {
	// Minimal stub: generate code for the left, ignore the index for now.
	if err := ie.Left.Accept(cg); err != nil {
		return err
	}
	// We are not producing anything real, so just clear or reuse lastValue
	cg.lastValue = nil
	return nil
}
