package compiler

import (
	"compiler/ast"
	"compiler/compiler/generator"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// Compiler is the main struct for the compiler.
type Compiler struct {
	backend CompilerBackend
	errors  []string
	output  string
}

// CompilerBackend is the backend for the compiler.
type CompilerBackend int

// CompilerResult is the result of the compilation.
type CompilerResult struct {
	Errors []string
	Output string
}

const (
	// LLVM is the LLVM backend.
	LLVM CompilerBackend = iota
)

// LLVMNode is an interface that all AST nodes should implement if they
// are to be converted into LLVM IR.
type LLVMNode interface {
	LLVMValue(m *ir.Module) value.Value
}

// NewCompiler creates a new compiler.
func NewCompiler(backend CompilerBackend) *Compiler {
	return &Compiler{backend: backend}
}

func (c *Compiler) Compile(program *ast.Program) *CompilerResult {
	result := &CompilerResult{}

	if c.backend == LLVM {
		codeGen := generator.NewCodeGenerator()
		err := program.Accept(codeGen)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			return result
		}

		// Optionally, you can optimize the LLVM IR here

		// Generate the final IR string
		c.output = codeGen.Module.String()
		result.Output = c.output
	}

	return result
}
