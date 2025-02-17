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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.tok.String(), tt.want)
			test.True(t, tt.tok.Is(tt.tok.Kind))
		})
	}
}
