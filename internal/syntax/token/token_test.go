package token_test

import (
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

func TestToken(t *testing.T) {
	tests := []struct {
		name string
		want string
		tok  token.Token
	}{
		{
			name: "eof",
			tok:  token.Token{Kind: token.EOF},
			want: `<Token::EOF text="", offset=0>`,
		},
		{
			name: "error",
			tok:  token.Token{Kind: token.Error, Text: []byte("bang"), Offset: 42},
			want: `<Token::Error text="bang", offset=42>`,
		},
		{
			name: "open paren",
			tok:  token.Token{Kind: token.OpenParen, Text: []byte("("), Offset: 1},
			want: `<Token::OpenParen text="(", offset=1>`,
		},
		{
			name: "close paren",
			tok:  token.Token{Kind: token.CloseParen, Text: []byte(")"), Offset: 1},
			want: `<Token::CloseParen text=")", offset=1>`,
		},
		{
			name: "open brace",
			tok:  token.Token{Kind: token.OpenBrace, Text: []byte("{"), Offset: 1},
			want: `<Token::OpenBrace text="{", offset=1>`,
		},
		{
			name: "close brace",
			tok:  token.Token{Kind: token.CloseBrace, Text: []byte("}"), Offset: 1},
			want: `<Token::CloseBrace text="}", offset=1>`,
		},
		{
			name: "comma",
			tok:  token.Token{Kind: token.Comma, Text: []byte(","), Offset: 27},
			want: `<Token::Comma text=",", offset=27>`,
		},
		{
			name: "dot",
			tok:  token.Token{Kind: token.Dot, Text: []byte("."), Offset: 2},
			want: `<Token::Dot text=".", offset=2>`,
		},
		{
			name: "minus",
			tok:  token.Token{Kind: token.Minus, Text: []byte("-"), Offset: 32},
			want: `<Token::Minus text="-", offset=32>`,
		},
		{
			name: "plus",
			tok:  token.Token{Kind: token.Plus, Text: []byte("+"), Offset: 185},
			want: `<Token::Plus text="+", offset=185>`,
		},
		{
			name: "semicolon",
			tok:  token.Token{Kind: token.SemiColon, Text: []byte(";"), Offset: 86},
			want: `<Token::SemiColon text=";", offset=86>`,
		},
		{
			name: "forward slash",
			tok:  token.Token{Kind: token.ForwardSlash, Text: []byte("/"), Offset: 17},
			want: `<Token::ForwardSlash text="/", offset=17>`,
		},
		{
			name: "star",
			tok:  token.Token{Kind: token.Star, Text: []byte("*"), Offset: 12},
			want: `<Token::Star text="*", offset=12>`,
		},
		{
			name: "bang",
			tok:  token.Token{Kind: token.Bang, Text: []byte("!"), Offset: 7},
			want: `<Token::Bang text="!", offset=7>`,
		},
		{
			name: "bang equal",
			tok:  token.Token{Kind: token.BangEqual, Text: []byte("!="), Offset: 1},
			want: `<Token::BangEqual text="!=", offset=1>`,
		},
		{
			name: "double equal",
			tok:  token.Token{Kind: token.DoubleEqual, Text: []byte("=="), Offset: 174},
			want: `<Token::DoubleEqual text="==", offset=174>`,
		},
		{
			name: "greater than",
			tok:  token.Token{Kind: token.GreaterThan, Text: []byte(">"), Offset: 22},
			want: `<Token::GreaterThan text=">", offset=22>`,
		},
		{
			name: "less than",
			tok:  token.Token{Kind: token.LessThan, Text: []byte("<"), Offset: 63},
			want: `<Token::LessThan text="<", offset=63>`,
		},
		{
			name: "greater than equal",
			tok:  token.Token{Kind: token.GreaterThanEqual, Text: []byte(">="), Offset: 3},
			want: `<Token::GreaterThanEqual text=">=", offset=3>`,
		},
		{
			name: "less than equal",
			tok:  token.Token{Kind: token.LessThanEqual, Text: []byte("<="), Offset: 7},
			want: `<Token::LessThanEqual text="<=", offset=7>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.tok.String(), tt.want)
			test.True(t, tt.tok.Is(tt.tok.Kind))
		})
	}
}
