package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types" // Added for types.Void check
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitFunctionDefinition(fn *ast.FunctionDefinition) error {
	fnName := "anon"
	if fn.Name != nil && fn.Name.Value != "" {
		fnName = fn.Name.Value
	} else {
		// This path should ideally not be hit for named functions
		// If it's an anonymous function AST node, it should be handled by VisitLambdaExpression
		fmt.Printf("[WARN] VisitFunctionDefinition called for potentially anonymous function AST node.\n")
		return fmt.Errorf("VisitFunctionDefinition encountered anonymous function definition")
	}

	// 1. Find the pre-declared function.
	irFunc, ok := cg.Functions[fnName]
	if !ok {
		return fmt.Errorf("function '%s' was not pre-declared before visiting definition", fnName)
	}

	if len(irFunc.Blocks) > 0 && irFunc.Blocks[0].Term != nil {
		fmt.Printf("[DEBUG] Function '%s' already has a body, skipping definition.\n", fnName)
		return nil
	}

	fmt.Printf("[DEBUG] Generating body for function '%s' (Sig: %s)\n", fnName, irFunc.Sig.String())

	// 2. Create/Get entry block and set current context.
	var entry *ir.Block
	if len(irFunc.Blocks) == 0 {
		entry = irFunc.NewBlock("entry")
	} else {
		entry = irFunc.Blocks[0]
		if entry.Term != nil {
			fmt.Printf("[WARN] Function '%s' entry block already has terminator, potential redefinition?\n", fnName)
			return nil
		}
		fmt.Printf("[DEBUG] Reusing existing empty entry block for function '%s'\n", fnName)
	}

	// Store current context and set up for function body
	oldBlock := cg.Block
	oldFunc := cg.currentFunc
	oldVars := cg.Variables

	// Create a new variable map that inherits from the outer scope (for closures later)
	cg.Variables = make(map[string]value.Value)

	cg.Block = entry
	cg.currentFunc = irFunc

	// 3. Allocate space for parameters and store initial values
	if len(irFunc.Params) != len(fn.Parameters) {
		cg.Block = oldBlock
		cg.currentFunc = oldFunc
		cg.Variables = oldVars
		return fmt.Errorf("parameter count mismatch for function '%s': AST has %d, IR has %d",
			fnName, len(fn.Parameters), len(irFunc.Params))
	}
	for _, param := range irFunc.Params {
		paramIRName := param.Name()
		alloca := cg.Block.NewAlloca(param.Typ)
		alloca.SetName(paramIRName + ".addr")
		cg.Block.NewStore(param, alloca)
		// Make the parameter accessible by its name
		cg.setVar(paramIRName, alloca)
		fmt.Printf("[DEBUG] Function '%s': Allocated and stored parameter '%s' (type %s)\n", fnName, paramIRName, param.Typ)
	}

	// 4. Visit the function body.
	var bodyErr error
	if fn.Body != nil {
		bodyErr = fn.Body.Accept(cg)
	}

	// 5. Add a default return if the body didn't end with one.
	if cg.Block != nil && cg.Block.Term == nil {
		if bodyErr == nil {
			retType := cg.currentFunc.Sig.RetType

			if cg.lastValue != nil && cg.lastValue.Type().Equal(retType) {
				cg.Block.NewRet(cg.lastValue)
				fmt.Printf("[DEBUG] Function '%s': Added implicit return from last expression value\n", fnName)
			} else if !retType.Equal(types.Void) {
				zero := constant.NewZeroInitializer(retType)
				cg.Block.NewRet(zero)
				fmt.Printf("[DEBUG] Function '%s': Added default zero return (type %s)\n", fnName, retType)
			} else {
				cg.Block.NewRet(nil)
				fmt.Printf("[DEBUG] Function '%s': Added default void return\n", fnName)
			}
		} else {
			fmt.Printf("[WARN] Function '%s': Body errored, potentially leaving block unterminated.\n", fnName)
		}
	} else if cg.Block == nil {
		fmt.Printf("[DEBUG] Function '%s': Body resulted in nil block (all paths returned).\n", fnName)
	}

	// 6. Restore the previous context.
	cg.Block = oldBlock
	cg.currentFunc = oldFunc
	cg.Variables = oldVars

	if bodyErr != nil {
		return fmt.Errorf("error generating body for function '%s': %w", fnName, bodyErr)
	}

	fmt.Printf("[DEBUG] Finished generating body for function '%s'\n", fnName)
	return nil
}
