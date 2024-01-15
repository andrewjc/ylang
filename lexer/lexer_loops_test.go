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
				{Type: TokenTypeFor, Literal: "for"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeLet, Literal: "let"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeNumber, Literal: "0"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLessThan, Literal: "<"},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypePlus, Literal: "+"},
				{Type: TokenTypeNumber, Literal: "1"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "Test Lambda For Loop",
			input: "for (i -> i < 5; i = i + 1) { print(i) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLambdaArrow, Literal: "->"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLessThan, Literal: "<"},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypePlus, Literal: "+"},
				{Type: TokenTypeNumber, Literal: "1"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "Test For Each Loop (forEach)",
			input: "for item in range(0, 5) { print(item) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for"},
				{Type: TokenTypeIdentifier, Literal: "item"},
				{Type: TokenTypeIn, Literal: "in"},
				{Type: TokenTypeRange, Literal: "range"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeNumber, Literal: "0"},
				{Type: TokenTypeComma, Literal: ","},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "item"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "For Each Lambda Loop (forEachLambda)",
			input: "for item in range(0, 5) -> (item) { print(item) }",
			want: []LangToken{
				{Type: TokenTypeFor, Literal: "for"},
				{Type: TokenTypeIdentifier, Literal: "item"},
				{Type: TokenTypeIn, Literal: "in"},
				{Type: TokenTypeRange, Literal: "range"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeNumber, Literal: "0"},
				{Type: TokenTypeComma, Literal: ","},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeLambdaArrow, Literal: "->"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "item"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "item"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "Classic While Loop (whileClassic)",
			input: "let i = 0;while (i < 5) { print(i); i = i + 1 }",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeNumber, Literal: "0"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeWhile, Literal: "while"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLessThan, Literal: "<"},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypePlus, Literal: "+"},
				{Type: TokenTypeNumber, Literal: "1"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "Lambda While Loop (whileLambda)",
			input: "let i = 0;while (i -> i < 5) { print(i); i = i + 1 }",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeNumber, Literal: "0"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeWhile, Literal: "while"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLambdaArrow, Literal: "->"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLessThan, Literal: "<"},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypePlus, Literal: "+"},
				{Type: TokenTypeNumber, Literal: "1"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "Classic Do-While Loop (doClassic)",
			input: "let i = 0;do { print(i);i = i + 1 } while (i < 5)",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeNumber, Literal: "0"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeDo, Literal: "do"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypePlus, Literal: "+"},
				{Type: TokenTypeNumber, Literal: "1"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeWhile, Literal: "while"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLessThan, Literal: "<"},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
		{
			name:  "Lambda Do-While Loop (doLambda)",
			input: "let i = 0;do { print(i); i = i + 1 } while (i -> i < 5)",
			want: []LangToken{
				{Type: TokenTypeLet, Literal: "let"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeNumber, Literal: "0"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeDo, Literal: "do"},
				{Type: TokenTypeLeftBrace, Literal: "{"},
				{Type: TokenTypeIdentifier, Literal: "print"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeSemicolon, Literal: ";"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeAssignment, Literal: "="},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypePlus, Literal: "+"},
				{Type: TokenTypeNumber, Literal: "1"},
				{Type: TokenTypeRightBrace, Literal: "}"},
				{Type: TokenTypeWhile, Literal: "while"},
				{Type: TokenTypeLeftParenthesis, Literal: "("},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLambdaArrow, Literal: "->"},
				{Type: TokenTypeIdentifier, Literal: "i"},
				{Type: TokenTypeLessThan, Literal: "<"},
				{Type: TokenTypeNumber, Literal: "5"},
				{Type: TokenTypeRightParenthesis, Literal: ")"},
				{Type: TokenTypeEOF, Literal: ""},
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
