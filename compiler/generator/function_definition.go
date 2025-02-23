package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitFunctionDefinition(fn *ast.FunctionDefinition) error {
	fnType := types.NewFunc(types.I32) // treat all as returning i32
	functionName := "anon"
	if fn.Name != nil && fn.Name.Value != "" {
		functionName = fn.Name.Value
	}

	irFunc := cg.Module.NewFunc(functionName, fnType.RetType)
	// Manually create params (replace old irFunc.NewParam calls)
	for range fn.Parameters {
		p := ir.NewParam("", types.I32)
		irFunc.Params = append(irFunc.Params, p)
	}

	cg.Functions[functionName] = irFunc

	entry := irFunc.NewBlock("entry")
	oldBlock := cg.Block
	oldFunc := cg.currentFunc

	cg.Block = entry
	cg.currentFunc = irFunc

	if fn.Body != nil {
		if err := fn.Body.Accept(cg); err != nil {
			return err
		}
	}

	// If no explicit return, default to 0
	if !cg.endsWithReturn(entry) {
		entry.NewRet(constant.NewInt(types.I32, 0))
	}

	// restore
	cg.Block = oldBlock
	cg.currentFunc = oldFunc
	return nil
}

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

func (cg *CodeGenerator) VisitBlockStatement(bs *ast.BlockStatement) error {
	for _, stmt := range bs.Statements {
		if stmt == nil {
			continue
		}
		if err := stmt.Accept(cg); err != nil {
			return err
		}
	}
	return nil
}
