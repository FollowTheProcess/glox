package ast

import "github.com/FollowTheProcess/glox/internal/syntax/token"

// Expression is an AST expression node.
type Expression interface {
	Node             // Marks the Expression as an AST node
	expressionNode() // Private method enforcing type safety
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
)

// [Node] implementations

func (i IdentExpression) Token() token.Token { return i.Tok }
func (n NumberLiteral) Token() token.Token   { return n.Tok }
func (u UnaryExpression) Token() token.Token { return u.Tok }

// [Expression] implementations

func (i IdentExpression) expressionNode() {}
func (n NumberLiteral) expressionNode()   {}
func (u UnaryExpression) expressionNode() {}
