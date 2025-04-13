// lexer/lexer_loops_test.go
package lexer

import (
	"reflect"
	"testing"
	// No bufio needed here as NewLexerFromString handles it
)

// TestLexer_LoopTokenization verifies tokenization of loop constructs with correct positions.
func TestLexer_LoopTokenization(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []LangToken
		wantErr bool
	}{
		// Line numbers are 0-based, Pos is 1-based STARTING column
		{
			name:  "Test Classic For Loop",
			input: "for (let i = 0; i < 5; i = i + 1) { print(i) }", // All on Line 0
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 1},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 5},
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 6},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 10},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 12},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 14},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 15},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 17},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 19},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 21},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 22},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 24},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 26},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 28},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 30},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 32},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 34},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 36},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 38},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 43},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 44},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 45},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 47},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 48}, // Position where EOF is detected
			},
			wantErr: false,
		},
		{
			name:  "Test Lambda For Loop",
			input: "for (i -> i < 5; i = i + 1) { print(i) }", // Line 0
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 1},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 5},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 6},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 8},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 11},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 13},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 15},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 16},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 18},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 20},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 22},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 24},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 26},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 28},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 30},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 32},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 37},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 38},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 39},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 41},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 42},
			},
			wantErr: false,
		},
		{
			name:  "Test For Each Loop (forEach)",
			input: "for item in range(0, 5) { print(item) }", // Line 0
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 1},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 5},
				{Type: TokenTypeIn, Literal: "in", Line: 0, Pos: 10},
				{Type: TokenTypeRange, Literal: "range", Line: 0, Pos: 13},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 18},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 19},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 20},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 22},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 23},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 25},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 27},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 32},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 33},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 37},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 39},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 40},
			},
			wantErr: false,
		},
		{
			name:  "For Each Lambda Loop (forEachLambda)",
			input: "for item in range(0, 5) -> (item) { print(item) }", // Line 0
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 1},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 5},
				{Type: TokenTypeIn, Literal: "in", Line: 0, Pos: 10},
				{Type: TokenTypeRange, Literal: "range", Line: 0, Pos: 13},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 18},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 19},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 20},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 22},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 23},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 25},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 28},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 29},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 33},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 35},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 37},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 42},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 43},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 47},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 49},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 50},
			},
			wantErr: false,
		},
		{
			name:  "Classic While Loop (whileClassic)",
			input: "let i = 0;while (i < 5) { print(i); i = i + 1 }", // Line 0
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 11},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 17},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 18},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 20},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 22},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 23},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 25},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 27},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 32},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 33},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 34},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 35},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 37},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 39},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 41},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 43},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 45},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 47},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 48},
			},
			wantErr: false,
		},
		{
			name:  "Lambda While Loop (whileLambda)",
			input: "let i = 0;while (i -> i < 5) { print(i); i = i + 1 }", // Line 0
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 11},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 17},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 18},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 20},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 23},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 25},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 27},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 28},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 30},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 32},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 37},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 38},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 39},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 40},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 42},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 44},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 46},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 48},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 50},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 52},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 53},
			},
			wantErr: false,
		},
		{
			name:  "Classic Do-While Loop (doClassic)",
			input: "let i = 0;do { print(i);i = i + 1 } while (i < 5)", // Line 0
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeDo, Literal: "do", Line: 0, Pos: 11},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 14},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 16},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 21},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 22},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 23},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 24},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 25}, // Start Pos is 25 because no space after ;
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 27},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 29},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 31},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 33},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 35},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 37},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 43},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 44},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 46},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 48},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 49},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 50},
			},
			wantErr: false,
		},
		{
			name:  "Lambda Do-While Loop (doLambda)",
			input: "let i = 0;do { print(i); i = i + 1 } while (i -> i < 5)", // Line 0
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeDo, Literal: "do", Line: 0, Pos: 11},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 14},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 16},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 21},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 22},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 23},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 24},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 26}, // Space after ; this time
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 28},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 30},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 32},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 34},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 36},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 38},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 44},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 45},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 47},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 50},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 52},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 54},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 55},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 56},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := NewLexerFromString(tt.input)
			if err != nil {
				t.Errorf("NewLexerFromString() error = %v", err)
				return
			}

			var generatedTokens []LangToken // Collect generated tokens for better debugging
			for i, expected := range tt.want {
				got, err := l.NextToken()
				generatedTokens = append(generatedTokens, got) // Store generated token

				if (err != nil) != tt.wantErr {
					t.Errorf("test[%d] - NextToken() error = %v, wantErr %v", i, err, tt.wantErr)
					t.Logf("Generated Tokens so far: %v", generatedTokens)
					t.Logf("Remaining expected: %v", tt.want[i:])
					break
				}

				// Compare the actual token with the expected token, including Line and Pos
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("test[%d] - NextToken() mismatch.\ngot = {Type:%q, Literal:%q, Line:%d, Pos:%d}\nwant= {Type:%q, Literal:%q, Line:%d, Pos:%d}",
						i, got.Type, got.Literal, got.Line, got.Pos,
						expected.Type, expected.Literal, expected.Line, expected.Pos)
				}

				if got.Type == TokenTypeEOF {
					if i < len(tt.want)-1 {
						t.Errorf("test[%d] - premature EOF. Expected %d more tokens.", i, len(tt.want)-1-i)
						t.Logf("Generated Tokens: %v", generatedTokens)
						t.Logf("Full expected: %v", tt.want)
					}
					break // Stop after EOF
				}
			}

			// Check for extra tokens *after* the loop comparing expected tokens
			extraToken, _ := l.NextToken()
			if extraToken.Type != TokenTypeEOF {
				t.Errorf("NextToken() produced extra token after expected EOF, got = {Type:%q, Literal:%q, Line:%d, Pos:%d}",
					extraToken.Type, extraToken.Literal, extraToken.Line, extraToken.Pos)
				t.Logf("Generated Tokens: %v", generatedTokens)
				t.Logf("Full expected: %v", tt.want)

			}
		})
	}
}
