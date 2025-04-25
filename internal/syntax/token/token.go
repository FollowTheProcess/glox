// Package token implements the set of lexical tokens in the lox language.
package token

import "fmt"

// Kind represents the type of token.
type Kind int

//go:generate stringer -type Kind -linecomment
const (
	EOF        Kind = iota // EOF
	Error                  // Error
	OpenParen              // OpenParen
	CloseParen             // CloseParen
	OpenBrace              // OpenBrace
	CloseBrace             // CloseBrace
	Comma                  // Comma
	Dot                    // Dot
	Minus                  // Minus
	Plus                   // Plus
	SemiColon              // SemiColon
	Slash                  // Slash
	Star                   // Star
	Bang                   // Bang
	BangEq                 // BangEq
	Eq                     // Eq
	DoubleEq               // DoubleEq
	Greater                // Greater
	GreaterEq              // GreaterEq
	Less                   // Less
	LessEq                 // LessEq
	Ident                  // Ident
	String                 // String
	Number                 // Number
	And                    // And
	Class                  // Class
	Else                   // Else
	False                  // False
	Fun                    // Fun
	For                    // For
	If                     // If
	Nil                    // Nil
	Or                     // Or
	Print                  // Print
	Return                 // Return
	Super                  // Super
	This                   // This
	True                   // True
	Var                    // Var
	While                  // While
)

// Token represents a single lexical token.
type Token struct {
	Kind  Kind // The token kind
	Line  int  // The line number (starting at 1)
	Start int  // Byte offset of the first character of the Token
	End   int  // Byte offset of the last character in the Token (=Start for 1 char tokens)
}

// String implement [fmt.Stringer] for a [Token].
func (t Token) String() string {
	return fmt.Sprintf("<Token::%s line=%d start=%d end=%d>", t.Kind, t.Line, t.Start, t.End)
}
