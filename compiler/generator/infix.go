package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitInfixExpression(ie *ast.InfixExpression) error {
	// Generate left
	if err := ie.Left.Accept(cg); err != nil {
		return err
	}
	leftVal := cg.lastValue
	if leftVal == nil {
		return nil
	}

	// Generate right
	if err := ie.Right.Accept(cg); err != nil {
		return err
	}
	rightVal := cg.lastValue
	if rightVal == nil {
		return nil
	}

	// Assume both are i32
	var result value.Value
	switch ie.Operator {
	case "+":
		result = cg.Block.NewAdd(leftVal, rightVal)
	case "-":
		result = cg.Block.NewSub(leftVal, rightVal)
	case "*":
		result = cg.Block.NewMul(leftVal, rightVal)
	case "/":
		result = cg.Block.NewSDiv(leftVal, rightVal)
	default:
		// fallback
		result = leftVal
	}
	cg.lastValue = result
	return nil
}
