package lexer

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

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
		{
			name: "Peek single character",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("ba")),
				position: 0,
				ch:       0,
			},
			want: 'a',
		},
		{
			name: "Peek next character",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("ab")),
				position: 1,
				ch:       'a',
			},
			want: 'b',
		},
		{
			name: "Peek at end of string",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("a")),
				position: 1,
				ch:       'a',
			},
			want: 0,
		},
		{
			name: "Peek with multiple characters",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("_hello")),
				position: 0,
				ch:       0,
			},
			want: 'h',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			// Initialize lexer's ch field by reading the first character
			l.readChar()
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
		{
			name: "Peek at index 0",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("test")),
				position: 0,
				ch:       0,
			},
			args: args{index: 0},
			want: 't',
		},
		{
			name: "Peek at index 1",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("test")),
				position: 0,
				ch:       0,
			},
			args: args{index: 1},
			want: 'e',
		},
		{
			name: "Peek at last index",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("test")),
				position: 0,
				ch:       0,
			},
			args: args{index: 3},
			want: 't',
		},
		{
			name: "Peek beyond string length",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("test")),
				position: 0,
				ch:       0,
			},
			args: args{index: 10},
			want: 0,
		},
		{
			name: "Peek with special characters",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("t€st")),
				position: 0,
				ch:       0,
			},
			args: args{index: 1},
			want: '€',
		},
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
