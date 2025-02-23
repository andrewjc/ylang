package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitArrayLiteral(al *ast.ArrayLiteral) error {
	count := len(al.Elements)
	arrayTy := types.NewArray(uint64(count), types.I32)

	var constVals []constant.Constant
	for _, elem := range al.Elements {
		if err := elem.Accept(cg); err != nil {
			return err
		}
		if c, ok := cg.lastValue.(constant.Constant); ok {
			constVals = append(constVals, c)
		} else {
			constVals = append(constVals, constant.NewInt(types.I32, 0))
		}
	}

	// Pass the array type itself, not just arrayTy.ElemType.
	arrConst := constant.NewArray(arrayTy, constVals...)
	cg.lastValue = arrConst
	return nil
}
