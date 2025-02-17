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
				{Kind: token.OpenParen, Text: []byte("("), Offset: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "close paren",
			src:  ")",
			want: []token.Token{
				{Kind: token.CloseParen, Text: []byte(")"), Offset: 1},
				{Kind: token.EOF, Offset: 1},
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
