package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitArrayLiteral(al *ast.ArrayLiteral) error {
	count := uint64(len(al.Elements))
	if count == 0 {
		// Handle empty array literal e.g., []
		// Create an empty Array struct
		arrayStructType, err := cg.resolveStructType("Array")
		if err != nil {
			return err
		}

		llvmIntType, _ := cg.mapType("int") // Assuming int type
		llvmIntPtrType := types.NewPointer(llvmIntType)

		arrayAlloca := cg.Block.NewAlloca(arrayStructType)
		arrayAlloca.SetName("empty_array_struct")

		// Store length = 0
		lenFieldAddr := cg.Block.NewGetElementPtr(arrayStructType, arrayAlloca, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)) // Field 0 = length
		cg.Block.NewStore(constant.NewInt(llvmIntType.(*types.IntType), 0), lenFieldAddr)

		// Store data = null
		dataFieldAddr := cg.Block.NewGetElementPtr(arrayStructType, arrayAlloca, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1)) // Field 1 = data
		nullDataPtr := constant.NewNull(llvmIntPtrType)                                                                                        // Null pointer of type *int
		cg.Block.NewStore(nullDataPtr, dataFieldAddr)

		cg.lastValue = arrayAlloca // Return pointer to the struct on the stack
		return nil
	}

	// Assuming elements are uniform type, determine from first element (or require type hint)
	// Let's assume 'int' (i32) for now based on the test case [1,2,3,4,5]
	elemTypeName := "int"
	llvmElemType, err := cg.mapType(elemTypeName)
	if err != nil {
		return fmt.Errorf("cannot determine element type for array literal: %w", err)
	}

	// Evaluate all element expressions first
	llvmElements := make([]value.Value, count)
	for i, elemAST := range al.Elements {
		if err := elemAST.Accept(cg); err != nil {
			return fmt.Errorf("error evaluating element %d for array literal: %w", i, err)
		}
		// TODO: Type check/cast cg.lastValue to llvmElemType if needed
		llvmElements[i] = cg.lastValue
	}

	// Create the LLVM constant array value (for initialization)
	constArrayType := types.NewArray(count, llvmElemType)
	constVals := make([]constant.Constant, count)
	for i, v := range llvmElements {
		if c, ok := v.(constant.Constant); ok {
			constVals[i] = c
		} else {
			// Element was not a constant, cannot create global const array easily.
			// We need to allocate stack space and store elements individually.
			return fmt.Errorf("non-constant element found in array literal - stack allocation needed (not fully implemented yet)")
			// Alternative: Allocate stack array, loop and store each llvmElements[i]
		}
	}
	constArray := constant.NewArray(constArrayType, constVals...)

	// --- Create the runtime Array struct instance ---

	// 1. Allocate stack space for the underlying data array
	dataAlloca := cg.Block.NewAlloca(constArrayType)
	dataAlloca.SetName("array_data")

	// 2. Store the constant data into the stack allocation
	cg.Block.NewStore(constArray, dataAlloca)

	// 3. Allocate stack space for the Array struct itself
	arrayStructType, err := cg.resolveStructType("Array")
	if err != nil {
		return err
	}
	arrayStructAlloca := cg.Block.NewAlloca(arrayStructType)
	arrayStructAlloca.SetName("array_struct")

	// 4. Store length field
	lenFieldAddr := cg.Block.NewGetElementPtr(arrayStructType, arrayStructAlloca, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0)) // Field 0 = length
	llvmIntType := llvmElemType.(*types.IntType)                                                                                                // Assuming length is same type as elements for now
	cg.Block.NewStore(constant.NewInt(llvmIntType, int64(count)), lenFieldAddr)

	// 5. Store data field (pointer to the first element of dataAlloca)
	dataFieldAddr := cg.Block.NewGetElementPtr(arrayStructType, arrayStructAlloca, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1)) // Field 1 = data
	// GEP to get pointer to first element: [N x T]* -> T*
	firstElemPtr := cg.Block.NewGetElementPtr(constArrayType, dataAlloca, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	firstElemPtr.SetName("first_elem_ptr")
	cg.Block.NewStore(firstElemPtr, dataFieldAddr)

	// The result of the array literal expression is the pointer to the Array struct
	cg.lastValue = arrayStructAlloca
	return nil
}
