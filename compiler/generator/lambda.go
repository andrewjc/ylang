package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value" // Added for value.Value type
)

var lambdaCount int // Consider moving this to CodeGenerator struct if concurrency becomes a concern

func (cg *CodeGenerator) VisitLambdaExpression(le *ast.LambdaExpression) error {
	fnName := cg.newLambdaName()

	paramTypes := make([]types.Type, len(le.Parameters))
	paramNames := make([]string, len(le.Parameters))
	for i, paramAST := range le.Parameters {
		paramTypes[i] = types.I32 // Assuming i32 for now
		paramNames[i] = paramAST.Value
	}
	retType := types.I32 // Assuming i32 return for now

	funcParams := make([]*ir.Param, len(paramNames))
	for i, pName := range paramNames {
		funcParams[i] = ir.NewParam(pName, paramTypes[i])
	}

	irFunc := cg.Module.NewFunc(fnName, retType, funcParams...)
	irFunc.Linkage = enum.LinkageInternal

	fmt.Printf("[DEBUG] Lambda '%s': Created Func. Checking Sig(): %s\n", fnName, irFunc.Sig.String())
	fmt.Printf("[DEBUG] Lambda '%s': Func.Params field has %d entries.\n", fnName, len(irFunc.Params))

	oldBlock := cg.Block
	oldFunc := cg.currentFunc
	oldVars := cg.Variables
	lambdaScopeVars := make(map[string]value.Value)
	cg.Variables = lambdaScopeVars

	// lambda context
	entry := irFunc.NewBlock("entry")
	cg.Block = entry
	cg.currentFunc = irFunc

	if len(irFunc.Params) != len(paramNames) {
		fmt.Printf("[ERROR] Lambda '%s': Mismatch between AST params (%d) and irFunc.Params (%d) after creation!\n", fnName, len(paramNames), len(irFunc.Params))

		cg.Variables = oldVars
		cg.Block = oldBlock
		cg.currentFunc = oldFunc
		return fmt.Errorf("internal error: parameter count mismatch for lambda %s", fnName)
	}
	for _, param := range irFunc.Params {
		alloca := cg.Block.NewAlloca(param.Typ)
		alloca.SetName(param.Name() + ".addr")
		cg.Block.NewStore(param, alloca)
		cg.setVar(param.Name(), alloca)
		fmt.Printf("[DEBUG] Lambda '%s': Allocated and stored parameter '%s' (type %s) into NEW scope %p\n", fnName, param.Name(), param.Typ, cg.Variables)
	}

	fmt.Printf("[DEBUG] Lambda '%s': Visiting body with scope %p containing vars: %v\n", fnName, cg.Variables, cg.Variables)
	var bodyErr error
	if le.Body != nil {
		bodyErr = le.Body.Accept(cg)
	} else {
		fmt.Printf("[WARN] Lambda '%s' has nil body.\n", fnName)
	}

	if cg.Block != nil && cg.Block.Term == nil {
		if bodyErr == nil {
			if cg.lastValue != nil && cg.lastValue.Type().Equal(retType) {
				cg.Block.NewRet(cg.lastValue)
				fmt.Printf("[DEBUG] Lambda '%s': Added implicit return from last expression value\n", fnName)
			} else if !retType.Equal(types.Void) {
				zero := constant.NewZeroInitializer(retType)
				cg.Block.NewRet(zero)
				fmt.Printf("[DEBUG] Lambda '%s': Added default zero return (type %s)\n", fnName, retType)
			} else {
				cg.Block.NewRet(nil)
				fmt.Printf("[DEBUG] Lambda '%s': Added default void return\n", fnName)
			}
		} else {
			fmt.Printf("[WARN] Lambda '%s': Body errored, potentially leaving block unterminated.\n", fnName)
		}
	} else if cg.Block == nil {
		fmt.Printf("[DEBUG] Lambda '%s': Body resulted in nil block (all paths returned).\n", fnName)
	}

	cg.Block = oldBlock
	cg.currentFunc = oldFunc
	cg.Variables = oldVars

	if bodyErr != nil {
		return fmt.Errorf("error generating body for lambda '%s': %w", fnName, bodyErr)
	}

	cg.lastValue = irFunc
	fmt.Printf("[DEBUG] Finished generating lambda '%s' with signature %s\n", fnName, irFunc.Sig.String())
	return nil
}

func (cg *CodeGenerator) newLambdaName() string {
	name := fmt.Sprintf("lambda_%d", lambdaCount)
	lambdaCount++
	return name
}
