package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// currentFunc holds the function we are generating code for at the moment.
func (cg *CodeGenerator) endsWithReturn(block *ir.Block) bool {
	if block.Term == nil {
		return false
	}
	_, isRet := block.Term.(*ir.TermRet)
	return isRet
}

// getVar returns the LLVM value for a named variable if it exists.
func (cg *CodeGenerator) getVar(name string) value.Value {
	if v, ok := cg.Variables[name]; ok {
		return v
	}
	return nil
}

// setVar sets the LLVM value for a named variable.
func (cg *CodeGenerator) setVar(name string, val value.Value) {
	cg.Variables[name] = val
}

// helper for ternary-like expressions
func (cg *CodeGenerator) ternaryLike(cond, trueExpr, falseExpr ast.ExpressionNode) error {
	// Evaluate condition
	if err := cond.Accept(cg); err != nil {
		return err
	}
	condVal := cg.lastValue
	iCmp := cg.Block.NewICmp(enum.IPredNE, condVal, constant.NewInt(types.I32, 0))

	thenBlock := cg.currentFunc.NewBlock("ternary_then")
	elseBlock := cg.currentFunc.NewBlock("ternary_else")
	mergeBlock := cg.currentFunc.NewBlock("ternary_merge")

	cg.Block.NewCondBr(iCmp, thenBlock, elseBlock)

	// THEN branch
	cg.Block = thenBlock
	if err := trueExpr.Accept(cg); err != nil {
		return err
	}
	trueVal := cg.lastValue
	if !cg.endsWithReturn(thenBlock) {
		thenBlock.NewBr(mergeBlock)
	}

	// ELSE branch
	cg.Block = elseBlock
	if err := falseExpr.Accept(cg); err != nil {
		return err
	}
	falseVal := cg.lastValue
	if !cg.endsWithReturn(elseBlock) {
		elseBlock.NewBr(mergeBlock)
	}

	// MERGE
	cg.Block = mergeBlock
	phi := mergeBlock.NewPhi(
		&ir.Incoming{X: trueVal, Pred: thenBlock},
		&ir.Incoming{X: falseVal, Pred: elseBlock},
	)
	cg.lastValue = phi
	return nil
}

// minimal helper to create a stub function
func (cg *CodeGenerator) ensureFuncExists(name string) *ir.Func {
	if existing, ok := cg.Functions[name]; ok {
		return existing
	}
	// create a stub function: i32 myFunc(i32)
	sig := types.NewFunc(types.I32, types.I32)
	fn := cg.Module.NewFunc(name, sig.RetType)
	fn.Params = append(fn.Params, ir.NewParam("p0", types.I32))
	blk := fn.NewBlock("entry")
	blk.NewRet(constant.NewInt(types.I32, 0))
	cg.Functions[name] = fn
	return fn
}
