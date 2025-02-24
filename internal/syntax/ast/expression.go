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
)

// Token implements the [Node] interface for Ident and returns the ident token.
func (i IdentExpression) Token() token.Token { return i.Tok }

func (i IdentExpression) expressionNode() {}
