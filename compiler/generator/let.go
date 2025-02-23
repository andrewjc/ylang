package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitLetStatement(ls *ast.LetStatement) error {
	// Allocate in current block (or ideally in function entry)
	alloca := cg.Block.NewAlloca(types.I32)
	var valueVal = alloca

	if ls.Value != nil {
		if err := ls.Value.Accept(cg); err != nil {
			return err
		}
		initVal := cg.lastValue
		if initVal == nil {
			initVal = constant.NewInt(types.I32, 0)
		}
		cg.Block.NewStore(initVal, alloca)
	}
	cg.setVar(ls.Name.Value, valueVal)
	return nil
}

func (cg *CodeGenerator) VisitMemberAccessExpression(mae *ast.MemberAccessExpression) error {
	// Minimal stub: generate code for the left, ignore the member for now.
	if err := mae.Left.Accept(cg); err != nil {
		return err
	}
	// We are not producing anything real, so just clear or reuse lastValue
	cg.lastValue = nil
	return nil
}
