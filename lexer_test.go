package compiler

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
			name:  "Test TokenType",
			input: "function let if else for while do switch case default data type asm",
			want: []LangToken{
				{Type: TokenTypeFunction, Literal: "function"},
				{Type: TokenTypeLet, Literal: "let"},
				{Type: TokenTypeIf, Literal: "if"},
				{Type: TokenTypeElse, Literal: "else"},
				{Type: TokenTypeFor, Literal: "for"},
				{Type: TokenTypeWhile, Literal: "while"},
				{Type: TokenTypeDo, Literal: "do"},
				{Type: TokenTypeSwitch, Literal: "switch"},
				{Type: TokenTypeCase, Literal: "case"},
				{Type: TokenTypeDefault, Literal: "default"},
				{Type: TokenTypeData, Literal: "data"},
				{Type: TokenTypeType, Literal: "type"},
				{Type: TokenTypeAssembly, Literal: "asm"},
				{Type: TokenTypeEOF, Literal: ""},
			},
			wantErr: false,
		},
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
		{
			name:  "Test Strings",
			input: `"Hello, world!" "Another string"`,
			want: []LangToken{
				{Type: TokenTypeString, Literal: "Hello, world!"},
				{Type: TokenTypeString, Literal: "Another string"},
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

func TestLexer_peekChar(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name   string
		fields fields
		want   rune
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			if got := l.peekChar(); got != tt.want {
				t.Errorf("peekChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_peekCharAtIndex(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   rune
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			if got := l.peekCharAtIndex(tt.args.index); got != tt.want {
				t.Errorf("peekCharAtIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_readAssembly(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			if got := l.readAssembly(); got != tt.want {
				t.Errorf("readAssembly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_readChar(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			l.readChar()
		})
	}
}

func TestLexer_readIdentifier(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name   string
		fields fields
		want   TokenType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			if got, _ := l.readIdentifier(); got != tt.want {
				t.Errorf("readIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_readNumber(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			if got := l.readNumber(); got != tt.want {
				t.Errorf("readNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_readString(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			if got := l.readString(); got != tt.want {
				t.Errorf("readString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLexer_skipWhitespace(t *testing.T) {
	type fields struct {
		reader   *bufio.Reader
		position int
		ch       rune
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			l.skipWhitespace()
		})
	}
}

func TestLookupIdent(t *testing.T) {
	type args struct {
		ident string
	}
	tests := []struct {
		name string
		args args
		want TokenType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LookupIdent(tt.args.ident); got != tt.want {
				t.Errorf("LookupIdent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLexer(t *testing.T) {
	type args struct {
		inputFile string
	}
	tests := []struct {
		name    string
		args    args
		want    *Lexer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLexer(tt.args.inputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLexer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLexer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLexerFromString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    *Lexer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLexerFromString(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLexerFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLexerFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDigit(t *testing.T) {
	type args struct {
		ch rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDigit(tt.args.ch); got != tt.want {
				t.Errorf("isDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLetter(t *testing.T) {
	type args struct {
		ch rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLetter(tt.args.ch); got != tt.want {
				t.Errorf("isLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newToken(t *testing.T) {
	type args struct {
		tokenType TokenType
		ch        rune
	}
	tests := []struct {
		name string
		args args
		want LangToken
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newToken(tt.args.tokenType, tt.args.ch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
