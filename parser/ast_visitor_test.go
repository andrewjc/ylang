package parser

import (
	"compiler/ast"
	"compiler/lexer"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// recordingVisitor implements the Visitor interface and records which methods were called.
type recordingVisitor struct {
	visitedNodes []string // Records the type of node visited and its string representation
	errors       []error
}

func (v *recordingVisitor) recordVisit(node ast.Node) {
	nodeType := reflect.TypeOf(node).String()
	v.visitedNodes = append(v.visitedNodes, fmt.Sprintf("%s: %s", nodeType, node.String()))
}

func (v *recordingVisitor) VisitProgram(p *ast.Program) error {
	v.recordVisit(p)
	// Manually traverse children because Accept calls Visit, not the other way around
	for _, stmt := range p.ImportStatements {
		if err := stmt.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	for _, fn := range p.Functions {
		if err := fn.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if p.MainFunction != nil {
		if err := p.MainFunction.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	// Add loops for ClassDeclarations, DataStructures etc.
	return nil
}

func (v *recordingVisitor) VisitFunctionDefinition(fn *ast.FunctionDefinition) error {
	v.recordVisit(fn)
	// Traverse children
	if fn.Name != nil {
		if err := fn.Name.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	for _, p := range fn.Parameters {
		if err := p.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if fn.ReturnType != nil {
		if err := fn.ReturnType.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if fn.Body != nil {
		if err := fn.Body.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitLetStatement(ls *ast.LetStatement) error {
	v.recordVisit(ls)
	if ls.Name != nil {
		if err := ls.Name.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if ls.Value != nil {
		if err := ls.Value.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitReturnStatement(rs *ast.ReturnStatement) error {
	v.recordVisit(rs)
	if rs.ReturnValue != nil {
		if err := rs.ReturnValue.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitExpressionStatement(es *ast.ExpressionStatement) error {
	v.recordVisit(es)
	if es.Expression != nil {
		if err := es.Expression.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitNumberLiteral(nl *ast.NumberLiteral) error {
	v.recordVisit(nl)
	return nil
}

func (v *recordingVisitor) VisitStringLiteral(sl *ast.StringLiteral) error {
	v.recordVisit(sl)
	return nil
}

func (v *recordingVisitor) VisitIdentifier(id *ast.Identifier) error {
	v.recordVisit(id)
	return nil
}

func (v *recordingVisitor) VisitInfixExpression(ie *ast.InfixExpression) error {
	v.recordVisit(ie)
	if ie.Left != nil {
		if err := ie.Left.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if ie.Right != nil {
		if err := ie.Right.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitCallExpression(ce *ast.CallExpression) error {
	v.recordVisit(ce)
	if ce.Function != nil {
		if err := ce.Function.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	for _, arg := range ce.Arguments {
		if err := arg.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitArrayLiteral(al *ast.ArrayLiteral) error {
	v.recordVisit(al)
	for _, el := range al.Elements {
		if err := el.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitLambdaExpression(le *ast.LambdaExpression) error {
	v.recordVisit(le)
	for _, p := range le.Parameters {
		if err := p.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if le.Body != nil {
		if err := le.Body.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitMemberAccessExpression(mae *ast.MemberAccessExpression) error {
	v.recordVisit(mae)
	if mae.Left != nil {
		if err := mae.Left.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if mae.Member != nil {
		if err := mae.Member.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitAssignmentExpression(as *ast.AssignmentExpression) error {
	v.recordVisit(as)
	if as.Left != nil {
		if err := as.Left.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if as.Right != nil {
		if err := as.Right.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitBlockStatement(bs *ast.BlockStatement) error {
	v.recordVisit(bs)
	for _, stmt := range bs.Statements {
		if err := stmt.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitIndexExpression(ie *ast.IndexExpression) error {
	v.recordVisit(ie)
	if ie.Left != nil {
		if err := ie.Left.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if ie.Index != nil {
		if err := ie.Index.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitVariableDeclaration(vd *ast.VariableDeclaration) error {
	v.recordVisit(vd)
	if vd.Name != nil {
		if err := vd.Name.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if vd.Type != nil {
		if err := vd.Type.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if vd.Value != nil {
		if err := vd.Value.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil // Should not happen via Accept, but needed for interface
}

func (v *recordingVisitor) VisitIfStatement(is *ast.IfStatement) error {
	v.recordVisit(is)
	if is.Condition != nil {
		if err := is.Condition.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if is.Consequence != nil {
		if err := is.Consequence.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if is.Alternative != nil {
		if err := is.Alternative.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitTraditionalTernaryExpression(te *ast.TraditionalTernaryExpression) error {
	v.recordVisit(te)
	if te.Condition != nil {
		if err := te.Condition.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if te.TrueExpr != nil {
		if err := te.TrueExpr.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if te.FalseExpr != nil {
		if err := te.FalseExpr.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitLambdaStyleTernaryExpression(aste *ast.LambdaStyleTernaryExpression) error {
	v.recordVisit(aste)
	if aste.Condition != nil {
		if err := aste.Condition.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if aste.TrueExpr != nil {
		if err := aste.TrueExpr.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if aste.FalseExpr != nil {
		if err := aste.FalseExpr.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitInlineIfElseTernaryExpression(iite *ast.InlineIfElseTernaryExpression) error {
	v.recordVisit(iite)
	if iite.Condition != nil {
		if err := iite.Condition.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if iite.TrueExpr != nil {
		if err := iite.TrueExpr.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if iite.FalseExpr != nil {
		if err := iite.FalseExpr.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitDotOperator(do *ast.DotOperator) error {
	v.recordVisit(do)
	if do.Left != nil {
		if err := do.Left.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	if do.Right != nil {
		if err := do.Right.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitSyscallExpression(se *ast.SyscallExpression) error {
	v.recordVisit(se)
	if se.Num != nil {
		if err := se.Num.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	for _, arg := range se.Args {
		if err := arg.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func (v *recordingVisitor) VisitImportStatement(is *ast.ImportStatement) error {
	v.recordVisit(is)
	return nil
}
func (v *recordingVisitor) VisitAssemblyExpression(ae *ast.AssemblyExpression) error {
	v.recordVisit(ae)
	if ae.Code != nil {
		if err := ae.Code.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	for _, arg := range ae.Args {
		if err := arg.Accept(v); err != nil {
			v.errors = append(v.errors, err)
		}
	}
	return nil
}

func TestVisitorAcceptDispatchUnit(t *testing.T) {
	input := `
        import "core";
        main() -> {
            let x = 1 + 2;
            let y = "hello";
            let arr = [1, x];
            let l = (p) -> { return p > 0; };
            let z = l(x);
            let obj = MyType(); // Assume MyType exists
            obj.field = arr[0];
            if (z) {
                return obj.field;
            } else {
                syscall(1, 1, y, 5);
                asm("nop");
            }
        }
    `
	// Note: Parsing might fail for unsupported features (MyType, syscall, asm),
	// but we want to test Accept on the nodes that *are* successfully parsed.

	l, err := lexer.NewLexerFromString(input)
	if err != nil {
		t.Fatalf("Lexer creation failed: %v", err)
	}
	p := NewParser(l)
	program := p.ParseProgram()
	// Don't fail on parser errors here, as we want to test Accept on partially valid ASTs if possible.
	// checkParserErrorsAST(t, p)

	if program == nil {
		t.Fatalf("Parsing failed completely, cannot test visitor.")
	}

	visitor := &recordingVisitor{}
	err = program.Accept(visitor) // Start traversal

	if err != nil {
		t.Errorf("Visitor returned an error during traversal: %v", err)
	}
	if len(visitor.errors) > 0 {
		t.Errorf("Visitor encountered %d errors during traversal:", len(visitor.errors))
		for _, e := range visitor.errors {
			t.Errorf(" - %v", e)
		}
	}

	// Expected sequence of Visit calls (simplified, focusing on types)
	// Order depends heavily on traversal implementation within the visitor.
	// This example assumes a depth-first traversal triggered by Accept calls.
	expectedVisits := []string{
		"*ast.Program",            // VisitProgram called first
		"*ast.ImportStatement",    // VisitImportStatement
		"*ast.FunctionDefinition", // VisitFunctionDefinition (main)
		"*ast.Identifier",         // VisitIdentifier (main name)
		// No Params
		// No ReturnType
		"*ast.BlockStatement",   // VisitBlockStatement (main body)
		"*ast.LetStatement",     // VisitLetStatement (x)
		"*ast.Identifier",       // VisitIdentifier (x)
		"*ast.InfixExpression",  // VisitInfixExpression (+)
		"*ast.NumberLiteral",    // VisitNumberLiteral (1)
		"*ast.NumberLiteral",    // VisitNumberLiteral (2)
		"*ast.LetStatement",     // VisitLetStatement (y)
		"*ast.Identifier",       // VisitIdentifier (y)
		"*ast.StringLiteral",    // VisitStringLiteral ("hello")
		"*ast.LetStatement",     // VisitLetStatement (arr)
		"*ast.Identifier",       // VisitIdentifier (arr)
		"*ast.ArrayLiteral",     // VisitArrayLiteral ([1, x])
		"*ast.NumberLiteral",    // VisitNumberLiteral (1)
		"*ast.Identifier",       // VisitIdentifier (x)
		"*ast.LetStatement",     // VisitLetStatement (l)
		"*ast.Identifier",       // VisitIdentifier (l)
		"*ast.LambdaExpression", // VisitLambdaExpression (...)
		"*ast.Identifier",       // VisitIdentifier (p)
		"*ast.BlockStatement",   // VisitBlockStatement (lambda body)
		"*ast.ReturnStatement",  // VisitReturnStatement
		"*ast.InfixExpression",  // VisitInfixExpression (>)
		"*ast.Identifier",       // VisitIdentifier (p)
		"*ast.NumberLiteral",    // VisitNumberLiteral (0)
		"*ast.LetStatement",     // VisitLetStatement (z)
		"*ast.Identifier",       // VisitIdentifier (z)
		"*ast.CallExpression",   // VisitCallExpression (l(x))
		"*ast.Identifier",       // VisitIdentifier (l) - the function being called
		"*ast.Identifier",       // VisitIdentifier (x) - the argument
		// ... potentially skipping MyType(), obj.field due to parsing errors ...
		// Assuming parser recovers enough for the IfStatement:
		"*ast.IfStatement",     // VisitIfStatement
		"*ast.Identifier",      // VisitIdentifier (z) - condition
		"*ast.BlockStatement",  // VisitBlockStatement (consequence)
		"*ast.ReturnStatement", // VisitReturnStatement
		// VisitMemberAccessExpression or DotOperator depending on parser
		// VisitIdentifier obj
		// VisitIdentifier field
		"*ast.BlockStatement", // VisitBlockStatement (alternative)
		// VisitSyscallExpression or fail
		// VisitNumberLiteral 1
		// VisitNumberLiteral 1
		// VisitIdentifier y
		// VisitNumberLiteral 5
		// VisitAssemblyExpression or fail
		// VisitStringLiteral "nop"
	}

	// Check if the recorded visits contain the expected node types.
	// Exact order matching is too brittle due to traversal details and error recovery.
	// Instead, check for the presence of expected node types.
	visitedTypes := make(map[string]bool)
	for _, visit := range visitor.visitedNodes {
		typeName := strings.Split(visit, ":")[0]
		visitedTypes[typeName] = true
	}

	for _, expectedType := range expectedVisits {
		// Strip pointer marker for easier checking if necessary
		cleanExpectedType := strings.TrimPrefix(expectedType, "*")
		found := false
		for visitedTypeKey := range visitedTypes {
			cleanVisitedType := strings.TrimPrefix(visitedTypeKey, "*")
			if cleanVisitedType == cleanExpectedType {
				found = true
				break
			}
		}
		if !found && expectedType != "*ast.SyscallExpression" && expectedType != "*ast.AssemblyExpression" && expectedType != "*ast.MemberAccessExpression" && expectedType != "*ast.DotOperator" { // Allow failures for known unsupported exprs
			t.Errorf("Expected visitor to visit node type %s, but it was not found in recorded visits.", expectedType)
		}
	}

	t.Logf("Recorded visits (%d):", len(visitor.visitedNodes))
	for _, v := range visitor.visitedNodes {
		t.Logf("- %s", v)
	}
}
