package generator

import (
	"compiler/ast"
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// makeSyscallInlineAsm builds an *ir.InlineAsm for a Linux x86-64 syscall with
// the provided operand types.  The first element of argTypes is always the
// syscall number (rax); subsequent elements map to rdi, rsi, rdx, r10, r8, r9.
// The return value is the i64 result left in rax after the instruction.
func makeSyscallInlineAsm(argTypes ...types.Type) *ir.InlineAsm {
	regNames := []string{"{rax}", "{rdi}", "{rsi}", "{rdx}", "{r10}", "{r8}", "{r9}"}
	parts := []string{"={rax}"}
	for i := range argTypes {
		if i < len(regNames) {
			parts = append(parts, regNames[i])
		}
	}
	parts = append(parts, "~{rcx}", "~{r11}", "~{memory}")

	funcType := types.NewFunc(types.I64, argTypes...)
	asm := ir.NewInlineAsm(types.NewPointer(funcType), "syscall", strings.Join(parts, ","))
	asm.SideEffect = true
	return asm
}

// coerceToI64 sign-extends or ptr-to-ints a value to i64 so it can be used
// as a syscall register operand.  Sign-extension is used for integer types so
// that negative values like AT_FDCWD (-100) are preserved correctly.
func coerceToI64(block *ir.Block, v value.Value) value.Value {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize == 64 {
			return v
		}
		return block.NewSExt(v, types.I64)
	case *types.PointerType:
		return block.NewPtrToInt(v, types.I64)
	default:
		return v
	}
}

func (cg *CodeGenerator) VisitSyscallExpression(se *ast.SyscallExpression) error {
	// 1) Evaluate syscall number.
	if err := se.Num.Accept(cg); err != nil {
		return err
	}
	numVal := coerceToI64(cg.Block, cg.lastValue)

	// 2) Evaluate up to 6 arguments, coercing each to i64.
	var argVals []value.Value
	for _, argNode := range se.Args {
		if err := argNode.Accept(cg); err != nil {
			return err
		}
		if cg.lastValue == nil {
			return fmt.Errorf("syscall argument produced no value")
		}
		argVals = append(argVals, coerceToI64(cg.Block, cg.lastValue))
	}

	// Pad missing args with 0.
	for len(argVals) < 6 {
		argVals = append(argVals, constant.NewInt(types.I64, 0))
	}

	// Build inline-asm operand type list: syscall# + 6 args, all i64.
	allArgTypes := make([]types.Type, 7)
	for i := range allArgTypes {
		allArgTypes[i] = types.I64
	}
	allArgs := append([]value.Value{numVal}, argVals...)

	asm := makeSyscallInlineAsm(allArgTypes...)
	call := cg.Block.NewCall(asm, allArgs...)
	cg.lastValue = call
	return nil
}
