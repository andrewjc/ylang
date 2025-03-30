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
	if dotOp, isDotOp := ce.Function.(*ast.DotOperator); isDotOp {
		isLHS := cg.inAssignmentLHS
		cg.inAssignmentLHS = true
		err := dotOp.Left.Accept(cg)
		cg.inAssignmentLHS = isLHS
		if err != nil {
			return fmt.Errorf("error evaluating method call receiver '%s': %w", dotOp.Left.String(), err)
		}
		objReceiver := cg.lastValue

		args, err := cg.evaluateArguments(ce.Arguments)
		if err != nil {
			return fmt.Errorf("error evaluating arguments for method call '%s': %w", dotOp.Right.Value, err)
		}

		methodName := dotOp.Right.Value

		return cg.handleMethodCall(objReceiver, methodName, args)
	}

	if err := ce.Function.Accept(cg); err != nil {
		return err
	}
	fnVal := cg.lastValue

	args, err := cg.evaluateArguments(ce.Arguments)
	if err != nil {
		funcNameStr := ce.Function.String()
		return fmt.Errorf("error evaluating arguments for function call '%s': %w", funcNameStr, err)
	}

	var callableFn value.Value
	var fnSig *types.FuncType

	switch fn := fnVal.(type) {
	case *ir.Func:
		callableFn = fn
		fnSig = fn.Sig
	case *ir.InstAlloca:
		elemType := fn.ElemType
		ptrType, isPtr := elemType.(*types.PointerType)
		if !isPtr {
			return fmt.Errorf("cannot call variable '%s', it does not store a pointer", fn.Name())
		}
		sig, isFunc := ptrType.ElemType.(*types.FuncType)
		if !isFunc {
			return fmt.Errorf("cannot call variable '%s', it does not store a function pointer", fn.Name())
		}
		loadedFnPtr := cg.Block.NewLoad(elemType, fn)
		callableFn = loadedFnPtr
		fnSig = sig
		fmt.Printf("[DEBUG] Loaded function pointer for call: %s from %s\n", loadedFnPtr.Ident(), fn.Ident())
	default:
		return fmt.Errorf("cannot call value of type %T", fnVal)
	}

	if len(fnSig.Params) != len(args) {
		return fmt.Errorf("argument count mismatch for call to '%s': expected %d, got %d", callableFn.String(), len(fnSig.Params), len(args))
	}

	call := cg.Block.NewCall(callableFn, args...)
	if !fnSig.RetType.Equal(types.Void) {
		cg.lastValue = call
	} else {
		cg.lastValue = nil
	}

	return nil
}

func (cg *CodeGenerator) handleMethodCall(objReceiver value.Value, methodName string, args []value.Value) error {
	var objType types.Type
	if ptrType, isPtr := objReceiver.Type().(*types.PointerType); isPtr {
		objType = ptrType.ElemType
	} else {
		return fmt.Errorf("method call receiver is not a pointer, but %T (%s)", objReceiver, objReceiver.Type())
	}

	fmt.Printf("[DEBUG] Method Call: receiverType=%s, method=%s\n", objType, methodName)

	if arrayType, isArray := objType.(*types.ArrayType); isArray {
		switch methodName {
		case "map":
			return cg.generateArrayMap(objReceiver, arrayType, args)
		case "forEach":
			return cg.generateArrayForEach(objReceiver, arrayType, args)
		default:
			return fmt.Errorf("unknown method '%s' for array type %s", methodName, objType)
		}
	}

	return fmt.Errorf("method call '%s' not supported on type %s", methodName, objType)
}
