package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"log"
)

// CodeGenerator implements the Visitor interface to generate LLVM IR.
type CodeGenerator struct {
	Module      *ir.Module
	Functions   map[string]*ir.Func
	Variables   map[string]value.Value
	Structs     map[string]*types.Type
	Block       *ir.Block
	currentFunc *ir.Func

	// lastValue holds the most recently produced LLVM value by a node visit.
	lastValue value.Value

	inAssignmentLHS bool // are we visiting the left-hand side of an assignment?
}

func NewCodeGenerator() *CodeGenerator {
	m := ir.NewModule()
	return &CodeGenerator{
		Module:    m,
		Functions: make(map[string]*ir.Func),
		Variables: make(map[string]value.Value),
		Structs:   make(map[string]*types.Type),

		Block:       nil,
		currentFunc: nil,
		lastValue:   nil,
	}
}

func (cg *CodeGenerator) VisitVariableDeclaration(vd *ast.VariableDeclaration) error {
	panic("implement me")
}

func (cg *CodeGenerator) VisitProgram(program *ast.Program) error {
	// Visit each normal function.
	for _, fn := range program.Functions {
		if err := fn.Accept(cg); err != nil {
			return err
		}
	}
	// Then visit the main function, if any.
	if program.MainFunction != nil {
		if err := program.MainFunction.Accept(cg); err != nil {
			return err
		}
	} else {
		log.Println("No main function found.")
	}
	return nil
}
