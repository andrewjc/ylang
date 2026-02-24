package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/types"
)

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

	// When storing an integer into a pointer whose element type is a narrower
	// (or wider) integer, truncate or sign-extend the RHS so the types match.
	// This is needed for byte-buffer writes such as: buf[0] = digit + 48
	if lhsPtrType, ok := lhsAddr.Type().(*types.PointerType); ok {
		if elemIntType, ok := lhsPtrType.ElemType.(*types.IntType); ok {
			if rhsIntType, ok := rhsVal.Type().(*types.IntType); ok {
				if rhsIntType.BitSize > elemIntType.BitSize {
					rhsVal = cg.Block.NewTrunc(rhsVal, elemIntType)
				} else if rhsIntType.BitSize < elemIntType.BitSize {
					rhsVal = cg.Block.NewSExt(rhsVal, elemIntType)
				}
			}
		}
	}

	// Do the store.
	cg.Block.NewStore(rhsVal, lhsAddr)

	// The "value" of an assignment expression is the RHS.
	cg.lastValue = rhsVal
	return nil
}
