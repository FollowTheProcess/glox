package token_test

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

func TestTokenString(t *testing.T) {
	var kinds []token.Kind
	for i := token.EOF; i <= token.Return; i++ {
		kinds = append(kinds, i)
	}

	for _, kind := range kinds {
		t.Run(kind.String(), func(t *testing.T) {
			start := rand.Int()
			end := rand.Int()
			tok := token.Token{Kind: kind, Start: start, End: end}

			want := fmt.Sprintf("<Token::%s start=%d, end=%d>", kind, start, end)
			test.Equal(t, tok.String(), want)
		})
	}
}

func TestKeyword(t *testing.T) {
	tests := []struct {
		name  string     // Name of the test
		ident string     // Ident to lookup in the set of keywords
		want  token.Kind // Expected return kind
		ok    bool       // Expected return bool
	}{
		{name: "empty", ident: "", want: token.Ident, ok: false},
		{name: "garbage", ident: "$%^&*vhbjj", want: token.Ident, ok: false},
		{name: "non keyword", ident: "some_variable", want: token.Ident, ok: false},
		{name: "if", ident: "if", want: token.If, ok: true},
		{name: "else", ident: "else", want: token.Else, ok: true},
		{name: "or", ident: "or", want: token.Or, ok: true},
		{name: "and", ident: "and", want: token.And, ok: true},
		{name: "for", ident: "for", want: token.For, ok: true},
		{name: "while", ident: "while", want: token.While, ok: true},
		{name: "true", ident: "true", want: token.True, ok: true},
		{name: "false", ident: "false", want: token.False, ok: true},
		{name: "class", ident: "class", want: token.Class, ok: true},
		{name: "super", ident: "super", want: token.Super, ok: true},
		{name: "this", ident: "this", want: token.This, ok: true},
		{name: "fun", ident: "fun", want: token.Fun, ok: true},
		{name: "var", ident: "var", want: token.Var, ok: true},
		{name: "nil", ident: "nil", want: token.Nil, ok: true},
		{name: "print", ident: "print", want: token.Print, ok: true},
		{name: "return", ident: "return", want: token.Return, ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := token.Keyword(tt.ident)

			test.Equal(t, ok, tt.ok)
			test.Equal(t, got, tt.want)
		})
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		kind token.Kind // The token kind under test
		want int        // The expected precedence
	}{
		{kind: token.Or, want: token.PrecedenceOr},
		{kind: token.And, want: token.PrecedenceAnd},
		{kind: token.Eq, want: token.PrecedenceComp},
		{kind: token.BangEq, want: token.PrecedenceComp},
		{kind: token.LessThan, want: token.PrecedenceComp},
		{kind: token.LessThanEq, want: token.PrecedenceComp},
		{kind: token.GreaterThan, want: token.PrecedenceComp},
		{kind: token.GreaterThanEq, want: token.PrecedenceComp},
		{kind: token.Plus, want: token.PrecedenceAddSubtract},
		{kind: token.Minus, want: token.PrecedenceAddSubtract},
		{kind: token.Star, want: token.PrecedenceMulDivide},
		{kind: token.ForwardSlash, want: token.PrecedenceMulDivide},
		{kind: token.Bang, want: token.PrecedenceMin},
		{kind: token.Bang, want: token.PrecedenceMin},
		{kind: token.Number, want: token.PrecedenceMin},
		{kind: token.Ident, want: token.PrecedenceMin},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			tok := token.Token{Kind: tt.kind}

			test.Equal(t, tok.Precedence(), tt.want)
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name string      // Name of the test case
		a, b token.Token // Tokens for comparison
		want bool        // Expected return value
	}{
		{
			name: "equal",
			a:    token.Token{Kind: token.String, Start: 0, End: 3},
			b:    token.Token{Kind: token.String, Start: 0, End: 3},
			want: true,
		},
		{
			name: "different kind",
			a:    token.Token{Kind: token.Number, Start: 0, End: 3},
			b:    token.Token{Kind: token.SemiColon, Start: 0, End: 3},
			want: false,
		},
		{
			name: "different start",
			a:    token.Token{Kind: token.Eq, Start: 1, End: 3},
			b:    token.Token{Kind: token.Eq, Start: 0, End: 3},
			want: false,
		},
		{
			name: "different end",
			a:    token.Token{Kind: token.Eq, Start: 0, End: 4},
			b:    token.Token{Kind: token.Eq, Start: 0, End: 3},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, token.Equal(tt.a, tt.b), tt.want)
		})
	}
}
