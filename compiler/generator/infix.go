package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// coerceToSameIntType zero-extends the narrower of two integer values to
// match the wider one. Pointer types are left unchanged (returned as-is).
func coerceToSameIntType(block *ir.Block, a, b value.Value) (value.Value, value.Value) {
	aInt, aOk := a.Type().(*types.IntType)
	bInt, bOk := b.Type().(*types.IntType)
	if !aOk || !bOk {
		return a, b
	}
	if aInt.BitSize == bInt.BitSize {
		return a, b
	}
	if aInt.BitSize > bInt.BitSize {
		return a, block.NewZExt(b, a.Type())
	}
	return block.NewZExt(a, b.Type()), b
}

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

	// Promote operands to the same width before applying the operator.
	leftVal, rightVal = coerceToSameIntType(cg.Block, leftVal, rightVal)

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
	case "==":
		result = cg.Block.NewICmp(enum.IPredEQ, leftVal, rightVal)
	case "<":
		result = cg.Block.NewICmp(enum.IPredSLT, leftVal, rightVal)
	case ">":
		result = cg.Block.NewICmp(enum.IPredSGT, leftVal, rightVal)
	case "<=":
		result = cg.Block.NewICmp(enum.IPredSLE, leftVal, rightVal)
	default:
		// fallback
		result = leftVal
	}
	cg.lastValue = result
	return nil
}
