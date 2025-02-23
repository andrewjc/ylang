package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"math"
)

func (cg *CodeGenerator) VisitNumberLiteral(nl *ast.NumberLiteral) error {
	f := nl.Value
	intPart, frac := math.Modf(f)
	if frac == 0.0 {
		cg.lastValue = constant.NewInt(types.I32, int64(intPart))
	} else {
		cg.lastValue = constant.NewFloat(types.Float, f)
	}
	return nil
}
