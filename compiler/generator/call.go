package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitCallExpression(ce *ast.CallExpression) error {
	// Evaluate the function expression (could be an identifier or lambda).
	if err := ce.Function.Accept(cg); err != nil {
		return err
	}
	fnVal := cg.lastValue

	var args []value.Value
	for _, argExpr := range ce.Arguments {
		if err := argExpr.Accept(cg); err != nil {
			return err
		}
		args = append(args, cg.lastValue)
	}

	// If fnVal is known function, we call it; else produce dummy i32 0.
	switch actual := fnVal.(type) {
	case *ir.Func:
		call := cg.Block.NewCall(actual, args...)
		cg.lastValue = call
	}

	return nil
}
