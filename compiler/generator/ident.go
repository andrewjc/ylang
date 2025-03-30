package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitIdentifier(id *ast.Identifier) error {
	identName := id.Value
	currentScope := cg.Variables // Capture the map instance being checked

	// 1. Check local variables in the current scope
	if allocaVal, ok := currentScope[identName]; ok { // Check the captured map instance
		fmt.Printf("[DEBUG] Identifier '%s': Found in current scope %p: %s\n", identName, currentScope, allocaVal.Ident())
		allocaInst, isAlloca := allocaVal.(*ir.InstAlloca)
		if !isAlloca {
			cg.lastValue = allocaVal
			fmt.Printf("[DEBUG] Identifier '%s' resolved to non-alloca variable: %s\n", identName, allocaVal.Ident())
			return nil
		}
		if cg.inAssignmentLHS {
			cg.lastValue = allocaInst
			fmt.Printf("[DEBUG] Identifier '%s' resolved to alloca (LHS): %s\n", identName, allocaInst.Ident())
			return nil
		}
		loaded := cg.Block.NewLoad(allocaInst.ElemType, allocaInst)
		cg.lastValue = loaded
		fmt.Printf("[DEBUG] Identifier '%s' resolved via load: %s from %s\n", identName, loaded.Ident(), allocaInst.Ident())
		return nil
	}

	// 2. Check global functions
	if fn, ok := cg.Functions[identName]; ok {
		cg.lastValue = fn
		fmt.Printf("[DEBUG] Identifier '%s' resolved to function: %s\n", identName, fn.Ident())
		return nil
	}

	// 3. Check global variables (if any added later)
	// todo!!!

	// 4. Not found - Implicit declaration (WARN)
	// Log the content of the scope map where it wasn't found
	fmt.Printf("[WARN] Identifier '%s' not found in current scope %p (Vars: %v), creating implicit external declaration i32()\n", identName, currentScope, currentScope)
	implicitSig := types.NewFunc(types.I32)
	implicitFunc := cg.Module.NewFunc(identName, implicitSig.RetType)
	cg.Functions[identName] = implicitFunc
	cg.lastValue = implicitFunc
	return nil
}
