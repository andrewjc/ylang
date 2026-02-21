package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitAssignmentStatement(as *ast.AssignmentStatement) error {
	//TODO implement me
	panic("implement me")
}

func (cg *CodeGenerator) VisitPrefixExpression(pe *ast.PrefixExpression) error {
	if err := pe.Right.Accept(cg); err != nil {
		return err
	}
	operand := cg.lastValue
	if operand == nil {
		return nil
	}
	switch pe.Operator {
	case "-":
		// Negate: 0 - operand
		zero := constant.NewInt(types.I32, 0)
		cg.lastValue = cg.Block.NewSub(zero, operand)
	case "!":
		// Logical not: icmp eq operand, 0
		cg.lastValue = cg.Block.NewICmp(enum.IPredEQ, operand, constant.NewInt(types.I32, 0))
	default:
		cg.lastValue = operand
	}
	return nil
}
