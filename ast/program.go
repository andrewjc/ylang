package ast

type Program struct {
	MainFunction      *FunctionDefinition
	ClassDeclarations []*ClassDeclaration
	Functions         []*FunctionDefinition
	DataStructures    []*DataStructure
}
