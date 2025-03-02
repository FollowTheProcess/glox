package ast

import (
	"fmt"
	"strconv"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// Expression is an AST expression node.
type Expression interface {
	Node             // Marks the Expression as an AST node
	expressionNode() // Private method enforcing type safety

	// Precedence returns a string clearly showing the precedence hierarchy of the expression.
	// e.g. for the expression `a + b + b;` something like `((a + b) + c);`.
	Precedence() string
}

// Concrete AST Expression node types, all implementing [Node] and [Expression].
type (
	// An Ident is the AST node representing an identifier, both
	// keyword and not.
	Ident struct {
		Name string      // The name of the ident
		Tok  token.Token // The underlying ident token
	}

	// A Number is the AST node representing a literal number. Note
	// that numbers in Lox are *all* float64s underneath
	//
	// See https://craftinginterpreters.com/the-lox-language.html#data-types
	Number struct {
		Value float64     // The value of the number
		Tok   token.Token // Underlying number token
	}

	// A Bool is the AST node representing a literal true|false.
	Bool struct {
		Value bool        // The value of the boolean
		Tok   token.Token // Underlying token.True|token.False
	}

	// A String is the AST node representing a string literal.
	String struct {
		Value string      // The underlying string literal
		Tok   token.Token // The string literal token
	}

	// A UnaryExpression is the AST node representing a unary expression
	// i.e. `-5` or `!true`.
	UnaryExpression struct {
		Value Expression  // The expression to be "unary'd"
		Tok   token.Token // The operator e.g. `-` or `!`
	}

	// A BinaryExpression is the AST node representing a binary expression
	// i.e. `x != y` or `5 + 5`.
	BinaryExpression struct {
		Left  Expression  // The lhs of the expression
		Right Expression  // The rhs of the expression
		Op    token.Token // The operator token
	}

	// GroupedExpression is the AST node representing a grouped (parenthesised)
	// expression i.e. `(x + y);`.
	GroupedExpression struct {
		Value  Expression  // The inner expression
		LParen token.Token // Opening '(' token
		RParen token.Token // Closing ')' token
	}
)

// [Node] implementations

func (i Ident) Token() token.Token             { return i.Tok }
func (n Number) Token() token.Token            { return n.Tok }
func (b Bool) Token() token.Token              { return b.Tok }
func (s String) Token() token.Token            { return s.Tok }
func (u UnaryExpression) Token() token.Token   { return u.Tok }
func (b BinaryExpression) Token() token.Token  { return b.Left.Token() }
func (g GroupedExpression) Token() token.Token { return g.LParen }

// [Expression] implementations

func (i Ident) expressionNode()             {}
func (n Number) expressionNode()            {}
func (b Bool) expressionNode()              {}
func (s String) expressionNode()            {}
func (u UnaryExpression) expressionNode()   {}
func (b BinaryExpression) expressionNode()  {}
func (g GroupedExpression) expressionNode() {}

// Precedence implementations

func (i Ident) Precedence() string {
	return i.Name // No precedence here, just the thing
}

func (n Number) Precedence() string {
	return strconv.FormatFloat(n.Value, 'g', -1, 64) // Same
}

func (b Bool) Precedence() string {
	return strconv.FormatBool(b.Value)
}

func (s String) Precedence() string {
	return s.Value
}

func (u UnaryExpression) Precedence() string {
	return fmt.Sprintf("(%s%s)", u.Tok.Kind.Lexeme(), u.Value.Precedence())
}

func (b BinaryExpression) Precedence() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.Precedence(), b.Op.Kind.Lexeme(), b.Right.Precedence())
}

func (g GroupedExpression) Precedence() string {
	return g.Value.Precedence()
}
