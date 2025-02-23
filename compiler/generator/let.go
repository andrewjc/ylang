package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitLetStatement(ls *ast.LetStatement) error {
	// We'll pick a default initializer type/value if there's no explicit one.
	var allocaType types.Type = types.I32
	var initValue value.Value = constant.NewInt(types.I32, 0)

	// If there's an explicit initializer, generate code for it.
	if ls.Value != nil {
		if err := ls.Value.Accept(cg); err != nil {
			return err
		}
		if cg.lastValue == nil {
			// fallback if somehow lastValue was not set
			cg.lastValue = constant.NewInt(types.I32, 0)
		}
		initValue = cg.lastValue

		valType := initValue.Type()
		switch t := valType.(type) {
		case *types.FuncType:
			// Bare function type => store a pointer to that function type
			// e.g. i32 () => we alloca i32 ()*
			allocaType = types.NewPointer(t)

		case *types.PointerType:
			// Already a pointer (e.g. i32 ()*).
			// Just use it directly, so the alloca is pointer-to-(that pointer).
			// For example, if it's i32 ()*, we do alloca i32 ()*, meaning the
			// store is (i32 ()*), (i32 ()**) which is valid.
			allocaType = t

		default:
			// Normal first-class type, e.g. i32
			allocaType = t
		}
	}

	// Allocate space for the variable (a stack allocation).
	allocaInst := cg.Block.NewAlloca(allocaType)
	cg.setVar(ls.Name.Value, allocaInst)

	// Store the initializer value if we had one.
	if ls.Value != nil {
		cg.Block.NewStore(initValue, allocaInst)
	}
	return nil
}

func (cg *CodeGenerator) VisitMemberAccessExpression(mae *ast.MemberAccessExpression) error {
	// Minimal stub: generate code for the left, ignore the member for now.
	if err := mae.Left.Accept(cg); err != nil {
		return err
	}
	// We are not producing anything real, so just clear or reuse lastValue
	cg.lastValue = nil
	return nil
}
