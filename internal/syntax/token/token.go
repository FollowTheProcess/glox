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
	case Ident:
		return "Ident"
	case If:
		return "If"
	case Else:
		return "Else"
	case Or:
		return "Or"
	case And:
		return "And"
	case For:
		return "For"
	case While:
		return "While"
	case True:
		return "True"
	case False:
		return "False"
	case Class:
		return "Class"
	case Super:
		return "Super"
	case This:
		return "This"
	case Fun:
		return "Fun"
	case Var:
		return "Var"
	case Nil:
		return "Nil"
	case Print:
		return "Print"
	case Return:
		return "Return"
	default:
		return "<BadToken>"
	}
}

// All keywords in the Lox language, mapped to their token Kind.
var keywords = map[string]Kind{
	"if":     If,
	"else":   Else,
	"or":     Or,
	"and":    And,
	"for":    For,
	"while":  While,
	"true":   True,
	"false":  False,
	"class":  Class,
	"super":  Super,
	"this":   This,
	"fun":    Fun,
	"var":    Var,
	"nil":    Nil,
	"print":  Print,
	"return": Return,
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

// Keyword looks up an identifier in the set of Lox keywords, returning it's
// [Kind] if it was a keyword.
//
// If the ident is a keyword, the keyword Kind is returned and ok = true, if not
// the [Error] Kind is returned with ok = false.
func Keyword(ident string) (kind Kind, ok bool) {
	kw, ok := keywords[ident]
	if !ok {
		return Error, false
	}
	return kw, true
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
