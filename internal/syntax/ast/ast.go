// Package ast implements the abstract syntax tree for Lox.
package ast

import "github.com/FollowTheProcess/glox/internal/syntax/token"

// Node is a generic Abstract Syntax Tree node.
//
// All AST nodes implement this interface.
type Node interface {
	// Token returns the first lexical token associated with the node.
	Token() token.Token
}

// Program represents an entire blob of Lox source in AST form.
type Program struct {
	Statements []Statement
}

// Token implements [Node] for a [Program] and returns the first token associated
// with the node, in this case the very first token of the entire program.
//
// If the program is empty, the returned token will be [token.EOF].
func (p Program) Token() token.Token {
	if len(p.Statements) == 0 {
		return token.Token{Kind: token.EOF}
	}

	return p.Statements[0].Token()
}

// Statement is an AST statement node.
type Statement interface {
	Node            // Marks the Statement as an AST node
	statementNode() // Private method enforcing type safety
}

// Expression is an AST expression node.
type Expression interface {
	Node             // Marks the Expression as an AST node
	expressionNode() // Private method enforcing type safety
}

// VarDeclaration is the ast node representing a var <ident> = <expression>; statement.
type VarDeclaration struct {
	Value Expression
	Ident Ident
}

// Token implements the [Node] interface for VarDeclaration and returns the first
// token associated with the node, in this case the ident token.
func (v VarDeclaration) Token() token.Token {
	return v.Ident.Token()
}

func (v VarDeclaration) statementNode() {}

// Ident is the ast node representing an identifier.
type Ident struct {
	Tok token.Token // The underlying ident token
}

// Token implements the [Node] interface for Ident and returns the ident token.
func (i Ident) Token() token.Token {
	return i.Tok
}

// Name returns the literal name of the ident.
func (i Ident) Name() string {
	// TODO(@FollowTheProcess): How do we get the name now tokens don't carry
	// text
	return "TODO"
}
