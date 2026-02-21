package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	c "compiler/compiler"
	l "compiler/lexer"
	p "compiler/parser"
)

// TestListFilesProgram specifies a y-lang program that lists all files in the
// working directory, compiles and runs it, then verifies the output matches
// the actual directory contents.
//
// The program is implemented entirely in Y-lang using Linux syscalls — no
// external C runtime or stubs are required.
func TestListFilesProgram(t *testing.T) {
	input := `
	import "stdlib/fs";

	main() -> {
		listdir();
		return 0;
	}`

	// 1. Lex the program.
	lexer, err := l.NewLexerFromString(input)
	if err != nil {
		t.Fatalf("Failed to create lexer: %v", err)
	}

	// 2. Parse the program.
	parser := p.NewParser(lexer)
	program := parser.ParseProgram()
	if len(parser.Errors()) != 0 {
		for i, parseErr := range parser.Errors() {
			t.Logf("Parser error %d: %v", i+1, parseErr)
		}
		t.FailNow()
	}

	// 3. Compile to LLVM IR.
	compilerInstance := c.NewCompiler(c.LLVM)
	result := compilerInstance.Compile(program)
	if len(result.Errors) != 0 {
		t.Fatalf("Compiler errors: %v", result.Errors)
	}
	if result.Output == "" {
		t.Fatal("Compiler produced no output")
	}
	t.Logf("Generated LLVM IR:\n%s", result.Output)

	// 4. Write the IR to a temporary directory.
	tmpDir := t.TempDir()

	irFile := filepath.Join(tmpDir, "listfiles.ll")
	if err := os.WriteFile(irFile, []byte(result.Output), 0o644); err != nil {
		t.Fatalf("Failed to write IR file: %v", err)
	}

	// 5. Compile IR with clang into an executable.
	//    No C runtime file is needed — the program is fully self-contained
	//    (all I/O goes through inline-asm syscall instructions).
	exeFile := filepath.Join(tmpDir, "listfiles")
	clangCmd := exec.Command("clang", irFile, "-o", exeFile)
	if out, err := clangCmd.CombinedOutput(); err != nil {
		t.Fatalf("clang compilation failed: %v\nOutput:\n%s", err, string(out))
	}

	// 6. Run the compiled executable in the current working directory.
	runCmd := exec.Command(exeFile)
	output, err := runCmd.Output()
	if err != nil {
		t.Fatalf("Program execution failed: %v", err)
	}

	// 7. Collect expected file names from the working directory.
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("Failed to read working directory: %v", err)
	}
	expected := make([]string, 0, len(entries))
	for _, e := range entries {
		expected = append(expected, e.Name())
	}
	sort.Strings(expected)

	// 8. Parse and sort the actual program output.
	rawLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	actual := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		if line != "" {
			actual = append(actual, line)
		}
	}
	sort.Strings(actual)

	// 9. Verify the output matches the expected directory contents.
	if len(actual) != len(expected) {
		t.Errorf("File count mismatch: program listed %d files, directory has %d\nActual:   %v\nExpected: %v",
			len(actual), len(expected), actual, expected)
		return
	}
	for i, name := range expected {
		if actual[i] != name {
			t.Errorf("File name mismatch at index %d: got %q, want %q", i, actual[i], name)
		}
	}
}

