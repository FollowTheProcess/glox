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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.tok.String(), tt.want)
			test.True(t, tt.tok.Is(tt.tok.Kind))
		})
	}
}
