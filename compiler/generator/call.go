package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) evaluateArguments(argNodes []ast.ExpressionNode) ([]value.Value, error) {
	args := make([]value.Value, 0, len(argNodes))
	for _, argExpr := range argNodes {
		if err := argExpr.Accept(cg); err != nil {
			return nil, err
		}
		argVal := cg.lastValue
		if argVal == nil {
			return nil, fmt.Errorf("argument expression produced no value")
		}

		if allocaInst, isAlloca := argVal.(*ir.InstAlloca); isAlloca {
			if ptrType, isPtr := allocaInst.ElemType.(*types.PointerType); isPtr {
				if _, isFunc := ptrType.ElemType.(*types.FuncType); isFunc {
					loadedFnPtr := cg.Block.NewLoad(allocaInst.ElemType, allocaInst)
					fmt.Printf("[DEBUG] Loaded function pointer for argument: %s from %s\n", loadedFnPtr.Ident(), allocaInst.Ident())
					argVal = loadedFnPtr
				}
			}
		}
		args = append(args, argVal)
	}
	return args, nil
}

func (cg *CodeGenerator) VisitCallExpression(ce *ast.CallExpression) error {
	args, err := cg.evaluateArguments(ce.Arguments)
	if err != nil {
		funcNameStr := "unknown_function"
		if ce.Function != nil {
			funcNameStr = ce.Function.String()
		}
		return fmt.Errorf("error evaluating arguments for call to '%s': %w", funcNameStr, err)
	}

	if memberAccessExpr, isMemberAccess := ce.Function.(*ast.MemberAccessExpression); isMemberAccess {
		fmt.Printf("[DEBUG] Detected method call: %s\n", memberAccessExpr.String())
		err := memberAccessExpr.Left.Accept(cg)
		if err != nil {
			return fmt.Errorf("error evaluating receiver for method call '%s': %w", memberAccessExpr.Member.Value, err)
		}
		objReceiver := cg.lastValue

		if objReceiver == nil {
			return fmt.Errorf("method call receiver '%s' evaluated to nil", memberAccessExpr.Left.String())
		}

		methodName := memberAccessExpr.Member.Value

		return cg.handleMethodCall(objReceiver, methodName, args)

	} else {
		fmt.Printf("[DEBUG] Detected regular function call: %s\n", ce.Function.String())

		if err := ce.Function.Accept(cg); err != nil {
			return fmt.Errorf("error evaluating function expression '%s': %w", ce.Function.String(), err)
		}
		fnVal := cg.lastValue // This should be an *ir.Func or a function pointer

		if fnVal == nil {
			return fmt.Errorf("function expression '%s' evaluated to nil", ce.Function.String())
		}

		var callableFn value.Value
		var fnSig *types.FuncType

		switch fn := fnVal.(type) {
		case *ir.Func:
			callableFn = fn
			fnSig = fn.Sig
			fmt.Printf("[DEBUG] Calling direct function: %s\n", fn.Name())
		case *ir.InstAlloca:
			elemType := fn.ElemType
			ptrType, isPtr := elemType.(*types.PointerType)
			if !isPtr {
				return fmt.Errorf("cannot call variable '%s', it does not store a pointer (type: %s)", fn.Name(), elemType)
			}
			sig, isFunc := ptrType.ElemType.(*types.FuncType)
			if !isFunc {
				return fmt.Errorf("cannot call variable '%s', it does not store a function pointer (stores pointer to: %s)", fn.Name(), ptrType.ElemType)
			}
			loadedFnPtr := cg.Block.NewLoad(elemType, fn) // Load the function pointer (e.g., i32 (...)**)
			callableFn = loadedFnPtr
			fnSig = sig
			fmt.Printf("[DEBUG] Calling loaded function pointer: %s from %s\n", loadedFnPtr.Ident(), fn.Ident())
		case value.Value:
			ptrType, isPtr := fn.Type().(*types.PointerType)
			if !isPtr {
				return fmt.Errorf("cannot call value of type %T, not a function or pointer", fnVal)
			}
			sig, isFunc := ptrType.ElemType.(*types.FuncType)
			if !isFunc {
				return fmt.Errorf("cannot call pointer value of type %s, does not point to a function", fn.Type())
			}
			callableFn = fn
			fnSig = sig
			fmt.Printf("[DEBUG] Calling function pointer value: %s\n", callableFn.Ident())
		default:
			return fmt.Errorf("cannot call value of type %T", fnVal)
		}

		if len(fnSig.Params) != len(args) {
			return fmt.Errorf("argument count mismatch for call to '%s': expected %d, got %d", callableFn.String(), len(fnSig.Params), len(args))
		}

		// TODO: Type check arguments against fnSig.Params

		call := cg.Block.NewCall(callableFn, args...)
		if !fnSig.RetType.Equal(types.Void) {
			cg.lastValue = call
		} else {
			cg.lastValue = nil
		}

		return nil
	}
}

func (cg *CodeGenerator) handleMethodCall(objReceiver value.Value, methodName string, args []value.Value) error {
	// 1. Get receiver type info
	objPtrType, ok := objReceiver.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("method call receiver is not a pointer, but %T", objReceiver.Type())
	}
	objStructType, ok := objPtrType.ElemType.(*types.StructType)
	if !ok {
		return fmt.Errorf("method call receiver does not point to a struct, but %T", objPtrType.ElemType)
	}
	typeName := objStructType.Name() // Get "Array"

	fmt.Printf("[DEBUG] Method Call: receiverType=%s, method=%s\n", typeName, methodName)

	// 2. Find the LLVM function for the method
	// How to link AST MethodDeclaration to LLVM ir.Func?
	// Assume a naming convention or store mapping during class processing.
	// Convention: ClassName_MethodName? e.g., "Array_map"
	mangledName := typeName + "_" + methodName // Simple mangling
	llvmMethodFunc, funcExists := cg.Functions[mangledName]

	if !funcExists {
		// Fallback: Maybe it's a built-in method implemented directly in Go?
		// This section needs refinement based on how stdlib/builtins are truly handled.
		// For now, let's assume all methods MUST be defined in YLang.
		return fmt.Errorf("method '%s' not found for type '%s' (tried mangled name '%s')", methodName, typeName, mangledName)
	}

	// 3. Prepare arguments (prepend self)
	allArgs := append([]value.Value{objReceiver}, args...) // objReceiver is 'self'

	// 4. Check argument count (LLVM func params should be N+1)
	if len(llvmMethodFunc.Sig.Params) != len(allArgs) {
		return fmt.Errorf("argument count mismatch for method call '%s.%s': expected %d (including self), got %d", typeName, methodName, len(llvmMethodFunc.Sig.Params), len(allArgs))
	}

	// 5. Type check arguments (TODO)
	// for i, arg := range allArgs {
	//    if !arg.Type().Equal(llvmMethodFunc.Sig.Params[i]) { ... error ... }
	// }

	// 6. Generate the call instruction
	callInst := cg.Block.NewCall(llvmMethodFunc, allArgs...)
	callInst.SetName(methodName + "_res")

	// 7. Set lastValue if method returns something
	if !llvmMethodFunc.Sig.RetType.Equal(types.Void) {
		cg.lastValue = callInst
	} else {
		cg.lastValue = nil
	}

	return nil
}
