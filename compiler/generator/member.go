package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitIndexExpression(ie *ast.IndexExpression) error {
	// 1. Evaluate the base pointer expression (e.g., self.data)
	err := ie.Left.Accept(cg)
	if err != nil {
		return fmt.Errorf("error evaluating base for index expression: %w", err)
	}
	basePtrVal := cg.lastValue // Should be a pointer type, e.g., i32*

	basePtrType, ok := basePtrVal.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("index expression base is not a pointer, but %T", basePtrVal.Type())
	}
	elemType := basePtrType.ElemType // The type of elements being pointed to, e.g., i32

	// 2. Evaluate the index expression
	err = ie.Index.Accept(cg)
	if err != nil {
		return fmt.Errorf("error evaluating index for index expression: %w", err)
	}
	indexVal := cg.lastValue // Should be an integer type

	// Ensure index is compatible type (e.g. i64 or i32, GEP usually wants i64 or i32)
	// LLVM GEP often uses i64 or i32 indices. Let's target i64 for flexibility.
	var indexValI64 value.Value
	if !indexVal.Type().Equal(types.I64) {
		indexValI64 = cg.Block.NewSExt(indexVal, types.I64) // Or ZExt? Assume signed.
	} else {
		indexValI64 = indexVal
	}

	// 3. Generate GetElementPtr (GEP)
	// gep %ElemType* %basePtrVal, i64 %indexValI64
	elemAddr := cg.Block.NewGetElementPtr(elemType, basePtrVal, indexValI64)
	elemAddr.SetName("elem_addr")

	// 4. Handle LHS vs RHS
	if cg.inAssignmentLHS { // e.g., array[i] = value
		cg.lastValue = elemAddr // Return address for store
		fmt.Printf("[DEBUG] IndexExpr (LHS): GEP -> %s\n", elemAddr.Ident())
	} else { // e.g., let x = array[i]
		loadedVal := cg.Block.NewLoad(elemType, elemAddr)
		loadedVal.SetName("elem_val")
		cg.lastValue = loadedVal // Return loaded value
		fmt.Printf("[DEBUG] IndexExpr (RHS): GEP -> %s, Load -> %s\n", elemAddr.Ident(), loadedVal.Ident())
	}

	return nil
}

// very hard coded to arrays for now, next stage is to make this more generic
func (cg *CodeGenerator) VisitMemberAccessExpression(mae *ast.MemberAccessExpression) error {
	// 1. Evaluate the left expression (the object/struct instance)
	isLHSOuter := cg.inAssignmentLHS
	cg.inAssignmentLHS = true
	err := mae.Left.Accept(cg)

	if err != nil {
		return fmt.Errorf("error evaluating base for member access '%s': %w", mae.Member, err)
	}
	basePtrVal := cg.lastValue // Should be a pointer to a struct, e.g. %Array*

	basePtrType, ok := basePtrVal.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("member access base is not a pointer, but %T", basePtrVal.Type())
	}
	structType, ok := basePtrType.ElemType.(*types.StructType)
	if !ok {
		return fmt.Errorf("member access base does not point to a struct, but %T", basePtrType.ElemType)
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
