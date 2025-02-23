package generator

import (
	"compiler/ast"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) VisitReturnStatement(rs *ast.ReturnStatement) error {
	if rs.ReturnValue != nil {
		if err := rs.ReturnValue.Accept(cg); err != nil {
			return err
		}
		if cg.lastValue != nil {
			cg.Block.NewRet(cg.lastValue)
			return nil
		}
	}
	// Default return 0 if no expression or no lastValue
	cg.Block.NewRet(constant.NewInt(types.I32, 0))
	return nil
}
