package lexer_test

import (
	"bytes"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name string        // Name of the test case
		src  string        // Source text to lex, string for convenience
		want []token.Token // Expected tokens
	}{
		{
			name: "empty",
			src:  "",
			want: []token.Token{
				{Kind: token.EOF},
			},
		},
		{
			name: "unexpected",
			src:  "%",
			want: []token.Token{
				{Kind: token.Error, Text: []byte("unexpected char %"), Offset: 0},
				{Kind: token.EOF, Offset: 0},
			},
		},
		{
			name: "open paren",
			src:  "(",
			want: []token.Token{
				{Kind: token.OpenParen, Text: []byte("("), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "close paren",
			src:  ")",
			want: []token.Token{
				{Kind: token.CloseParen, Text: []byte(")"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "open brace",
			src:  "{",
			want: []token.Token{
				{Kind: token.OpenBrace, Text: []byte("{"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "close brace",
			src:  "}",
			want: []token.Token{
				{Kind: token.CloseBrace, Text: []byte("}"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "comma",
			src:  ",",
			want: []token.Token{
				{Kind: token.Comma, Text: []byte(","), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "dot",
			src:  ".",
			want: []token.Token{
				{Kind: token.Dot, Text: []byte("."), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "minus",
			src:  "-",
			want: []token.Token{
				{Kind: token.Minus, Text: []byte("-"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "plus",
			src:  "+",
			want: []token.Token{
				{Kind: token.Plus, Text: []byte("+"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "semicolon",
			src:  ";",
			want: []token.Token{
				{Kind: token.SemiColon, Text: []byte(";"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "forward slash",
			src:  "/",
			want: []token.Token{
				{Kind: token.ForwardSlash, Text: []byte("/"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "star",
			src:  "*",
			want: []token.Token{
				{Kind: token.Star, Text: []byte("*"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "bang",
			src:  "!",
			want: []token.Token{
				{Kind: token.Bang, Text: []byte("!"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "bang equal",
			src:  "!=",
			want: []token.Token{
				{Kind: token.BangEqual, Text: []byte("!="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "equal",
			src:  "=",
			want: []token.Token{
				{Kind: token.Equal, Text: []byte("="), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "double equal",
			src:  "==",
			want: []token.Token{
				{Kind: token.DoubleEqual, Text: []byte("=="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "greater than",
			src:  ">",
			want: []token.Token{
				{Kind: token.GreaterThan, Text: []byte(">"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "greater than equal",
			src:  ">=",
			want: []token.Token{
				{Kind: token.GreaterThanEqual, Text: []byte(">="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "less than",
			src:  "<",
			want: []token.Token{
				{Kind: token.LessThan, Text: []byte("<"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "less than equal",
			src:  "<=",
			want: []token.Token{
				{Kind: token.LessThanEqual, Text: []byte("<="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "comment",
			src:  "// I'm a comment to be ignored",
			want: []token.Token{
				{Kind: token.EOF, Offset: 30},
			},
		},
		{
			name: "ignore whitespace",
			src:  "  \t\t\n\n ()!=",
			want: []token.Token{
				{Kind: token.OpenParen, Text: []byte("("), Offset: 7, Width: 1},
				{Kind: token.CloseParen, Text: []byte(")"), Offset: 8, Width: 1},
				{Kind: token.BangEqual, Text: []byte("!="), Offset: 9, Width: 2},
				{Kind: token.EOF, Offset: 11},
			},
		},
		{
			name: "string",
			src:  `"I'm a string literal"`,
			want: []token.Token{
				{Kind: token.String, Text: []byte(`"I'm a string literal"`), Offset: 0, Width: 22},
				{Kind: token.EOF, Offset: 22},
			},
		},
		{
			name: "unterminated string",
			src:  `"I'm a string literal`,
			want: []token.Token{
				{Kind: token.Error, Text: []byte("unterminated string literal"), Offset: 21},
				{Kind: token.EOF, Offset: 21},
			},
		},
		{
			name: "integer",
			src:  "42",
			want: []token.Token{
				{Kind: token.Number, Text: []byte("42"), Offset: 0},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "float",
			src:  "42.69",
			want: []token.Token{
				{Kind: token.Number, Text: []byte("42.69"), Offset: 0},
				{Kind: token.EOF, Offset: 5},
			},
		},
		{
			name: "ident",
			src:  "some_variable",
			want: []token.Token{
				{Kind: token.Ident, Text: []byte("some_variable"), Offset: 0},
				{Kind: token.EOF, Offset: 13},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := collect(tt.src)
			test.EqualFunc(t, tokens, tt.want, tokenStreamEqual)
		})
	}
}

// collect gathers the emitted tokens into a slice for comparison.
func collect(src string) []token.Token {
	var tokens []token.Token
	l := lexer.New([]byte(src))
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Kind == token.EOF {
			break
		}
	}
	return tokens
}

// tokenStreamEqual compares to slices of tokens for equality.
func tokenStreamEqual(t1, t2 []token.Token) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i := range t1 {
		if t1[i].Kind != t2[i].Kind {
			return false
		}
		if !bytes.Equal(t1[i].Text, t2[i].Text) {
			return false
		}

		if t1[i].Offset != t2[i].Offset {
			return false
		}
	}
	return true
}
