package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitIdentifier(id *ast.Identifier) error {
	// Look up the alloca we stored under this name (if any).
	allocaVal := cg.getVar(id.Value)
	if allocaVal == nil {
		// fallback if we never stored anything under this name
		cg.lastValue = constant.NewInt(types.I32, 0)
		return nil
	}

	// If we are on the left side of an assignment, we want the "address" itself.
	if cg.inAssignmentLHS {
		cg.lastValue = allocaVal
		return nil
	}

	// Otherwise, we want to "load" the value from that address.
	allocaInst, ok := allocaVal.(*ir.InstAlloca)
	if !ok {
		return fmt.Errorf("variable %s is not an alloca instruction", id.Value)
	}

	// The type we gave to allocaInst is the "element type" we'll load.
	loaded := cg.Block.NewLoad(allocaInst.ElemType, allocaInst)
	cg.lastValue = loaded
	return nil
}
