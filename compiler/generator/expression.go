package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// condAsBool converts any integer value into an i1 suitable for a conditional
// branch.  If the value is already i1 (produced by a comparison instruction),
// it is returned unchanged.  Otherwise an ICmp NE against the appropriate zero
// constant is emitted.
func condAsBool(block *ir.Block, v value.Value) value.Value {
	if iType, ok := v.Type().(*types.IntType); ok && iType.BitSize == 1 {
		return v
	}
	var zero constant.Constant
	if intType, ok := v.Type().(*types.IntType); ok {
		zero = constant.NewInt(intType, 0)
	} else {
		zero = constant.NewInt(types.I32, 0)
	}
	return block.NewICmp(enum.IPredNE, v, zero)
}

func (cg *CodeGenerator) VisitExpressionStatement(es *ast.ExpressionStatement) error {
	// Generate code for the expression, discard its value.
	if es.Expression == nil {
		return nil
	}
	return es.Expression.Accept(cg)
}

func (cg *CodeGenerator) VisitIfStatement(is *ast.IfStatement) error {
	// Evaluate condition
	if err := is.Condition.Accept(cg); err != nil {
		return err
	}
	condVal := cg.lastValue

	thenBlock := cg.newBlock("if_then")
	elseBlock := cg.newBlock("if_else")
	mergeBlock := cg.newBlock("if_merge")

	cg.Block.NewCondBr(condAsBool(cg.Block, condVal), thenBlock, elseBlock)

	// THEN branch — visit and then wire the current block to merge.
	cg.Block = thenBlock
	if is.Consequence != nil {
		if err := is.Consequence.Accept(cg); err != nil {
			return err
		}
	}
	if cg.Block != nil && cg.Block.Term == nil {
		cg.Block.NewBr(mergeBlock)
	}

	// ELSE branch
	cg.Block = elseBlock
	if is.Alternative != nil {
		if err := is.Alternative.Accept(cg); err != nil {
			return err
		}
	}
	if cg.Block != nil && cg.Block.Term == nil {
		cg.Block.NewBr(mergeBlock)
	}

	// Merge
	cg.Block = mergeBlock
	cg.lastValue = constant.NewInt(types.I32, 0)
	return nil
}

func (cg *CodeGenerator) VisitWhileStatement(ws *ast.WhileStatement) error {
	condBlock := cg.newBlock("while_cond")
	bodyBlock := cg.newBlock("while_body")
	exitBlock := cg.newBlock("while_exit")

	// Fall into the condition check.
	cg.Block.NewBr(condBlock)

	// Condition block: evaluate condition and branch.
	cg.Block = condBlock
	if err := ws.Condition.Accept(cg); err != nil {
		return err
	}
	cg.Block.NewCondBr(condAsBool(cg.Block, cg.lastValue), bodyBlock, exitBlock)

	// Body block — visit and loop back.
	cg.Block = bodyBlock
	if ws.Body != nil {
		if err := ws.Body.Accept(cg); err != nil {
			return err
		}
	}
	if cg.Block != nil && cg.Block.Term == nil {
		cg.Block.NewBr(condBlock)
	}

	// Exit block.
	cg.Block = exitBlock
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
