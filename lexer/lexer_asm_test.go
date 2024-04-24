package lexer

import (
	"bufio"
	"strings"
	"testing"
)

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
		{
			name: "Read simple assembly command",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("mov eax, ebx")),
				position: 0,
				ch:       0,
			},
			want: "mov",
		},
		{
			name: "Read assembly with numbers",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("add eax, 10")),
				position: 0,
				ch:       0,
			},
			want: "add",
		},
		{
			name: "Read assembly with special characters",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("cmp eax, [ebx+esi*4]")),
				position: 0,
				ch:       0,
			},
			want: "cmp",
		},
		{
			name: "Read assembly with whitespace",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("  jmp label")),
				position: 0,
				ch:       0,
			},
			want: "jmp",
		},
		{
			name: "Read multi-word assembly",
			fields: fields{
				reader:   bufio.NewReader(strings.NewReader("push dword ptr [ebp-4]")),
				position: 0,
				ch:       0,
			},
			want: "push",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				reader:   tt.fields.reader,
				Position: tt.fields.position,
				ch:       tt.fields.ch,
			}
			l.ReadChar() // Initialize the lexer's ch field
			if got := l.readAssembly(); got != tt.want {
				t.Errorf("readAssembly() = %v, want %v", got, tt.want)
			}
		})
	}
}
