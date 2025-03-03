package ast

import "github.com/FollowTheProcess/glox/internal/syntax/token"

// Declaration is an AST declaration node.
type Declaration interface {
	Node              // Marks the declaration as an AST node
	declarationNode() // Private method enforcing type safety between node types
}

// Concrete AST Declaration node types, all implementing [Node] and [Declaration].
type (
	// A VarDeclaration is the AST node representing a var declaration
	// i.e. `var <ident> = <expression>;`.
	VarDeclaration struct {
		Value Expression
		Ident Ident
	}
)

// [Node] implementations

func (v VarDeclaration) Token() token.Token { return v.Ident.Tok }

// [Declaration] implementations

func (v VarDeclaration) declarationNode() {}
