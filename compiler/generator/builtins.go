package generator

import (
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type BuiltInManager struct {
	funcs map[string]*ir.Func
}

func NewBuiltInManager(m *ir.Module) *BuiltInManager {
	bm := &BuiltInManager{
		funcs: make(map[string]*ir.Func),
	}
	bm.initBuiltInFuncs(m)
	return bm
}

func (bm *BuiltInManager) GetProvidedFunctionsMap() map[string]*ir.Func {
	return bm.funcs
}

func (bm *BuiltInManager) initBuiltInFuncs(m *ir.Module) {
	// --- Builtin: malloc ---
	mallocSig := types.NewFunc(types.NewPointer(types.I8), types.I64) // void* malloc(size_t size) -> i8* malloc(i64 size)
	mallocFunc := m.NewFunc("malloc", mallocSig.RetType, ir.NewParam("size", mallocSig.Params[0]))
	bm.funcs["malloc"] = mallocFunc // Add to known functions

	// --- Builtin: builtin_print_int(i32) -> void ---
	printIntSig := types.NewFunc(types.Void, types.I32)
	printIntFunc := m.NewFunc("builtin_print_int", printIntSig.RetType, ir.NewParam("val", printIntSig.Params[0]))
	printIntEntry := printIntFunc.NewBlock("entry")
	printIntEntry.NewRet(nil) // Return void
	bm.funcs["builtin_print_int"] = printIntFunc

	// --- Builtin: builtin_print_newline() -> void ---
	printNewlineSig := types.NewFunc(types.Void)
	printNewlineFunc := m.NewFunc("builtin_print_newline", printNewlineSig.RetType)
	printNewlineEntry := printNewlineFunc.NewBlock("entry")
	printNewlineEntry.NewRet(nil) // Return void
	bm.funcs["builtin_print_newline"] = printNewlineFunc

	// --- Array builtins just placeholders for now ---
	arrayMapRetType := types.NewPointer(types.I32)
	arrayMapFunc := m.NewFunc("builtin_array_map", arrayMapRetType)
	bm.funcs["builtin_array_map"] = arrayMapFunc

	// ForEach: Takes array ptr, callback ptr. Returns void.
	arrayForEachRetType := types.Void
	arrayForEachFunc := m.NewFunc("builtin_array_forEach", arrayForEachRetType)
	bm.funcs["builtin_array_forEach"] = arrayForEachFunc
}

// generateArrayMap generates LLVM IR for array.map(callback)
// objPtrVal: Pointer to the array (e.g., [N x T]*)
// arrayType: The underlying array type ([N x T])
// args: LLVM values for arguments passed to map (should contain the callback function)
func (cg *CodeGenerator) generateArrayMap(objPtrVal value.Value, arrayType *types.ArrayType, args []value.Value) error {
	if len(args) != 1 {
		return fmt.Errorf("array.map expects exactly 1 argument (callback function), got %d", len(args))
	}
	callbackFnVal := args[0]

	// Ensure callbackFnVal is a callable function pointer
	callbackFnPtrType, ok := callbackFnVal.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("argument to map is not a function pointer type: %s", callbackFnVal.Type())
	}
	callbackFnSig, ok := callbackFnPtrType.ElemType.(*types.FuncType)
	if !ok {
		return fmt.Errorf("argument to map is not a function pointer type: %s", callbackFnVal.Type())
	}

	// --- Basic Type/Size Info ---
	arrayLen := arrayType.Len
	elemType := arrayType.ElemType          // Type of elements in the *source* array
	resultElemType := callbackFnSig.RetType // Type of elements in the *result* array

	if len(callbackFnSig.Params) != 1 || !callbackFnSig.Params[0].Equal(elemType) {
		return fmt.Errorf("map callback signature mismatch: expected func(%s) %s, got %s", elemType, resultElemType, callbackFnSig)
	}

	// --- Allocate Result Array ---
	resultArrayType := types.NewArray(arrayLen, resultElemType)
	resultArrayAlloca := cg.Block.NewAlloca(resultArrayType)
	resultArrayAlloca.SetName("map_result_arr")

	// --- Loop Setup ---
	// We need an index variable
	indexAlloca := cg.Block.NewAlloca(types.I64) // Use i64 for index
	indexAlloca.SetName("map_idx")
	cg.Block.NewStore(constant.NewInt(types.I64, 0), indexAlloca)

	loopCondBlock := cg.currentFunc.NewBlock("map_loop_cond")
	loopBodyBlock := cg.currentFunc.NewBlock("map_loop_body")
	loopEndBlock := cg.currentFunc.NewBlock("map_loop_end")

	cg.Block.NewBr(loopCondBlock) // Jump to condition check first

	// --- Loop Condition ---
	cg.Block = loopCondBlock
	currentIndex := cg.Block.NewLoad(types.I64, indexAlloca)
	currentIndex.SetName("map_current_idx")
	cond := cg.Block.NewICmp(enum.IPredULT, currentIndex, constant.NewInt(types.I64, int64(arrayLen))) // Use unsigned less than
	cond.SetName("map_loop_cond_check")
	cg.Block.NewCondBr(cond, loopBodyBlock, loopEndBlock)

	// --- Loop Body ---
	cg.Block = loopBodyBlock
	currentIndexForBody := cg.Block.NewLoad(types.I64, indexAlloca) // Reload index for use in GEP
	currentIndexForBody.SetName("map_idx_body")

	// Get pointer to element in SOURCE array: gep [N x T]* %objPtrVal, i64 0, i64 %currentIndexForBody
	srcElemPtr := cg.Block.NewGetElementPtr(arrayType, objPtrVal,
		constant.NewInt(types.I64, 0), // First index for pointer deref
		currentIndexForBody,           // Second index for array element
	)
	srcElemPtr.SetName("map_src_elem_ptr")

	// Load the element from the source array
	srcElemVal := cg.Block.NewLoad(elemType, srcElemPtr)
	srcElemVal.SetName("map_src_elem_val")

	// Call the callback function
	callResult := cg.Block.NewCall(callbackFnVal, srcElemVal)
	callResult.SetName("map_callback_res")

	// Get pointer to element in RESULT array: gep [N x ResT]* %resultArrayAlloca, i64 0, i64 %currentIndexForBody
	destElemPtr := cg.Block.NewGetElementPtr(resultArrayType, resultArrayAlloca,
		constant.NewInt(types.I64, 0),
		currentIndexForBody,
	)
	destElemPtr.SetName("map_dest_elem_ptr")

	// Store the result in the destination array
	cg.Block.NewStore(callResult, destElemPtr)

	// Increment index: index = index + 1
	nextIndex := cg.Block.NewAdd(currentIndexForBody, constant.NewInt(types.I64, 1))
	nextIndex.SetName("map_next_idx")
	cg.Block.NewStore(nextIndex, indexAlloca)

	cg.Block.NewBr(loopCondBlock) // Jump back to condition

	// --- Loop End ---
	cg.Block = loopEndBlock
	// The result of map is the pointer to the new array
	cg.lastValue = resultArrayAlloca
	return nil
}

// generateArrayForEach generates LLVM IR for array.forEach(callback)
// objPtrVal: Pointer to the array (e.g., [N x T]*)
// arrayType: The underlying array type ([N x T])
// args: LLVM values for arguments passed to forEach (should contain the callback function)
func (cg *CodeGenerator) generateArrayForEach(objPtrVal value.Value, arrayType *types.ArrayType, args []value.Value) error {
	if len(args) != 1 {
		return fmt.Errorf("array.forEach expects exactly 1 argument (callback function), got %d", len(args))
	}
	callbackFnVal := args[0]

	// Ensure callbackFnVal is a callable function pointer
	callbackFnPtrType, ok := callbackFnVal.Type().(*types.PointerType)
	if !ok {
		return fmt.Errorf("argument to forEach is not a function pointer type: %s", callbackFnVal.Type())
	}
	callbackFnSig, ok := callbackFnPtrType.ElemType.(*types.FuncType)
	if !ok {
		return fmt.Errorf("argument to forEach is not a function pointer type: %s", callbackFnVal.Type())
	}

	// --- Basic Type/Size Info ---
	arrayLen := arrayType.Len
	elemType := arrayType.ElemType

	// Check signature: Expects func(ElemType) void or func(ElemType) T
	if len(callbackFnSig.Params) != 1 || !callbackFnSig.Params[0].Equal(elemType) {
		return fmt.Errorf("forEach callback signature mismatch: expected func(%s), got %s", elemType, callbackFnSig)
	}

	// --- Loop Setup ---
	indexAlloca := cg.Block.NewAlloca(types.I64)
	indexAlloca.SetName("fe_idx")
	cg.Block.NewStore(constant.NewInt(types.I64, 0), indexAlloca)

	loopCondBlock := cg.currentFunc.NewBlock("fe_loop_cond")
	loopBodyBlock := cg.currentFunc.NewBlock("fe_loop_body")
	loopEndBlock := cg.currentFunc.NewBlock("fe_loop_end")

	cg.Block.NewBr(loopCondBlock)

	// --- Loop Condition ---
	cg.Block = loopCondBlock
	currentIndex := cg.Block.NewLoad(types.I64, indexAlloca)
	currentIndex.SetName("fe_current_idx")
	cond := cg.Block.NewICmp(enum.IPredULT, currentIndex, constant.NewInt(types.I64, int64(arrayLen)))
	cond.SetName("fe_loop_cond_check")
	cg.Block.NewCondBr(cond, loopBodyBlock, loopEndBlock)

	// --- Loop Body ---
	cg.Block = loopBodyBlock
	currentIndexForBody := cg.Block.NewLoad(types.I64, indexAlloca) // Reload index
	currentIndexForBody.SetName("fe_idx_body")

	// Get pointer to element: gep [N x T]* %objPtrVal, i64 0, i64 %currentIndexForBody
	elemPtr := cg.Block.NewGetElementPtr(arrayType, objPtrVal,
		constant.NewInt(types.I64, 0),
		currentIndexForBody,
	)
	elemPtr.SetName("fe_elem_ptr")

	// Load the element
	elemVal := cg.Block.NewLoad(elemType, elemPtr)
	elemVal.SetName("fe_elem_val")

	// Call the callback function
	cg.Block.NewCall(callbackFnVal, elemVal) // Ignore return value if any

	// Increment index
	nextIndex := cg.Block.NewAdd(currentIndexForBody, constant.NewInt(types.I64, 1))
	nextIndex.SetName("fe_next_idx")
	cg.Block.NewStore(nextIndex, indexAlloca)

	cg.Block.NewBr(loopCondBlock)

	// --- Loop End ---
	cg.Block = loopEndBlock
	// forEach usually doesn't return a meaningful value for chaining in this style,
	// but to allow obj.forEach(...).anotherMethod(...) we might return the original object ptr.
	cg.lastValue = objPtrVal
	return nil
}
