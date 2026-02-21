package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitIndexExpression(ie *ast.IndexExpression) error {
	// 1. Load the variable holding the Array struct pointer
	savedLHS := cg.inAssignmentLHS
	cg.inAssignmentLHS = true
	err := ie.Left.Accept(cg)
	cg.inAssignmentLHS = savedLHS
	if err != nil {
		return fmt.Errorf("error evaluating base for index expression: %w", err)
	}
	allocaVal := cg.lastValue // the alloca instruction (e.g. %myArr of type %Array**)

	allocaPtrType, ok := allocaVal.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("index expression base alloca is not a pointer, but %T", allocaVal.Type())
	}
	storedType := allocaPtrType.ElemType // type stored in the alloca (e.g. %Array*)

	// Determine the data pointer (i32*) from which to index elements
	var dataPtrVal value.Value

	// Case: alloca stores a pointer-to-struct (%Array**)
	if ptrToStruct, isPtrToStruct := storedType.(*types.PointerType); isPtrToStruct {
		if arrayStructType, isStruct := ptrToStruct.ElemType.(*types.StructType); isStruct {
			// Load %Array** → %Array*
			structPtr := cg.Block.NewLoad(storedType, allocaVal)
			// GEP %Array*, field index 1 (data field) → i32**
			dataFieldAddr := cg.Block.NewGetElementPtr(arrayStructType, structPtr,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, 1))
			// Load i32** → i32*
			dataFieldType := arrayStructType.Fields[1]
			dataPtrVal = cg.Block.NewLoad(dataFieldType, dataFieldAddr)
		}
	}

	// Case: alloca directly stores a struct (%Array*)
	if dataPtrVal == nil {
		if arrayStructType, isStruct := storedType.(*types.StructType); isStruct {
			// GEP %Array, field index 1 → i32**
			dataFieldAddr := cg.Block.NewGetElementPtr(arrayStructType, allocaVal,
				constant.NewInt(types.I32, 0),
				constant.NewInt(types.I32, 1))
			dataFieldType := arrayStructType.Fields[1]
			dataPtrVal = cg.Block.NewLoad(dataFieldType, dataFieldAddr)
		}
	}

	// Case: alloca stores a plain pointer (e.g. i8*, i32*) or an integer
	// treated as a memory address (e.g. i64 from an mmap syscall return).
	if dataPtrVal == nil {
		loadedVal := cg.Block.NewLoad(storedType, allocaVal)
		if _, isPtr := loadedVal.Type().(*types.PointerType); isPtr {
			// Already a pointer — use directly.
			dataPtrVal = loadedVal
		} else if _, isInt := loadedVal.Type().(*types.IntType); isInt {
			// Integer value (e.g. i64 from an mmap syscall return, or i32 holding
			// a small offset). Zero-extend to i64 (preserving the bit pattern
			// of unsigned addresses) then reinterpret as i8*.
			var i64Val value.Value = loadedVal
			if !loadedVal.Type().Equal(types.I64) {
				i64Val = cg.Block.NewZExt(loadedVal, types.I64)
			}
			dataPtrVal = cg.Block.NewIntToPtr(i64Val, types.NewPointer(types.I8))
		} else {
			dataPtrVal = loadedVal
		}
	}

	// 2. Evaluate the index expression
	err = ie.Index.Accept(cg)
	if err != nil {
		return fmt.Errorf("error evaluating index for index expression: %w", err)
	}
	indexVal := cg.lastValue

	// Ensure index is i64 for GEP
	var indexValI64 value.Value
	if !indexVal.Type().Equal(types.I64) {
		indexValI64 = cg.Block.NewSExt(indexVal, types.I64)
	} else {
		indexValI64 = indexVal
	}

	// 3. Determine element type from dataPtrVal
	dataPtrType, ok := dataPtrVal.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("data pointer for index expression is not a pointer type: %T", dataPtrVal.Type())
	}
	elemType := dataPtrType.ElemType

	// 4. GEP to element address, using unique debug names only when the
	//    preferred name is not yet taken (avoids duplicate-value errors for
	//    functions that contain several index expressions).
	elemAddr := cg.Block.NewGetElementPtr(elemType, dataPtrVal, indexValI64)
	cg.trySetName(elemAddr, "elem_addr")

	// 5. Handle LHS vs RHS
	if cg.inAssignmentLHS {
		cg.lastValue = elemAddr
		fmt.Printf("[DEBUG] IndexExpr (LHS): GEP -> %s\n", elemAddr.Ident())
	} else {
		loadedVal := cg.Block.NewLoad(elemType, elemAddr)
		cg.trySetName(loadedVal, "elem_val")
		cg.lastValue = loadedVal
		fmt.Printf("[DEBUG] IndexExpr (RHS): GEP -> %s, Load -> %s\n", elemAddr.Ident(), loadedVal.Ident())
	}

	return nil
}

// very hard coded to arrays for now, next stage is to make this more generic
func (cg *CodeGenerator) VisitMemberAccessExpression(mae *ast.MemberAccessExpression) error {
	// 1. Evaluate the left expression (the object/struct instance) - get alloca
	isLHSOuter := cg.inAssignmentLHS
	cg.inAssignmentLHS = true
	err := mae.Left.Accept(cg)

	if err != nil {
		return fmt.Errorf("error evaluating base for member access '%s': %w", mae.Member, err)
	}
	baseVal := cg.lastValue // Alloca or pointer value

	// Resolve to a pointer-to-struct (%Array*)
	// If the alloca holds a %Array* (%Array**), load it once to get %Array*
	var basePtrVal value.Value
	if ptrType, ok := baseVal.Type().(*types.PointerType); ok {
		if _, isStruct := ptrType.ElemType.(*types.StructType); isStruct {
			// Already %Array* — use directly
			basePtrVal = baseVal
		} else if innerPtr, isInnerPtr := ptrType.ElemType.(*types.PointerType); isInnerPtr {
			if _, isStruct := innerPtr.ElemType.(*types.StructType); isStruct {
				// %Array** — load to get %Array*
				basePtrVal = cg.Block.NewLoad(ptrType.ElemType, baseVal)
			}
		}
	}
	if basePtrVal == nil {
		return fmt.Errorf("member access base does not resolve to a struct pointer, got %T", baseVal.Type())
	}

	structPtrType, ok := basePtrVal.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("member access base is not a pointer, but %T", basePtrVal.Type())
	}
	structType, ok := structPtrType.ElemType.(*types.StructType)
	if !ok {
		return fmt.Errorf("member access base does not point to a struct, but %T", structPtrType.ElemType)
	}

	// 2. Find the field index
	fieldName := mae.Member
	fieldIndex := -1

	// TODO needs proper mapping from AST node to LLVM struct layout.
	// hardcode for Array for now: 0=length, 1=data
	if structType.Name() == "Array" { // Use Name() which we set
		switch fieldName.Value {
		case "length":
			fieldIndex = 0
		case "data":
			fieldIndex = 1
		}
	}

	if fieldIndex == -1 {
		return fmt.Errorf("field '%s' not found in struct type '%s'", fieldName, structType.Name())
	}

	// 3. Generate GetElementPtr (GEP) instruction
	// gep %StructType* %basePtrVal, i32 0, i32 %fieldIndex
	memberAddr := cg.Block.NewGetElementPtr(structType, basePtrVal,
		constant.NewInt(types.I32, 0),                 // Index for the struct
		constant.NewInt(types.I32, int64(fieldIndex)), // Index for the field
	)
	memberAddr.SetName(fieldName.Value + "_addr")

	// 4. Handle LHS vs RHS context
	if isLHSOuter { // If the *overall* expression is LHS (e.g., self.length = ...)
		cg.lastValue = memberAddr // Return the address for storing
		fmt.Printf("[DEBUG] MemberAccess '%s' (LHS): GEP -> %s\n", fieldName, memberAddr.Ident())
	} else { // If RHS (e.g., let x = self.length)
		loadedVal := cg.Block.NewLoad(structType.Fields[fieldIndex], memberAddr)
		loadedVal.SetName(fieldName.Value + "_val")
		cg.lastValue = loadedVal // Return the loaded value
		fmt.Printf("[DEBUG] MemberAccess '%s' (RHS): GEP -> %s, Load -> %s\n", fieldName, memberAddr.Ident(), loadedVal.Ident())
	}

	return nil
}
