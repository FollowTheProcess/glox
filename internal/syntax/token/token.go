// Package token defines the lexical tokens of the Lox language.
package token

import "fmt"

// Kind is the kind of token.
type Kind int

const (
	EOF        Kind = iota // EOF
	OpenParen              // '('
	CloseParen             // ')'
)

// String returns the string representation of [Kind].
func (k Kind) String() string {
	switch k {
	case EOF:
		return "EOF"
	case OpenParen:
		return "OpenParen"
	case CloseParen:
		return "CloseParen"
	default:
		return "<BadToken>"
	}
}

// Token is a lexical token in Lox.
type Token struct {
	Text   []byte // The src text of the token
	Kind   Kind   // The kind of token
	Offset int    // The offset in bytes starting from 0 from the start of the input
}

// String returns the string representation of [Token].
func (t Token) String() string {
	return fmt.Sprintf("<Token::%s text=%q, offset=%d>", t.Kind, t.Text, t.Offset)
}

// Is reports whether the token is of a given [Kind].
func (t Token) Is(kind Kind) bool {
	return t.Kind == kind
}
