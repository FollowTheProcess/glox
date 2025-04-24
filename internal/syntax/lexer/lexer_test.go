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
