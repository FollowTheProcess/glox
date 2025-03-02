package ast

import (
	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// Statement is an AST statement node.
type Statement interface {
	Node            // Marks the Statement as an AST node
	statementNode() // Private method enforcing type safety
}

// Concrete AST Statement node types, all implementing [Node] and [Statement].
type (
	// A ReturnStatement is the AST node representing a return statement
	// i.e. `return <expression>;`.
	ReturnStatement struct {
		Value Expression
		Tok   token.Token // The "return" keyword
	}

	// A PrintStatement is the AST node representing a print statement
	// i.e. `print <expression>;`.
	PrintStatement struct {
		Value Expression
		Tok   token.Token // The "print" keyword
	}

	// An ExpressionStatement is the AST node representing a top level
	// expression as a statement.
	ExpressionStatement struct {
		Value Expression
	}

	// A DeclarationStatement is the AST node representing a top level
	// declaration as a statement.
	DeclarationStatement struct {
		Value Declaration
	}
)

// [Node] implementations

func (r ReturnStatement) Token() token.Token      { return r.Tok }
func (p PrintStatement) Token() token.Token       { return p.Tok }
func (e ExpressionStatement) Token() token.Token  { return e.Value.Token() }
func (d DeclarationStatement) Token() token.Token { return d.Value.Token() }

// [Statement] implementations

func (r ReturnStatement) statementNode()      {}
func (p PrintStatement) statementNode()       {}
func (e ExpressionStatement) statementNode()  {}
func (d DeclarationStatement) statementNode() {}
