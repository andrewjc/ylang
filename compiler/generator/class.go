package generator

import (
	"compiler/ast"
	"fmt"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func (cg *CodeGenerator) VisitClassDeclaration(cd *ast.ClassDeclaration) error {
	typeName := cd.Name.Value
	if _, exists := cg.Structs[typeName]; exists {
		fmt.Printf("[WARN] Type '%s' already defined, skipping definition processing.\n", typeName)
		return nil
	}

	fmt.Printf("[DEBUG] Processing definition for type '%s'\n", typeName)

	// --- Step 1: Determine Field Types ---
	var fieldTypes []types.Type
	var fieldNames []string // For debugging/ordering clarity

	// This requires resolving types from AST nodes (VariableDeclaration, Parameter, etc.)
	// Placeholder: Hardcoding for 'Array' based on stdlib definition
	if typeName == "Array" {
		// Assuming 'int' maps to i32 and '*int' maps to i32*
		llvmIntType, errInt := cg.mapType("int")        // Use mapType to get LLVM type
		llvmIntPtrType, errIntPtr := cg.mapType("*int") // Needs mapType to handle '*'
		if errInt != nil || errIntPtr != nil {
			return fmt.Errorf("failed to map basic types for Array struct: %v, %v", errInt, errIntPtr)
		}
		if llvmIntType == nil || llvmIntPtrType == nil { // Check if mapType returned nil
			return fmt.Errorf("mapType returned nil for Array struct fields")
		}
		fieldTypes = append(fieldTypes, llvmIntType) // length
		fieldNames = append(fieldNames, "length")
		fieldTypes = append(fieldTypes, llvmIntPtrType) // data
		fieldNames = append(fieldNames, "data")
	} else {
		// General case: Iterate cd.Members to find VariableDeclarations
		for _, member := range cd.Members {
			if member.VariableDeclaration != nil {
				varDecl := member.VariableDeclaration
				var fieldName string
				if varDecl.Name != nil {
					fieldName = varDecl.Name.Value
				} else {
					return fmt.Errorf("unnamed field found in type '%s'", typeName) // Fields must have names
				}

				var fieldTypeName string
				if varDecl.Type != nil { // Check for explicit type annotation : Type
					fieldTypeName = varDecl.Type.Value
				} else {
					// TODO: Type inference from varDecl.Value if available?
					// For now, require explicit type annotations for struct fields.
					return fmt.Errorf("type annotation missing for field '%s' in type '%s'", fieldName, typeName)
				}

				llvmFieldType, err := cg.mapType(fieldTypeName)
				if err != nil {
					return fmt.Errorf("could not map type '%s' for field '%s' in type '%s': %w", fieldTypeName, fieldName, typeName, err)
				}
				if llvmFieldType == nil { // Check if mapType returned nil
					return fmt.Errorf("mapType returned nil for field '%s' in type '%s'", fieldName, typeName)
				}
				fieldTypes = append(fieldTypes, llvmFieldType)
				fieldNames = append(fieldNames, fieldName)
			}
			// Ignore MethodDeclarations for field layout
		}
	}

	// --- Step 2: Define the Struct Type in the Module ---
	// Create the struct type with its fields directly.
	structType := types.NewStruct(fieldTypes...)

	// Use NewTypeDef to associate the created struct type with the name in the module.
	// This allows referencing the type by name (e.g., %Array = type { i32, i32* })
	definedType := cg.Module.NewTypeDef(typeName, structType)
	// Note: definedType is essentially the same as structType here for non-opaque cases,
	// but using the result of NewTypeDef is cleaner.

	// --- Step 3: Store in Compiler's Map ---
	// Store the defined struct type (which implements types.Type) in our map.
	cg.Structs[typeName] = definedType // Store the *types.StructType returned by NewTypeDef

	fmt.Printf("[DEBUG] Defined struct type '%s' with fields %v -> %s\n", typeName, fieldNames, definedType)

	// --- Step 4: Process Methods ---
	// Iterate members again to find and generate functions for methods.
	for _, member := range cd.Members {
		if member.MethodDeclaration != nil {
			methodAST := member.MethodDeclaration
			// Pass the *defined* struct type (not the name) to generateMethod if needed,
			// though generateMethod currently re-resolves it via cg.Structs.
			if err := cg.generateMethod(typeName, methodAST); err != nil {
				return fmt.Errorf("error generating method '%s' for type '%s': %w", methodAST.Name.Value, typeName, err)
			}
		}
	}

	return nil
}

func (cg *CodeGenerator) generateMethod(className string, methodAST *ast.MethodDeclaration) error {
	methodName := methodAST.Name.Value
	mangledName := className + "_" + methodName // Simple name mangling

	if _, exists := cg.Functions[mangledName]; exists {
		fmt.Printf("[WARN] Method '%s' already declared/generated, skipping.\n", mangledName)
		return nil
	}

	// 1. Determine Parameter Types (including implicit 'self')
	selfType, ok := cg.Structs[className] // Get the types.Type (should be *types.StructType)
	if !ok {
		return fmt.Errorf("internal error: struct type '%s' not found when generating method '%s'", className, methodName)
	}
	// Ensure selfType is actually a struct type before creating a pointer to it
	_, isStruct := selfType.(*types.StructType)
	if !isStruct {
		return fmt.Errorf("internal error: type '%s' resolved but is not a struct type (%T) when generating method '%s'", className, selfType, methodName)
	}
	selfPtrType := types.NewPointer(selfType) // Pointer to the struct type

	paramTypes := []types.Type{selfPtrType} // First param is 'self' pointer
	paramNames := []string{"self"}

	for _, paramAST := range methodAST.Parameters {
		var paramTypeName string
		if paramAST.Type != nil { // Assuming Parameter AST has Type field
			paramTypeName = paramAST.Type.Value
		} else {
			// Require type annotations for method parameters
			return fmt.Errorf("type annotation missing for parameter '%s' in method '%s'", paramAST.Name.Value, methodName)
		}
		llvmParamType, err := cg.mapType(paramTypeName)
		if err != nil {
			return fmt.Errorf("could not map type '%s' for parameter '%s' in method '%s': %w", paramTypeName, paramAST.Name.Value, methodName, err)
		}
		if llvmParamType == nil { // Check nil from mapType
			return fmt.Errorf("mapType returned nil for parameter '%s' type '%s' in method '%s'", paramAST.Name.Value, paramTypeName, methodName)
		}
		paramTypes = append(paramTypes, llvmParamType)
		paramNames = append(paramNames, paramAST.Name.Value)
	}

	// 2. Determine Return Type
	var retType types.Type = types.Void // Default to void
	if methodAST.ReturnType != nil {
		retTypeName := methodAST.ReturnType.Value
		llvmRetType, err := cg.mapType(retTypeName)
		if err != nil {
			return fmt.Errorf("could not map return type '%s' for method '%s': %w", retTypeName, methodName, err)
		}
		if llvmRetType == nil { // Check nil from mapType
			return fmt.Errorf("mapType returned nil for return type '%s' in method '%s'", retTypeName, methodName)
		}

		// If the YLang code says `-> Array`, we decided to return Array*
		if retTypeName == className {
			retType = selfPtrType // Return pointer to the struct
		} else {
			retType = llvmRetType
		}

	} // else defaults to void

	// 3. Create LLVM Function Signature & Definition
	funcParams := make([]*ir.Param, len(paramTypes))
	for i, pName := range paramNames {
		funcParams[i] = ir.NewParam(pName, paramTypes[i])
	}

	llvmFunc := cg.Module.NewFunc(mangledName, retType, funcParams...)
	cg.Functions[mangledName] = llvmFunc // Store the function

	fmt.Printf("[DEBUG] Declared method '%s' as LLVM function '%s' with signature %s\n", methodName, mangledName, llvmFunc.Sig)

	// 4. Generate Function Body (similar to VisitFunctionDefinition)
	// Store current context
	oldBlock := cg.Block
	oldFunc := cg.currentFunc
	oldVars := cg.Variables
	cg.Variables = make(map[string]value.Value) // New scope for the method

	entry := llvmFunc.NewBlock("entry")
	cg.Block = entry
	cg.currentFunc = llvmFunc

	// Allocate space for parameters ('self' and others) and store initial values
	for _, param := range llvmFunc.Params {
		paramIRName := param.Name()
		alloca := cg.Block.NewAlloca(param.Typ) // Allocate stack space for the parameter value/pointer
		alloca.SetName(paramIRName + ".addr")
		cg.Block.NewStore(param, alloca) // Store the incoming parameter value into the allocation
		cg.setVar(paramIRName, alloca)   // Make the variable name point to the stack allocation
		fmt.Printf("[DEBUG] Method '%s': Allocated and stored parameter '%s' (type %s)\n", mangledName, paramIRName, param.Typ)
	}

	// Visit the method body AST node
	var bodyErr error
	if methodAST.Body != nil {
		bodyErr = methodAST.Body.Accept(cg)
	} else {
		// Method has no body AST node? Add default return.
		fmt.Printf("[WARN] Method '%s' has nil body AST node.\n", mangledName)
	}

	// Add default return if necessary
	if cg.Block != nil && cg.Block.Term == nil {
		if bodyErr == nil {
			funcRetType := cg.currentFunc.Sig.RetType
			if !funcRetType.Equal(types.Void) {
				// If lastValue matches return type, use it? Better require explicit return.
				// Check if the body already returned. If not, add zero return.
				// We are here precisely because cg.Block.Term is nil, so no return happened.
				zero := constant.NewZeroInitializer(funcRetType)
				cg.Block.NewRet(zero)
				fmt.Printf("[DEBUG] Method '%s': Added default zero return (type %s)\n", mangledName, funcRetType)
			} else {
				cg.Block.NewRet(nil) // Return void
				fmt.Printf("[DEBUG] Method '%s': Added default void return\n", mangledName)
			}
		} else {
			fmt.Printf("[WARN] Method '%s': Body errored, potentially leaving block unterminated.\n", mangledName)
			// Still add a terminator? If it's void, maybe. If not, it's invalid IR.
			if cg.currentFunc.Sig.RetType.Equal(types.Void) {
				cg.Block.NewRet(nil)
			} else {
				// Maybe unreachable or ret undef? Let's add Unreachable for now.
				cg.Block.NewUnreachable()
			}
		}
	} else if cg.Block == nil {
		fmt.Printf("[DEBUG] Method '%s': Body resulted in nil block (all paths returned).\n", mangledName)
	}

	// Restore context
	cg.Block = oldBlock
	cg.currentFunc = oldFunc
	cg.Variables = oldVars

	if bodyErr != nil {
		// Don't mask the original body error
		return fmt.Errorf("error generating body for method '%s': %w", mangledName, bodyErr)
	}

	fmt.Printf("[DEBUG] Finished generating body for method '%s'\n", mangledName)
	return nil
}
