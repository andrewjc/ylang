package generator

import (
	"fmt"
	"github.com/llir/llvm/ir/types"
	"strings"
)

// Ensure mapType can handle pointers and basic types including void.
// Needs to look up defined struct types from cg.Structs.
func (cg *CodeGenerator) mapType(typeName string) (types.Type, error) {
	// Handle pointer types
	if strings.HasSuffix(typeName, "*") {
		baseTypeName := strings.TrimSuffix(typeName, "*")
		baseType, err := cg.mapType(baseTypeName) // Recursive call
		if err != nil {
			return nil, fmt.Errorf("unknown base type '%s' for pointer type '%s'", baseTypeName, typeName)
		}
		return types.NewPointer(baseType), nil
	}

	// Handle primitive types
	switch typeName {
	case "int":
		// Let's use i32 consistently for now, matches test cases better.
		// Consider i64 for sizes/indices if needed later.
		return types.I32, nil
	case "float":
		return types.Float, nil
	case "bool":
		return types.I1, nil
	case "string":
		// Represent string as char* (i8*)
		return types.NewPointer(types.I8), nil
	case "void": // For function return types
		return types.Void, nil
	// Add other primitives like uint8, i64 etc. as needed
	case "i8":
		return types.I8, nil
	case "i16":
		return types.I16, nil
	case "i32":
		return types.I32, nil
	case "i64":
		return types.I64, nil

	default:
		// Check if it's a user-defined struct type we've already processed
		if definedType, exists := cg.Structs[typeName]; exists {
			return definedType, nil
		}

		// If not found, it's an error at this stage.
		// Type resolution should ideally happen in passes or before codegen.
		return nil, fmt.Errorf("unsupported or undefined type: %s", typeName)
	}
}
