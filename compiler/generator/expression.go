package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitExpressionStatement(es *ast.ExpressionStatement) error {
	// Generate code for the expression, discard its value.
	if es.Expression == nil {
		return nil
	}
	return es.Expression.Accept(cg)
}

func (cg *CodeGenerator) VisitIndexExpression(ie *ast.IndexExpression) error {
	// Minimal stub: generate code for the left, ignore the index for now.
	if err := ie.Left.Accept(cg); err != nil {
		return err
	}
	// We are not producing anything real, so just clear or reuse lastValue
	cg.lastValue = nil
	return nil
}

func (cg *CodeGenerator) VisitIfStatement(is *ast.IfStatement) error {
	// Evaluate condition
	if err := is.Condition.Accept(cg); err != nil {
		return err
	}
	condVal := cg.lastValue

	// Compare condVal != 0 (treat 0 as false, non-zero as true)
	iCmp := cg.Block.NewICmp(enum.IPredNE, condVal, constant.NewInt(types.I32, 0))

	thenBlock := cg.currentFunc.NewBlock("if_then")
	elseBlock := cg.currentFunc.NewBlock("if_else")
	mergeBlock := cg.currentFunc.NewBlock("if_merge")

	cg.Block.NewCondBr(iCmp, thenBlock, elseBlock)

	// THEN branch
	cg.Block = thenBlock
	if is.Consequence != nil {
		if err := is.Consequence.Accept(cg); err != nil {
			return err
		}
	}
	if !cg.endsWithReturn(thenBlock) {
		thenBlock.NewBr(mergeBlock)
	}

	// ELSE branch
	cg.Block = elseBlock
	if is.Alternative != nil {
		if err := is.Alternative.Accept(cg); err != nil {
			return err
		}
	}
	if !cg.endsWithReturn(elseBlock) {
		elseBlock.NewBr(mergeBlock)
	}

	// Merge
	cg.Block = mergeBlock
	cg.lastValue = constant.NewInt(types.I32, 0)
	return nil
}

func (cg *CodeGenerator) VisitTraditionalTernaryExpression(te *ast.TraditionalTernaryExpression) error {
	//TODO implement me
	panic("implement me")
}

func (cg *CodeGenerator) VisitLambdaStyleTernaryExpression(aste *ast.LambdaStyleTernaryExpression) error {
	//TODO implement me
	panic("implement me")
}

func (cg *CodeGenerator) VisitInlineIfElseTernaryExpression(iite *ast.InlineIfElseTernaryExpression) error {
	//TODO implement me
	panic("implement me")
}

func (cg *CodeGenerator) VisitDotOperator(do *ast.DotOperator) error {
	// Evaluate the left side. For method calls, VisitCallExpression will re-evaluate
	// with inAssignmentLHS=true to get the pointer. For simple member access
	// (not implemented yet), this might need adjustment.
	if err := do.Left.Accept(cg); err != nil {
		return err
	}
	leftVal := cg.lastValue

	// Debug-print the 'leftVal'
	fmt.Printf("[DEBUG] DotOperator Left expression yielded: %v (Type: %s)\n", leftVal, leftVal.Type())

	methodName := do.Right.Value
	fmt.Printf("[DEBUG] DotOperator right side is '%s'\n", methodName)

	// For now, the value of a dot operation used outside a call context
	// is just the left value. This might need refinement for field access.
	cg.lastValue = leftVal
	return nil
}
