package lexer

import (
	"bufio"
	"reflect"
	"testing"
)

func TestLexer_NumericTests(t *testing.T) {
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
			name:  "Test Numbers",
			input: "123 456.789",
			want: []LangToken{
				{Type: TokenTypeNumber, Literal: "123"},
				{Type: TokenTypeNumber, Literal: "456.789"},
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
