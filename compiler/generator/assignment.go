package generator

import "compiler/ast"

func (cg *CodeGenerator) VisitAssignmentExpression(ae *ast.AssignmentExpression) error {
	// Evaluate the left side in "LHS mode" so we get the address, not the loaded value.
	cg.inAssignmentLHS = true
	err := ae.Left.Accept(cg)
	cg.inAssignmentLHS = false
	if err != nil {
		return err
	}
	lhsAddr := cg.lastValue

	// Now evaluate the RHS in normal mode (we want the loaded result).
	if err := ae.Right.Accept(cg); err != nil {
		return err
	}
	rhsVal := cg.lastValue

	// Do the store.
	cg.Block.NewStore(rhsVal, lhsAddr)

	// The "value" of an assignment expression is the RHS.
	cg.lastValue = rhsVal
	return nil
}
