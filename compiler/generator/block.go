package generator

import "compiler/ast"

func (cg *CodeGenerator) VisitBlockStatement(bs *ast.BlockStatement) error {
	for _, stmt := range bs.Statements {
		if stmt == nil {
			continue
		}
		if err := stmt.Accept(cg); err != nil {
			return err
		}
	}
	return nil
}
