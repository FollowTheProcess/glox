// Package token defines the lexical tokens of the Lox language.
package token

import "fmt"

// Kind is the kind of token.
type Kind int

const (
	EOF              Kind = iota // EOF
	Error                        // Lex error e.g. unexpected token
	OpenParen                    // '('
	CloseParen                   // ')'
	OpenBrace                    // '{'
	CloseBrace                   // '}'
	Comma                        // ','
	Dot                          // '.'
	Minus                        // '-'
	Plus                         // '+'
	SemiColon                    // ';'
	ForwardSlash                 // '/'
	Star                         // '*'
	Bang                         // '!'
	Equal                        // '='
	BangEqual                    // '!='
	DoubleEqual                  // '=='
	GreaterThan                  // '>'
	LessThan                     // '<'
	GreaterThanEqual             // '>='
	LessThanEqual                // '<='
	String                       // String literal, quoted with '"'
	Number                       // Number literal
)

// String returns the string representation of [Kind].
func (k Kind) String() string { //nolint: cyclop // This is technically high but obviously trivial
	switch k {
	case EOF:
		return "EOF"
	case Error:
		return "Error"
	case OpenParen:
		return "OpenParen"
	case CloseParen:
		return "CloseParen"
	case OpenBrace:
		return "OpenBrace"
	case CloseBrace:
		return "CloseBrace"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case SemiColon:
		return "SemiColon"
	case ForwardSlash:
		return "ForwardSlash"
	case Star:
		return "Star"
	case Bang:
		return "Bang"
	case Equal:
		return "Equal"
	case BangEqual:
		return "BangEqual"
	case DoubleEqual:
		return "DoubleEqual"
	case GreaterThan:
		return "GreaterThan"
	case LessThan:
		return "LessThan"
	case GreaterThanEqual:
		return "GreaterThanEqual"
	case LessThanEqual:
		return "LessThanEqual"
	case String:
		return "String"
	case Number:
		return "Number"
	default:
		return "<BadToken>"
	}
}

// Token is a lexical token in Lox.
type Token struct {
	Text   []byte // The src text of the token
	Kind   Kind   // The kind of token
	Offset int    // The offset in bytes starting from 0 from the start of the input to the start of this token
	Width  int    // The width in bytes of this token's raw src text
}

// String returns the string representation of [Token].
func (t Token) String() string {
	return fmt.Sprintf("<Token::%s text=%q, offset=%d, width=%d>", t.Kind, t.Text, t.Offset, t.Width)
}

// Is reports whether the token is of a given [Kind].
func (t Token) Is(kind Kind) bool {
	return t.Kind == kind
}
