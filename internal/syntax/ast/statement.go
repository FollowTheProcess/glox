package ast

import "github.com/FollowTheProcess/glox/internal/syntax/token"

// Statement is an AST statement node.
type Statement interface {
	Node            // Marks the Statement as an AST node
	statementNode() // Private method enforcing type safety
}

// TODO(@FollowTheProcess): VarStatement will actually get turned into VarDeclaration
// so we'll need a new file and node type

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
		Tok   token.Token // The first token of the expression
	}
)

// [Node] implementations

func (v VarStatement) Token() token.Token        { return v.Ident.Tok }
func (r ReturnStatement) Token() token.Token     { return r.Tok }
func (p PrintStatement) Token() token.Token      { return p.Tok }
func (e ExpressionStatement) Token() token.Token { return e.Tok }

// [Statement] implementations

func (r ReturnStatement) statementNode()     {}
func (v VarStatement) statementNode()        {}
func (p PrintStatement) statementNode()      {}
func (e ExpressionStatement) statementNode() {}
