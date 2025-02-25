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

	// Number literal is the AST node representing a literal number. Note
	// that numbers in Lox are *all* float64s underneath
	//
	// See https://craftinginterpreters.com/the-lox-language.html#data-types
	NumberLiteral struct {
		Value float64
		Tok   token.Token
	}
)

// [Node] implementations

func (i IdentExpression) Token() token.Token { return i.Tok }
func (n NumberLiteral) Token() token.Token   { return n.Tok }

// [Expression] implementations

func (i IdentExpression) expressionNode() {}
func (n NumberLiteral) expressionNode()   {}
