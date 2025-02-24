package ast

import "github.com/FollowTheProcess/glox/internal/syntax/token"

// Statement is an AST statement node.
type Statement interface {
	Node            // Marks the Statement as an AST node
	statementNode() // Private method enforcing type safety
}

// Concrete AST Statement node types, all implementing [Node] and [Statement].
type (
	// A VarStatement is the AST node representing a var declaration
	// i.e. `var <ident> = <expression>;`.
	VarStatement struct {
		Value Expression
		Ident IdentExpression
	}

	// A ReturnStatement is the AST node representing a return statement
	// i.e. `return <expression>;`.
	ReturnStatement struct {
		Value Expression
		Tok   token.Token // The "return" keyword
	}
)

// Token implements the [Node] interface for VarDeclaration and returns the first
// token associated with the node, in this case the ident token.
func (v VarStatement) Token() token.Token { return v.Ident.Tok }

// Token implements [Node] for a ReturnStatement and returns the token
// corresponding to the "return" keyword.
func (r ReturnStatement) Token() token.Token { return r.Tok }

func (r ReturnStatement) statementNode() {}
func (v VarStatement) statementNode()    {}
