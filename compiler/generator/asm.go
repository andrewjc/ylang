package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitAssemblyExpression(ae *ast.AssemblyExpression) error {
	asmCode := ae.Code.Value

	// Evaluate arguments
	var args []value.Value
	for _, argExpr := range ae.Args {
		if err := argExpr.Accept(cg); err != nil {
			return fmt.Errorf("error evaluating argument for asm '%s': %w", asmCode, err)
		}
		if cg.lastValue == nil {
			return fmt.Errorf("argument expression for asm '%s' produced no value", asmCode)
		}
		args = append(args, cg.lastValue)
	}

	switch asmCode {
	case "builtin_print_int":
		if fn, ok := cg.Functions["builtin_print_int"]; ok {
			if len(args) != 1 {
				return fmt.Errorf("asm 'builtin_print_int' expects 1 argument, got %d", len(args))
			}
			cg.Block.NewCall(fn, args...)
			cg.lastValue = nil
			return nil
		} else {
			return fmt.Errorf("builtin function 'builtin_print_int' not declared")
		}
	case "builtin_print_newline":
		if fn, ok := cg.Functions["builtin_print_newline"]; ok {
			if len(args) != 0 {
				return fmt.Errorf("asm 'builtin_print_newline' expects 0 arguments, got %d", len(args))
			}
			cg.Block.NewCall(fn, args...)
			cg.lastValue = nil
			return nil
		} else {
			return fmt.Errorf("builtin function 'builtin_print_newline' not declared")
		}

	case "builtin_map":
		fmt.Println("[WARN] asm 'builtin_map' not fully implemented")
		cg.lastValue = constant.NewNull(types.NewPointer(types.I32))
		return nil
	case "builtin_forEach":
		fmt.Println("[WARN] asm 'builtin_forEach' not fully implemented")
		cg.lastValue = nil
		return nil

	default:
		return fmt.Errorf("unsupported or unknown asm code: '%s'", asmCode)
	}
}
