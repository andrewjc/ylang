package module

import (
	"compiler/lexer"
	"compiler/parser"
	"fmt"
	"os"
	"path/filepath"
)

type ModuleManager struct {

	// Modules is a map of module paths to modules.
	modules map[string]*Module

	searchPaths []string
}

func NewModuleManager() *ModuleManager {
	return &ModuleManager{
		modules: make(map[string]*Module),
		// e.g. searchPaths could default to
		searchPaths: []string{"./", "./lib", "./contrib", "./vendor"},
	}
}

func (mm *ModuleManager) LoadModule(modulePath string) (*Module, error) {
	// If already loaded, return it:
	if mod, ok := mm.modules[modulePath]; ok {
		return mod, nil
	}

	// otherwise, locate the file
	filePath, err := mm.findModuleFile(modulePath)
	if err != nil {
		return nil, err
	}

	// parse it
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lex, err := lexer.NewLexerFromString(string(src))
	if err != nil {
		return nil, err
	}

	p := parser.NewParser(lex)
	astProg := p.ParseProgram()
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parser errors in %s: %v", modulePath, p.Errors())
	}

	mod := &Module{
		Name: modulePath,
		AST:  astProg,
	}
	mm.modules[modulePath] = mod
	return mod, nil
}

// Simplistic approach
func (mm *ModuleManager) findModuleFile(modulePath string) (string, error) {
	// If it’s "std/core", maybe we search "std/core.y" or something
	for _, sp := range mm.searchPaths {
		candidate := filepath.Join(sp, modulePath+".y")
		// if that file exists, return candidate
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		// or maybe its a directory with an index.y file
		moduleName := filepath.Base(modulePath)
		candidate = filepath.Join(sp, modulePath, moduleName+".y")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("module %s not found in search paths", modulePath)
}
