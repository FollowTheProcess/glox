package lexer_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax"
	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

func TestBasics(t *testing.T) {
	tests := []struct {
		name string        // Name of the test case
		src  string        // Source code to scan
		want []token.Token // Expected tokens
	}{
		{
			name: "OpenParen",
			src:  "(",
			want: []token.Token{
				{Kind: token.OpenParen, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "CloseParen",
			src:  ")",
			want: []token.Token{
				{Kind: token.CloseParen, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "OpenBrace",
			src:  "{",
			want: []token.Token{
				{Kind: token.OpenBrace, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "CloseBrace",
			src:  "}",
			want: []token.Token{
				{Kind: token.CloseBrace, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Comma",
			src:  ",",
			want: []token.Token{
				{Kind: token.Comma, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Dot",
			src:  ".",
			want: []token.Token{
				{Kind: token.Dot, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Minus",
			src:  "-",
			want: []token.Token{
				{Kind: token.Minus, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Plus",
			src:  "+",
			want: []token.Token{
				{Kind: token.Plus, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "SemiColon",
			src:  ";",
			want: []token.Token{
				{Kind: token.SemiColon, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Star",
			src:  "*",
			want: []token.Token{
				{Kind: token.Star, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Bang",
			src:  "!",
			want: []token.Token{
				{Kind: token.Bang, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Slash",
			src:  "/",
			want: []token.Token{
				{Kind: token.Slash, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "Comment",
			src:  "// A comment here",
			want: []token.Token{},
		},
		{
			name: "Comment newline",
			src:  "// A comment here\n",
			want: []token.Token{},
		},
		{
			name: "Eq",
			src:  "=",
			want: []token.Token{
				{Kind: token.Eq, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "BangEq",
			src:  "!=",
			want: []token.Token{
				{Kind: token.BangEq, Line: 1, Start: 0, End: 2},
			},
		},
		{
			name: "DoubleEq",
			src:  "==",
			want: []token.Token{
				{Kind: token.DoubleEq, Line: 1, Start: 0, End: 2},
			},
		},
		{
			name: "Greater",
			src:  ">",
			want: []token.Token{
				{Kind: token.Greater, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "GreaterEq",
			src:  ">=",
			want: []token.Token{
				{Kind: token.GreaterEq, Line: 1, Start: 0, End: 2},
			},
		},
		{
			name: "Less",
			src:  "<",
			want: []token.Token{
				{Kind: token.Less, Line: 1, Start: 0, End: 1},
			},
		},
		{
			name: "LessEq",
			src:  "<=",
			want: []token.Token{
				{Kind: token.LessEq, Line: 1, Start: 0, End: 2},
			},
		},
		{
			name: "String",
			src:  `"I'm a string literal"`,
			want: []token.Token{
				{Kind: token.String, Line: 1, Start: 0, End: 22},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.New(tt.name, tt.src, testFailHandler(t))

			var got []token.Token
			for tok := lex.NextToken(); tok.Kind != token.EOF; tok = lex.NextToken() {
				got = append(got, tok)
			}

			test.EqualFunc(t, got, tt.want, slices.Equal)
		})
	}
}

// testFailHandler returns a [syntax.ErrorHandler] that handles lexing errors by failing
// the enclosing test.
func testFailHandler(tb testing.TB) syntax.ErrorHandler {
	tb.Helper()
	return func(pos syntax.Position, msg string) {
		tb.Fatalf("%s: %s\n", pos, msg)
	}
}
