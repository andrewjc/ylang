// lexer/lexer_asm_test.go
package lexer

import (
	"reflect"
	"testing"
	// No bufio needed, using NewLexerFromString
)

// TestLexer_AssemblyTokenization tests that NextToken correctly identifies the 'asm' keyword
// and related tokens in assembly constructs.
func TestLexer_AssemblyTokenization(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedTokens []LangToken // Expect the sequence of tokens for the asm construct
		wantErr        bool
	}{
		{
			name:  "Simple asm call",
			input: `asm("mov eax, ebx")`, // Common syntax: asm keyword, parens, string
			expectedTokens: []LangToken{
				{Type: TokenTypeAssembly, Literal: "asm", Line: 0, Pos: 0, Length: 3}, // 'asm' is identified first
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 3, Length: 1},
				{Type: TokenTypeString, Literal: "mov eax, ebx", Line: 0, Pos: 4, Length: 12},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 18, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 18, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "asm call with args (lexing only)",
			input: `asm("syscall", 1, fd, buf)`, // Lexer just tokenizes, parser handles args
			expectedTokens: []LangToken{
				{Type: TokenTypeAssembly, Literal: "asm", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 3, Length: 1},
				{Type: TokenTypeString, Literal: "syscall", Line: 0, Pos: 4, Length: 7},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 13, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 15, Length: 1},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 16, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "fd", Line: 0, Pos: 18, Length: 2},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 20, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "buf", Line: 0, Pos: 22, Length: 3},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 25, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 25, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "asm not as keyword (part of identifier)",
			input: `asm = 1`,
			expectedTokens: []LangToken{
				{Type: TokenTypeIdentifier, Literal: "asm", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 12, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 15},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 16},
			},
			wantErr: false,
		},
		{
			name:  "asm with immediate whitespace",
			input: `asm ("nop")`,
			expectedTokens: []LangToken{
				{Type: TokenTypeAssembly, Literal: "asm", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 4, Length: 1}, // Skips space
				{Type: TokenTypeString, Literal: "nop", Line: 0, Pos: 5, Length: 3},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 10, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 10, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "Simple asm call (Corrected based on LookupIdent)",
			input: `asm("mov eax, ebx")`,
			expectedTokens: []LangToken{
				{Type: TokenTypeAssembly, Literal: "asm", Line: 0, Pos: 0, Length: 3}, // Should be TokenTypeAssembly now
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 3, Length: 1},
				{Type: TokenTypeString, Literal: "mov eax, ebx", Line: 0, Pos: 4, Length: 12},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 18, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 18, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "asm call with args (Corrected based on LookupIdent)",
			input: `asm("syscall", 1, fd, buf)`,
			expectedTokens: []LangToken{
				{Type: TokenTypeAssembly, Literal: "asm", Line: 0, Pos: 0, Length: 3}, // Should be TokenTypeAssembly
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 3, Length: 1},
				{Type: TokenTypeString, Literal: "syscall", Line: 0, Pos: 4, Length: 7},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 13, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 15, Length: 1},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 16, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "fd", Line: 0, Pos: 18, Length: 2},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 20, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "buf", Line: 0, Pos: 22, Length: 3},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 25, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 25, Length: 0},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use NewLexerFromString for initialization
			l, err := NewLexerFromString(tt.input)
			if err != nil {
				t.Fatalf("NewLexerFromString() error = %v", err)
				return
			}

			generatedTokens := []LangToken{}
			for i, expected := range tt.expectedTokens {
				got, err := l.NextToken()
				generatedTokens = append(generatedTokens, got)

				if (err != nil) != tt.wantErr {
					t.Errorf("test[%d] - NextToken() error = %v, wantErr %v", i, err, tt.wantErr)
					t.Logf("Generated Tokens so far: %v", generatedTokens)
					t.Logf("Remaining expected: %v", tt.expectedTokens[i:])
					break
				}

				// Compare the actual token with the expected token
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("test[%d] - NextToken() mismatch.\ngot = {Type:%q, Literal:%q, Line:%d, Pos:%d, Len=%d}\nwant= {Type:%q, Literal:%q, Line:%d, Pos:%d, Len: %d}",
						i, got.Type, got.Literal, got.Line, got.Pos, got.Length,
						expected.Type, expected.Literal, expected.Line, expected.Pos, expected.Length)
				}

				if got.Type == TokenTypeEOF {
					if i < len(tt.expectedTokens)-1 {
						t.Errorf("test[%d] - premature EOF. Expected %d more tokens.", i, len(tt.expectedTokens)-1-i)
						t.Logf("Generated Tokens: %v", generatedTokens)
						t.Logf("Full expected: %v", tt.expectedTokens)
					}
					break // Stop after EOF
				}
			}

			// Check for extra tokens
			extraToken, _ := l.NextToken()
			if extraToken.Type != TokenTypeEOF {
				t.Errorf("NextToken() produced extra token after expected EOF, got = {Type:%q, Literal:%q, Line:%d, Pos:%d}",
					extraToken.Type, extraToken.Literal, extraToken.Line, extraToken.Pos)
				t.Logf("Generated Tokens: %v", generatedTokens)
				t.Logf("Full expected: %v", tt.expectedTokens)
			}
		})
	}
}
