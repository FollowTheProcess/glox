// Package token defines the lexical tokens of the Lox language.
package token

import (
	"fmt"
)

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
	Ident                        // Identifier
	If                           // Keyword: 'if'
	Else                         // Keyword: 'else'
	Or                           // Keyword: 'or'
	And                          // Keyword: 'and'
	For                          // Keyword: 'for'
	While                        // Keyword: 'while'
	True                         // Keyword: 'true'
	False                        // Keyword: 'false'
	Class                        // Keyword: 'class'
	Super                        // Keyword: 'super'
	This                         // Keyword: 'this'
	Fun                          // Keyword: 'fun'
	Var                          // Keyword: 'var'
	Nil                          // Keyword: 'nil'
	Print                        // Keyword: 'print'
	Return                       // Keyword: 'return'
)

const (
	PrecedenceMin         = 0 // Lowest operator precedence
	PrecedenceOr          = 1 // Precedence of the 'or' binary operator
	PrecedenceAnd         = 2 // Precedence of the 'and' binary operator
	PrecedenceComp        = 3 // Precedence of comparson operators like '==', '!=' etc.
	PrecedenceAddSubtract = 4 // Precedence of addition '+' and subtraction '-'
	PrecedenceMulDivide   = 5 // Precedence of multiplication '*' and division '/'
)

var tokenStrings = [...]string{
	EOF:              "EOF",
	Error:            "Error",
	OpenParen:        "OpenParen",
	CloseParen:       "CloseParen",
	OpenBrace:        "OpenBrace",
	CloseBrace:       "CloseBrace",
	Comma:            "Comma",
	Dot:              "Dot",
	Minus:            "Minus",
	Plus:             "Plus",
	SemiColon:        "SemiColon",
	ForwardSlash:     "ForwardSlash",
	Star:             "Star",
	Bang:             "Bang",
	Equal:            "Equal",
	BangEqual:        "BangEqual",
	DoubleEqual:      "DoubleEqual",
	GreaterThan:      "GreaterThan",
	LessThan:         "LessThan",
	GreaterThanEqual: "GreaterThanEqual",
	LessThanEqual:    "LessThanEqual",
	String:           "String",
	Number:           "Number",
	Ident:            "Ident",
	If:               "If",
	Else:             "Else",
	Or:               "Or",
	And:              "And",
	For:              "For",
	While:            "While",
	True:             "True",
	False:            "False",
	Class:            "Class",
	Super:            "Super",
	This:             "This",
	Fun:              "Fun",
	Var:              "Var",
	Nil:              "Nil",
	Print:            "Print",
	Return:           "Return",
}

// String returns the string representation of [Kind].
func (k Kind) String() string {
	if 0 <= k && k < Kind(len(tokenStrings)) {
		return tokenStrings[k]
	}

	return "<BadToken>"
}

// Token is a lexical token in Lox.
//
// It stores the text, type and the offset only. Line and column information
// is calculated when needed by the parser based on this offset.
type Token struct {
	Text   []byte // The src text of the token
	Kind   Kind   // The kind of token
	Offset int    // The offset in bytes starting from 0 from the start of the input to the start of this token
	Width  int    // The width in bytes of this token's raw src text
}

// String returns the string representation of [Token].
func (t Token) String() string {
	if t.Kind == String {
		// Don't double quote the value if it's a string
		return fmt.Sprintf("<Token::%s text=%s, offset=%d, width=%d>", t.Kind, t.Text, t.Offset, t.Width)
	}
	return fmt.Sprintf("<Token::%s text=%q, offset=%d, width=%d>", t.Kind, t.Text, t.Offset, t.Width)
}

// Is reports whether the token is of a given [Kind].
func (t Token) Is(kind Kind) bool {
	return t.Kind == kind
}

// Keyword looks up an identifier in the set of Lox keywords, returning it's
// [Kind] if it was a keyword.
//
// If the ident is a keyword, the keyword Kind is returned and ok = true, if not
// the [Ident] Kind is returned with ok = false.
func Keyword(ident string) (kind Kind, ok bool) {
	switch ident {
	case "if":
		return If, true
	case "else":
		return Else, true
	case "or":
		return Or, true
	case "and":
		return And, true
	case "for":
		return For, true
	case "while":
		return While, true
	case "true":
		return True, true
	case "false":
		return False, true
	case "class":
		return Class, true
	case "super":
		return Super, true
	case "this":
		return This, true
	case "fun":
		return Fun, true
	case "var":
		return Var, true
	case "nil":
		return Nil, true
	case "print":
		return Print, true
	case "return":
		return Return, true
	default:
		return Ident, false
	}
}

// Precedence returns the precedence of the binary expression operator token.
//
// If the token is not a binary operator, the lowest precedence is returned.
func (t Token) Precedence() int {
	// TODO(@FollowTheProcess): What about a binding power thing like
	// https://matklad.github.io/2020/04/13/simple-but-powerful-pratt-parsing.html
	// we could return left and right binding power?
	switch t.Kind {
	case Or:
		return PrecedenceOr
	case And:
		return PrecedenceAnd
	case Equal, BangEqual, LessThan, LessThanEqual, GreaterThan, GreaterThanEqual:
		return PrecedenceComp
	case Plus, Minus:
		return PrecedenceAddSubtract
	case Star, ForwardSlash:
		return PrecedenceMulDivide
	default:
		return PrecedenceMin
	}
}
