package generator

import (
	"compiler/ast"
	"compiler/module"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// CodeGenerator implements the Visitor interface to generate LLVM IR.
type CodeGenerator struct {
	ModuleManager *module.ModuleManager
	Module        *ir.Module
	Functions     map[string]*ir.Func
	Variables     map[string]value.Value
	Structs       map[string]types.Type
	Block         *ir.Block
	currentFunc   *ir.Func

	// lastValue holds the most recently produced LLVM value by a node visit.
	lastValue value.Value

	inAssignmentLHS bool

	methodCallReceiver value.Value
}

func NewCodeGenerator() *CodeGenerator {
	m := ir.NewModule()
	mm := module.NewModuleManager()
	builtInManager := NewBuiltInManager(m)

	return &CodeGenerator{
		ModuleManager: mm,
		Module:        m,
		Functions:     builtInManager.GetProvidedFunctionsMap(),
		Variables:     make(map[string]value.Value),
		Structs:       make(map[string]types.Type),
		Block:         nil,
		currentFunc:   nil,
		lastValue:     nil,
	}
}

func (cg *CodeGenerator) VisitVariableDeclaration(vd *ast.VariableDeclaration) error {
	fmt.Printf("[WARN] VisitVariableDeclaration called directly - currently only LetStatement handles local vars.\n")
	return nil
}

func (cg *CodeGenerator) VisitProgram(program *ast.Program) error {
	for _, is := range program.ImportStatements {
		if err := is.Accept(cg); err != nil {
			return fmt.Errorf("error visiting import %s: %w", is.Path, err)
		}
	}

	// Pre-declare all functions (including main) to handle forward references
	// and allow module integration to find them.
	if program.MainFunction != nil {
		if err := cg.declareFunction(program.MainFunction); err != nil {
			return fmt.Errorf("error declaring main function: %w", err)
		}
	}
	for _, fn := range program.Functions {
		if err := cg.declareFunction(fn); err != nil {
			return fmt.Errorf("error declaring function %s: %w", fn.Name.Value, err)
		}
	}

	// Visit each normal function definition to generate its body.
	for _, fn := range program.Functions {
		if err := fn.Accept(cg); err != nil {
			return fmt.Errorf("error visiting function %s: %w", fn.Name.Value, err)
		}
	}
	// Then visit the main function definition, if any.
	if program.MainFunction != nil {
		if err := program.MainFunction.Accept(cg); err != nil {
			return fmt.Errorf("error visiting main function: %w", err)
		}
	}
	return nil
}

func (cg *CodeGenerator) declareFunction(fn *ast.FunctionDefinition) error {
	fnName := "anon"
	if fn.Name != nil && fn.Name.Value != "" {
		fnName = fn.Name.Value
	} else {
		return fmt.Errorf("declareFunction received anonymous function AST node")
	}

	if existingFunc, exists := cg.Functions[fnName]; exists {
		fmt.Printf("[DEBUG] Function '%s' already declared (Sig: %s), skipping.\n", fnName, existingFunc.Sig.String())
		return nil
	}

	// Determine parameter types and names
	paramTypes := make([]types.Type, len(fn.Parameters))
	paramNames := make([]string, len(fn.Parameters))
	fmt.Printf("[DEBUG] declareFunction '%s': Processing %d AST parameters.\n", fnName, len(fn.Parameters))
	for i, paramAST := range fn.Parameters {
		paramTypes[i] = types.I32 // Assuming i32 for now fix later
		paramNames[i] = paramAST.Value
		fmt.Printf("[DEBUG] declareFunction '%s': Found param %d: Name='%s' (assuming Type=i32)\n", fnName, i, paramNames[i])
	}

	// Determine return type
	var retType types.Type = types.I32

	if fn.ReturnType != nil {
		mappedType, err := cg.mapType(fn.ReturnType.Value)
		if err == nil {
			retType = mappedType
			fmt.Printf("[DEBUG] declareFunction '%s': Using explicit return type '%s' -> %s\n", fnName, fn.ReturnType.Value, retType)
		} else {
			fmt.Printf("[WARN] declareFunction '%s': Could not map explicit return type '%s': %v. Defaulting to i32.\n", fnName, fn.ReturnType.Value, err)
			retType = types.I32
		}
	} else {
		canInferVoid := false
		if fn.Body != nil {
			if block, ok := fn.Body.(*ast.BlockStatement); ok {
				if len(block.Statements) > 0 {
					allVoidAsm := true
					for _, stmt := range block.Statements {
						expStmt, isExpStmt := stmt.(*ast.ExpressionStatement)
						if !isExpStmt {
							allVoidAsm = false
							break
						}
						asmExpr, isAsm := expStmt.Expression.(*ast.AssemblyExpression)
						if !isAsm {
							allVoidAsm = false
							break
						}
						asmCode := asmExpr.Code.Value
						isVoidBuiltin := false
						if builtinFunc, exists := cg.Functions[asmCode]; exists {
							if builtinFunc.Sig.RetType.Equal(types.Void) {
								isVoidBuiltin = true
							}
						}
						if !isVoidBuiltin {
							allVoidAsm = false
							break
						}
					}
					if allVoidAsm {
						canInferVoid = true
					}
				} else {
					canInferVoid = true
				}
			}
		} else {
			canInferVoid = true
		}

		if canInferVoid {
			retType = types.Void
			fmt.Printf("[DEBUG] declareFunction '%s': Inferred return type void.\n", fnName)
		} else {
			retType = types.I32
			fmt.Printf("[DEBUG] declareFunction '%s': Could not infer void return type, defaulting to i32.\n", fnName)
		}
	}

	funcParams := make([]*ir.Param, len(paramNames))
	for i, pName := range paramNames {
		funcParams[i] = ir.NewParam(pName, paramTypes[i])
	}

	irFunc := cg.Module.NewFunc(fnName, retType, funcParams...)

	fmt.Printf("[DEBUG] declareFunction '%s': Created Func. Checking Sig(): %s\n", fnName, irFunc.Sig.String())

	fmt.Printf("[DEBUG] declareFunction '%s': Func.Params field has %d entries.\n", fnName, len(irFunc.Params))

	cg.Functions[fnName] = irFunc
	fmt.Printf("[DEBUG] Stored function '%s' in Functions map.\n", fnName)
	return nil
}

func (cg *CodeGenerator) resolveStructType(typeName string) (*types.StructType, error) {
	t, ok := cg.Structs[typeName]
	if !ok {
		// Type definition should have been processed before trying to resolve it.
		return nil, fmt.Errorf("struct type '%s' not defined or resolved", typeName)
	}
	// Check if the retrieved type is actually a struct type
	st, ok := t.(*types.StructType)
	if !ok {
		return nil, fmt.Errorf("type '%s' found but is not a struct type (%T)", typeName, t)
	}
	return st, nil
}
