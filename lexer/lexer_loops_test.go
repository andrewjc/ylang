package lexer

import (
	"bufio"
	"reflect"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
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
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 3},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 5}, // Skip space
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 9},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 11},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 13},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 15},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 16},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 18},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 20},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 22},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 23},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 25},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 27},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 29},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 31},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 33},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 35},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 37},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 43},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 44},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 45},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 47},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 49},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 50}, // Pos after last char
			},
			wantErr: false,
		},
		{
			name:  "Test Lambda For Loop",
			input: "for (i -> i < 5; i = i + 1) { print(i) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 3},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 5},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 7},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 10},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 12},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 14},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 16},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 17},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 19},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 21},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 23},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 25},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 27},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 29},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 31},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 37},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 38},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 39},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 41},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 43},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 44},
			},
			wantErr: false,
		},
		{
			name:  "Test For Each Loop (forEach)",
			input: "for item in range(0, 5) { print(item) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 3},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 8},
				{Type: TokenTypeIn, Literal: "in", Line: 0, Pos: 11},
				{Type: TokenTypeRange, Literal: "range", Line: 0, Pos: 17},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 18},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 19},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 20},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 22},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 24},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 26},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 32},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 33},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 37},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 39},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 41},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 42},
			},
			wantErr: false,
		},
		{
			name:  "For Each Lambda Loop (forEachLambda)",
			input: "for item in range(0, 5) -> (item) { print(item) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for", Line: 0, Pos: 3},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 8},
				{Type: TokenTypeIn, Literal: "in", Line: 0, Pos: 11},
				{Type: TokenTypeRange, Literal: "range", Line: 0, Pos: 17},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 18},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 19},
				{Type: TokenTypeComma, Literal: ",", Line: 0, Pos: 20},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 22},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 24},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 27},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 29},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 33},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 35},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 37},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 43},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 44},
				{Type: TokenTypeIdentifier, Literal: "item", Line: 0, Pos: 48},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 50},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 52},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 53},
			},
			wantErr: false,
		},
		{
			name:  "Classic While Loop (whileClassic)",
			input: "let i = 0;while (i < 5) { print(i); i = i + 1 }",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 15},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 17},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 19},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 21},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 23},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 25},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 27},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 33},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 34},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 35},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 37},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 38},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 40},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 42},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 44},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 46},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 48},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 50},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 51},
			},
			wantErr: false,
		},
		{
			name:  "Lambda While Loop (whileLambda)",
			input: "let i = 0;while (i -> i < 5) { print(i); i = i + 1 }",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 15},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 17},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 19},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 22},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 24},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 26},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 28},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 30},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 32},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 38},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 39},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 40},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 42},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 43},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 45},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 47},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 49},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 51},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 53},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 55},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 56},
			},
			wantErr: false,
		},
		{
			name:  "Classic Do-While Loop (doClassic)",
			input: "let i = 0;do { print(i);i = i + 1 } while (i < 5)",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeDo, Literal: "do", Line: 0, Pos: 13},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 15},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 21},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 22},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 23},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 25},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 26},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 27}, // Note: No space before i
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 29},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 31},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 33},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 35},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 37},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 43},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 45},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 47},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 49},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 51},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 53},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 54},
			},
			wantErr: false,
		},
		{
			name:  "Lambda Do-While Loop (doLambda)",
			input: "let i = 0;do { print(i); i = i + 1 } while (i -> i < 5)",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let", Line: 0, Pos: 3},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 5},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 7},
				{Type: TokenTypeNumber, Literal: "0", Line: 0, Pos: 9},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 10},
				{Type: TokenTypeDo, Literal: "do", Line: 0, Pos: 13},
				{Type: TokenTypeLeftBrace, Literal: "{", Line: 0, Pos: 15},
				{Type: TokenTypeIdentifier, Literal: "print", Line: 0, Pos: 21},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 22},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 23},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 25},
				{Type: TokenTypeSemicolon, Literal: ";", Line: 0, Pos: 26},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 28},
				{Type: TokenTypeAssignment, Literal: "=", Line: 0, Pos: 30},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 32},
				{Type: TokenTypePlus, Literal: "+", Line: 0, Pos: 34},
				{Type: TokenTypeNumber, Literal: "1", Line: 0, Pos: 36},
				{Type: TokenTypeRightBrace, Literal: "}", Line: 0, Pos: 38},
				{Type: TokenTypeWhile, Literal: "while", Line: 0, Pos: 44},
				{Type: TokenTypeLeftParenthesis, Literal: "(", Line: 0, Pos: 46},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 48},
				{Type: TokenTypeLambdaArrow, Literal: "->", Line: 0, Pos: 51},
				{Type: TokenTypeIdentifier, Literal: "i", Line: 0, Pos: 53},
				{Type: TokenTypeLessThan, Literal: "<", Line: 0, Pos: 55},
				{Type: TokenTypeNumber, Literal: "5", Line: 0, Pos: 57},
				{Type: TokenTypeRightParenthesis, Literal: ")", Line: 0, Pos: 59},
				{Type: TokenTypeEOF, Literal: "", Line: 0, Pos: 60},
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
			if extraToken != (LangToken{}) {
				t.Errorf("NextToken() produced extra token, got = %v", extraToken)
			}
		})
	}
}
