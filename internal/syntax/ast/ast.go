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

// BinaryExpression represents a binary expression e.g. x != y.
type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator token.Token
}

// Token implements [Node] for [BinaryExpression] and returns the first
// token associated with the node, in this case the token of the left [Expression].
func (b BinaryExpression) Token() token.Token {
	return b.Left.Token()
}

// expressionNode marks a [BinaryExpression] as an [Expression].
func (b BinaryExpression) expressionNode() {}
