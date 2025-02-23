package generator

import (
	"fmt"
	"github.com/llir/llvm/ir/types"
)

func (cg *CodeGenerator) mapType(typeName string) (types.Type, error) {
	switch typeName {
	case "int":
		return types.I32, nil
	case "uint8":
		return types.I8, nil
	case "uint16":
		return types.I16, nil
	case "uint32":
		return types.I32, nil
	case "uint64":
		return types.I64, nil
	case "float":
		return types.Float, nil
	case "double":
		return types.Double, nil
	case "bool":
		return types.I1, nil
	case "string":
		// Represent string as a pointer to i8 (C-style string)
		return types.NewPointer(types.I8), nil
	default:
		// Check if it's a user-defined struct
		if structType, exists := cg.Structs[typeName]; exists {
			return structType, nil
		}
		// If not defined yet, attempt to define it based on AST
		// This requires access to AST; adjust as necessary
		return nil, fmt.Errorf("unsupported or undefined type: %s", typeName)
	}
}
