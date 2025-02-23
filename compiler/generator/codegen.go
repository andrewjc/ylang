package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir"
	ivt "github.com/llir/llvm/ir/types"
	irv "github.com/llir/llvm/ir/value"
	"log"
)

// CodeGenerator implements the Visitor interface to generate LLVM IR.
type CodeGenerator struct {
	Module    *ir.Module
	Functions map[string]*ir.Func
	Variables map[string]irv.Value
	Structs   map[string]*ivt.StructType
}

func NewCodeGenerator() *CodeGenerator {
	m := ir.NewModule()
	return &CodeGenerator{
		Module:    m,
		Functions: make(map[string]*ir.Func),
		Variables: make(map[string]irv.Value),
		Structs:   make(map[string]*ivt.StructType),
	}
}

func (cg *CodeGenerator) VisitProgram(program *ast.Program) error {
	// Process Class Declarations
	/*for _, class := range p.ClassDeclarations {
		if err := class.Visit(g); err != nil {
			return err
		}
	}

	// Process Data Structures
	for _, ds := range p.DataStructures {
		if err := ds.Visit(g); err != nil {
			return err
		}
	}*/

	// Handle Function Declarations
	for _, fn := range program.Functions {
		err := fn.Visit(cg)
		if err != nil {
			return err
		}
	}

	// Handle Main Function
	if program.MainFunction != nil {
		err := program.MainFunction.Visit(cg)
		if err != nil {
			return err
		}
	} else {
		log.Println("No main function found.")
	}

	return nil
}
