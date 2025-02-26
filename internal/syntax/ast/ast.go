// Package ast implements the abstract syntax tree for Lox.
package ast

import (
	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// TODO(@FollowTheProcess): Make all AST nodes pretty print themselves and then tests can just be
// comparing a diff of the produced strings, possibly even in .txtar files with automatic snapshot
// updating

// Node is a generic Abstract Syntax Tree node. All AST nodes implement this interface.
type Node interface {
	// Token returns the first lexical token associated with the node.
	Token() token.Token
}

// A Program is the root AST node, and contains all the statements in the source.
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
