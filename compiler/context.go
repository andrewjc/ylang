package compiler

import (
	ir "github.com/llir/llvm/ir"
	irv "github.com/llir/llvm/ir/value"
)

// compiler context holding symbol table, blocks, modules etc
type CompilerContext struct {
	SymbolTable map[string]*irv.Value
	Function    *ir.Func
	Module      *ir.Module
	Block       *ir.Block
}
