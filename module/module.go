package module

import "compiler/ast"

// Module represents a module in the compiler.
type Module struct {
	Name string

	AST *ast.Program

	// Path is the path to the module.
	Path string

	// Compiled is true if the module has been compiled.
	Compiled bool

	// Imports is a list of modules that this module imports.
	Imports []*Module

	// ImportStatements is a list of import statements in the module.
	ImportStatements []*ast.ImportStatement

	// TopLevelItems is a list of top-level items in the module.
	TopLevelItems []ast.Node

	// Errors is a list of errors that occurred while compiling the module.
	Errors []string

	// The module manager that manages this module.
	ModuleManager *ModuleManager
}
