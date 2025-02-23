package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitIdentifier(id *ast.Identifier) error {
	val := cg.getVar(id.Value)
	if val != nil {
		// For llir >= v0.4.0, loads require a type.
		// We assume everything is i32 for simplicity.
		loaded := cg.Block.NewLoad(types.I32, val)
		cg.lastValue = loaded
	} else {
		// fallback to i32 0
		cg.lastValue = nil
	}
	return nil
}
