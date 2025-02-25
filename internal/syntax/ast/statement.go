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

	// A PrintStatement is the AST node representing a print statement
	// i.e. `print <expression>;`.
	PrintStatement struct {
		Value Expression
		Tok   token.Token // The "print" keyword
	}
)

// [Node] implementations

func (v VarStatement) Token() token.Token    { return v.Ident.Tok }
func (r ReturnStatement) Token() token.Token { return r.Tok }
func (p PrintStatement) Token() token.Token  { return p.Tok }

// [Statement] implementations

func (r ReturnStatement) statementNode() {}
func (v VarStatement) statementNode()    {}
func (p PrintStatement) statementNode()  {}
