package generator

import "compiler/ast"

func (cg *CodeGenerator) VisitImportStatement(is *ast.ImportStatement) error {
	// load the module's AST, then compile it if not compiled
	mod, err := cg.ModuleManager.LoadModule(is.Path)
	if err != nil {
		return err
	}
	// Now visit the module’s AST to generate IR for all its top-level items
	return mod.AST.Accept(cg)
}
