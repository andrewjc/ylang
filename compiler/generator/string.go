package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

// We'll store string constants globally and return an i8* pointer to them.
var stringCounter int

func (cg *CodeGenerator) VisitStringLiteral(sl *ast.StringLiteral) error {
	strName := fmt.Sprintf("str_%d", stringCounter)
	stringCounter++

	// Null-terminate
	raw := []byte(sl.Value + "\x00")
	arrType := types.NewArray(uint64(len(raw)), types.I8)

	// Create global definition
	g := cg.Module.NewGlobalDef(strName, constant.NewArray(arrType, bytesToConstants(raw)...))
	// Use Linkage = enum.LinkagePrivate or another suitable value
	g.Linkage = enum.LinkagePrivate
	g.Immutable = true

	zero := constant.NewInt(types.I32, 0)
	// GEP to get i8* pointer
	gep := cg.Block.NewGetElementPtr(arrType, g, zero, zero)
	cg.lastValue = gep
	return nil
}

func bytesToConstants(data []byte) []constant.Constant {
	elems := make([]constant.Constant, len(data))
	for i, b := range data {
		elems[i] = constant.NewInt(types.I8, int64(b))
	}
	return elems
}
