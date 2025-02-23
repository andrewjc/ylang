package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

var lambdaCount int

func (cg *CodeGenerator) VisitLambdaExpression(le *ast.LambdaExpression) error {
	fnName := cg.newLambdaName()
	fnType := types.NewFunc(types.I32)
	irFunc := cg.Module.NewFunc(fnName, fnType.RetType)

	// Instead of irFunc.NewParam, do it manually
	for range le.Parameters {
		param := ir.NewParam("", types.I32)
		irFunc.Params = append(irFunc.Params, param)
	}

	oldBlock := cg.Block
	oldFunc := cg.currentFunc
	entry := irFunc.NewBlock("entry")

	cg.Block = entry
	cg.currentFunc = irFunc

	if le.Body != nil {
		if err := le.Body.Accept(cg); err != nil {
			return err
		}
	}

	if !cg.endsWithReturn(entry) {
		entry.NewRet(constant.NewInt(types.I32, 0))
	}

	cg.Block = oldBlock
	cg.currentFunc = oldFunc
	cg.lastValue = irFunc
	return nil
}

func (cg *CodeGenerator) newLambdaName() string {
	name := "lambda_"
	name = name + string(rune('A'+lambdaCount))
	lambdaCount++
	return name
}
