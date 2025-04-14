package lexer

import (
	"bufio"
	"reflect"
	"testing"
)

func TestLexer_LoopTests(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name    string
		input   string
		want    []LangToken
		wantErr bool
	}{
		{
			name:  "Test Classic For Loop",
			input: "for (let i = 0; i < 5; i = i + 1) { print(i) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 4, Length: 1},
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 5, Length: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 9, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 11, Length: 1},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 13, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 14, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 16, Length: 1},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 18, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 20, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 21, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 23, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 25, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 27, Length: 1},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 29, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 31, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 32, Length: 1},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 34, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 36, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 41, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 42, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 43, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 45, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 45, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "Test Lambda For Loop",
			input: "for (i -> i < 5; i = i + 1) { print(i) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 4, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5, Length: 1},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 7, Length: 2},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 10, Length: 1},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 12, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 14, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 15, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 17, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 19, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 21, Length: 1},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 23, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 25, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 26, Length: 1},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 28, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 30, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 35, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 36, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 37, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 39, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 39, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "Test For Each Loop (forEach)",
			input: "for item in range(0, 5) { print(item) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 4, Length: 4},
				{Type: TokenTypeIn, Literal: "in", Line: 0, Pos: 9, Length: 2},
				{Type: TokenTypeRange, Literal: "range", Line: 0, Pos: 12, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 17, Length: 1},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 18, Length: 1},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 19, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 21, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 22, Length: 1},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 24, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 26, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 31, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 32, Length: 4},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 36, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 38, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 38, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "For Each Lambda Loop (forEachLambda)",
			input: "for item in range(0, 5) -> (item) { print(item) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 4, Length: 4},
				{Type: TokenTypeIn, Literal: "in", Line: 0, Pos: 9, Length: 2},
				{Type: TokenTypeRange, Literal: "range", Line: 0, Pos: 12, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 17, Length: 1},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 18, Length: 1},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 19, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 21, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 22, Length: 1},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 24, Length: 2},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 27, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 28, Length: 4},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 32, Length: 1},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 34, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 36, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 41, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 42, Length: 4},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 46, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 48, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 48, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "Classic While Loop (whileClassic)",
			input: "let i = 0;while (i < 5) { print(i); i = i + 1 }",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 4, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 6, Length: 1},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 8, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 9, Length: 1},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 10, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 16, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 17, Length: 1},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 19, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 21, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 22, Length: 1},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 24, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 26, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 31, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 32, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 33, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 34, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 36, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 38, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 40, Length: 1},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 42, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 44, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 46, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 46, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "Lambda While Loop (whileLambda)",
			input: "let i = 0;while (i -> i < 5) { print(i); i = i + 1 }",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 4, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 6, Length: 1},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 8, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 9, Length: 1},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 10, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 16, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 17, Length: 1},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 19, Length: 2},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 22, Length: 1},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 24, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 26, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 27, Length: 1},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 29, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 31, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 36, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 37, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 38, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 39, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 41, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 43, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 45, Length: 1},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 47, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 49, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 51, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 51, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "Classic Do-While Loop (doClassic)",
			input: "let i = 0;do { print(i);i = i + 1 } while (i < 5)",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 4, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 6, Length: 1},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 8, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 9, Length: 1},
				{Type: TokenTypeDo, Literal: "do", Line: 0, Pos: 10, Length: 2},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 13, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 15, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 20, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 21, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 22, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 23, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 24, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 26, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 28, Length: 1},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 30, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 32, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 34, Length: 1},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 36, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 42, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 43, Length: 1},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 45, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 47, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 48, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 48, Length: 0},
			},
			wantErr: false,
		},
		{
			name:  "Lambda Do-While Loop (doLambda)",
			input: "let i = 0;do { print(i); i = i + 1 } while (i -> i < 5)",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 0, Length: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 4, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 6, Length: 1},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 8, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 9, Length: 1},
				{Type: TokenTypeDo, Literal: "do", Line: 0, Pos: 10, Length: 2},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 13, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 15, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 20, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 21, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 22, Length: 1},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 23, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 25, Length: 1},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 27, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 29, Length: 1},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 31, Length: 1},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 33, Length: 1},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 35, Length: 1},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 37, Length: 5},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 43, Length: 1},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 44, Length: 1},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 46, Length: 2},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 49, Length: 1},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 51, Length: 1},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 53, Length: 1},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 54, Length: 1},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 54, Length: 0},
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

			for _, expected := range tt.want {
				got, err := l.NextToken()
				if (err != nil) != tt.wantErr {
					t.Errorf("NextToken() error = %v, wantErr %v", err, tt.wantErr)
					break
				}
				// Check if no tokens are expected and none are returned
				if len(tt.want) == 0 {
					if got != (LangToken{}) {
						t.Errorf("NextToken() unexpected token, got = %v", got)
					}
					continue
				}

				// Compare the actual token with the expected token
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("NextToken() got = %v, want %v", got, expected)
				}
			}
			// Additionally, verify that no extra tokens are produced
			extraToken, _ := l.NextToken()
			if extraToken.Type != TokenTypeEOF {
				t.Errorf("NextToken() produced extra token, got = %v", extraToken)
			}
		})
	}
}
