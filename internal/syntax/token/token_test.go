package token_test

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

func TestString(t *testing.T) {
	// All we really care about is the format, let's let quick handle it!
	f := func(tok token.Token) bool {
		return tok.String() == fmt.Sprintf(
			"<Token::%s line=%d start=%d end=%d>",
			tok.Kind.String(),
			tok.Line,
			tok.Start,
			tok.End,
		)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}
}
