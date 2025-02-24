package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitSyscallExpression(se *ast.SyscallExpression) error {
	// 1) Evaluate se.Num => syscall number
	if err := se.Num.Accept(cg); err != nil {
		return err
	}
	numVal := cg.lastValue // e.g. i64

	// 2) Evaluate each of the up to 6 arguments
	var argVals []value.Value
	for _, argNode := range se.Args {
		if err := argNode.Accept(cg); err != nil {
			return err
		}
		argVals = append(argVals, cg.lastValue)
	}

	// We'll store them in the registers:
	//   rax = syscall number
	//   rdi = arg0
	//   rsi = arg1
	//   rdx = arg2
	//   r10 = arg3
	//   r8  = arg4
	//   r9  = arg5
	//
	// Then do 'syscall' instruction. We'll use inline assembly.

	// We'll want to handle variable argument counts. Let's define defaults:
	var reg0 value.Value = constant.NewInt(types.I64, 0)
	var reg1 value.Value = constant.NewInt(types.I64, 0)
	var reg2 value.Value = constant.NewInt(types.I64, 0)
	var reg3 value.Value = constant.NewInt(types.I64, 0)
	var reg4 value.Value = constant.NewInt(types.I64, 0)
	var reg5 value.Value = constant.NewInt(types.I64, 0)

	if len(argVals) > 0 {
		reg0 = argVals[0]
	}
	if len(argVals) > 1 {
		reg1 = argVals[1]
	}
	if len(argVals) > 2 {
		reg2 = argVals[2]
	}
	if len(argVals) > 3 {
		reg3 = argVals[3]
	}
	if len(argVals) > 4 {
		reg4 = argVals[4]
	}
	if len(argVals) > 5 {
		reg5 = argVals[5]
	}

	// We'll build inline assembly with placeholders
	// In LLVM: "mov rax,$0; mov rdi,$1; mov rsi,$2; ...; syscall"
	// Then we feed the arguments in the correct order.

	asmTemplate := `
    mov rax, $0
    mov rdi, $1
    mov rsi, $2
    mov rdx, $3
    mov r10, $4
    mov r8, $5
    mov r9, $6
    syscall
    `
	// We'll produce one i64 return (the result in RAX)
	inlineAsmTy := types.I64
	inlineAsm := ir.NewInlineAsm(
		inlineAsmTy,
		asmTemplate,
		"", // constraints... we can do "r,r,r,r,r,r,r"
	)
	inlineAsm.SideEffect = true // verify this?

	call := cg.Block.NewCall(
		inlineAsm,
		// The arguments in the order of placeholders $0...$6
		numVal, reg0, reg1, reg2, reg3, reg4, reg5,
	)
	cg.lastValue = call

	return nil
}
