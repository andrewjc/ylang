package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	c "compiler/compiler"
	l "compiler/lexer"
	p "compiler/parser"
)

// TestPrimesProgram compiles a Y-lang program that calculates the first 1000
// prime numbers and prints them to stdout, then runs the binary and verifies
// that the output matches the expected primes.
func TestPrimesProgram(t *testing.T) {
	input := `
	import "stdlib/primes";

	main() -> {
		print_primes(1000);
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

	irFile := filepath.Join(tmpDir, "primes.ll")
	if err := os.WriteFile(irFile, []byte(result.Output), 0o644); err != nil {
		t.Fatalf("Failed to write IR file: %v", err)
	}

	// 5. Compile the IR with clang into an executable.
	exeFile := filepath.Join(tmpDir, "primes")
	clangCmd := exec.Command("clang", irFile, "-o", exeFile)
	if out, clangErr := clangCmd.CombinedOutput(); clangErr != nil {
		t.Fatalf("clang compilation failed: %v\nOutput:\n%s", clangErr, string(out))
	}

	// 6. Run the compiled executable.
	runCmd := exec.Command(exeFile)
	output, runErr := runCmd.Output()
	if runErr != nil {
		t.Fatalf("Program execution failed: %v", runErr)
	}

	// 7. Parse the program output into a slice of integers.
	rawLines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var got []int
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, convErr := strconv.Atoi(line)
		if convErr != nil {
			t.Errorf("Output line is not an integer: %q", line)
			continue
		}
		got = append(got, n)
	}

	// 8. Verify the count.
	if len(got) != 1000 {
		t.Fatalf("Expected 1000 primes, got %d\nFirst few: %v\nLast few: %v",
			len(got), got[:minLen(got, 5)], got[maxStart(got, 5):])
	}

	// 9. Verify each prime matches the expected value.
	expected := firstNPrimes(1000)
	for i, want := range expected {
		if got[i] != want {
			t.Errorf("Prime #%d: got %d, want %d", i+1, got[i], want)
		}
	}

	fmt.Printf("First prime: %d, 1000th prime: %d\n", got[0], got[999])
}

// firstNPrimes returns the first n prime numbers using trial division.
func firstNPrimes(n int) []int {
	primes := make([]int, 0, n)
	for candidate := 2; len(primes) < n; candidate++ {
		prime := true
		for _, p := range primes {
			if p*p > candidate {
				break
			}
			if candidate%p == 0 {
				prime = false
				break
			}
		}
		if prime {
			primes = append(primes, candidate)
		}
	}
	return primes
}

func minLen(s []int, n int) int {
	if n > len(s) {
		return len(s)
	}
	return n
}

func maxStart(s []int, n int) int {
	if len(s)-n < 0 {
		return 0
	}
	return len(s) - n
}
