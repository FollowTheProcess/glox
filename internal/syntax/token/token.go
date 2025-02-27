// Package token defines the lexical tokens of the Lox language.
package token

import (
	"fmt"
)

// Kind is the kind of token.
type Kind int

const (
	EOF           Kind = iota // EOF
	Error                     // Lex error e.g. unexpected token
	OpenParen                 // '('
	CloseParen                // ')'
	OpenBrace                 // '{'
	CloseBrace                // '}'
	Comma                     // ','
	Dot                       // '.'
	Minus                     // '-'
	Plus                      // '+'
	SemiColon                 // ';'
	ForwardSlash              // '/'
	Star                      // '*'
	Bang                      // '!'
	Eq                        // '='
	BangEq                    // '!='
	DoubleEq                  // '=='
	GreaterThan               // '>'
	LessThan                  // '<'
	GreaterThanEq             // '>='
	LessThanEq                // '<='
	String                    // String literal, quoted with '"'
	Number                    // Number literal
	Ident                     // Identifier
	If                        // Keyword: 'if'
	Else                      // Keyword: 'else'
	Or                        // Keyword: 'or'
	And                       // Keyword: 'and'
	For                       // Keyword: 'for'
	While                     // Keyword: 'while'
	True                      // Keyword: 'true'
	False                     // Keyword: 'false'
	Class                     // Keyword: 'class'
	Super                     // Keyword: 'super'
	This                      // Keyword: 'this'
	Fun                       // Keyword: 'fun'
	Var                       // Keyword: 'var'
	Nil                       // Keyword: 'nil'
	Print                     // Keyword: 'print'
	Return                    // Keyword: 'return'
)

const (
	PrecedenceMin         = iota // Lowest operator precedence
	PrecedenceEquals             // Precedence of comparson operators like '==', '!=' etc.
	PrecedenceComp               // Less than, greater than, and, or etc.
	PrecedenceAddSubtract        // Precedence of addition '+' and subtraction '-'
	PrecedenceMulDivide          // Precedence of multiplication '*' and division '/'
	PrecedenceUnary              // Unary expressions like '!true'
	PrecedenceMax                // Highest precedence, things like function calls, selectors
)

var tokenStrings = [...]string{
	EOF:           "EOF",
	Error:         "Error",
	OpenParen:     "OpenParen",
	CloseParen:    "CloseParen",
	OpenBrace:     "OpenBrace",
	CloseBrace:    "CloseBrace",
	Comma:         "Comma",
	Dot:           "Dot",
	Minus:         "Minus",
	Plus:          "Plus",
	SemiColon:     "SemiColon",
	ForwardSlash:  "ForwardSlash",
	Star:          "Star",
	Bang:          "Bang",
	Eq:            "Eq",
	BangEq:        "BangEq",
	DoubleEq:      "DoubleEq",
	GreaterThan:   "GreaterThan",
	LessThan:      "LessThan",
	GreaterThanEq: "GreaterThanEq",
	LessThanEq:    "LessThanEq",
	String:        "String",
	Number:        "Number",
	Ident:         "Ident",
	If:            "If",
	Else:          "Else",
	Or:            "Or",
	And:           "And",
	For:           "For",
	While:         "While",
	True:          "True",
	False:         "False",
	Class:         "Class",
	Super:         "Super",
	This:          "This",
	Fun:           "Fun",
	Var:           "Var",
	Nil:           "Nil",
	Print:         "Print",
	Return:        "Return",
}

var tokenLexemes = [...]string{
	EOF:           "EOF",
	Error:         "Error",
	OpenParen:     "(",
	CloseParen:    ")",
	OpenBrace:     "{",
	CloseBrace:    "}",
	Comma:         ",",
	Dot:           ".",
	Minus:         "-",
	Plus:          "+",
	SemiColon:     ";",
	ForwardSlash:  "/",
	Star:          "*",
	Bang:          "!",
	Eq:            "=",
	BangEq:        "!=",
	DoubleEq:      "==",
	GreaterThan:   ">",
	LessThan:      "<",
	GreaterThanEq: ">=",
	LessThanEq:    "<=",
	String:        "String",
	Number:        "Number",
	Ident:         "Ident",
	If:            "if",
	Else:          "else",
	Or:            "or",
	And:           "and",
	For:           "for",
	While:         "while",
	True:          "true",
	False:         "false",
	Class:         "class",
	Super:         "super",
	This:          "this",
	Fun:           "fun",
	Var:           "var",
	Nil:           "nil",
	Print:         "print",
	Return:        "return",
}

// String returns the string representation of [Kind].
func (k Kind) String() string {
	if 0 <= k && k < Kind(len(tokenStrings)) {
		return tokenStrings[k]
	}

	return "<BadToken>"
}

// Lexeme returns the actual text of a [Kind].
//
// For the kinds that are known ahead of time e.g. SemiColon (";"), OpenParen ("(") etc. this
// returns the underlying character, for the likes of Ident, String etc. this returns their name
// like [Kind.String] does.
func (k Kind) Lexeme() string {
	if 0 <= k && k < Kind(len(tokenLexemes)) {
		return tokenLexemes[k]
	}

	return "<BadToken>"
}

// Token is a lexical token in Lox.
//
// It stores the text, type and the offset only. Line and column information
// is calculated when needed by the parser based on this offset.
type Token struct {
	Kind  Kind // The kind of token
	Start int  // The offset in bytes starting from 0 from the start of the input to the first char in this token
	End   int  // The byte offset of the last character in this token
}

// String returns the string representation of [Token].
func (t Token) String() string {
	return fmt.Sprintf("<Token::%s start=%d, end=%d>", t.Kind, t.Start, t.End)
}

// Is reports whether the token is of a given [Kind].
func (t Token) Is(kind Kind) bool {
	return t.Kind == kind
}

// Lexeme returns the lexeme associated with the token's kind.
//
// For the kinds that are known ahead of time e.g. SemiColon (";"), OpenParen ("(") etc. this
// returns the underlying character, for the likes of Ident, String etc. this returns their name
// like [Kind.String] does.
func (t Token) Lexeme() string {
	return t.Kind.Lexeme()
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
	switch t.Kind {
	case DoubleEq, BangEq:
		return PrecedenceEquals
	case LessThan, LessThanEq, GreaterThan, GreaterThanEq, Or, And:
		return PrecedenceComp
	case Plus, Minus:
		return PrecedenceAddSubtract
	case Star, ForwardSlash:
		return PrecedenceMulDivide
	default:
		return PrecedenceMin
	}
}
