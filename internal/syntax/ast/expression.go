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
	// An IdentExpression is the AST node representing an identifier, both
	// keyword and not.
	IdentExpression struct {
		Name string      // The name of the ident
		Tok  token.Token // The underlying ident token
	}

	// A NumberLiteral is the AST node representing a literal number. Note
	// that numbers in Lox are *all* float64s underneath
	//
	// See https://craftinginterpreters.com/the-lox-language.html#data-types
	NumberLiteral struct {
		Value float64     // The value of the number
		Tok   token.Token // Underlying number token
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
)

// [Node] implementations

func (i IdentExpression) Token() token.Token  { return i.Tok }
func (n NumberLiteral) Token() token.Token    { return n.Tok }
func (u UnaryExpression) Token() token.Token  { return u.Tok }
func (b BinaryExpression) Token() token.Token { return b.Left.Token() }

// [Expression] implementations

func (i IdentExpression) expressionNode()  {}
func (n NumberLiteral) expressionNode()    {}
func (u UnaryExpression) expressionNode()  {}
func (b BinaryExpression) expressionNode() {}

// Precedence implementations

func (i IdentExpression) Precedence() string {
	return i.Name // No precedence here, just the thing
}

func (n NumberLiteral) Precedence() string {
	return strconv.FormatFloat(n.Value, 'g', -1, 64) // Same
}

func (u UnaryExpression) Precedence() string {
	return fmt.Sprintf("(%s%s)", u.Tok.Kind.Lexeme(), u.Value.Precedence())
}

func (b BinaryExpression) Precedence() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.Precedence(), b.Op.Kind.Lexeme(), b.Right.Precedence())
}
