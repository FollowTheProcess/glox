// Package parser implements the glox parser.
package parser

import (
	"errors"
	"fmt"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// Parser is the glox parser.
type Parser struct {
	tokeniser    lexer.Tokeniser
	src          []byte
	errs         []error
	currentToken token.Token
	nextToken    token.Token
}

// New returns a new Parser.
func New(src []byte, tokeniser lexer.Tokeniser) *Parser {
	parser := &Parser{
		src:       src,
		tokeniser: tokeniser,
	}

	// Read 2 tokens so current and next are set
	parser.next()
	parser.next()

	return parser
}

// next advances the parser by a single token.
func (p *Parser) next() {
	p.currentToken = p.nextToken
	p.nextToken = p.tokeniser.NextToken()
}

// expect advances the parser by a single token, and asserts that the token
// is of a particular kind. If not, a parsing error will be produced and
// appended to the error list.
func (p *Parser) expect(kind token.Kind) {
	if !p.nextToken.Is(kind) {
		p.errorf("expected %s, got %s: %q", kind, p.nextToken.Kind, string(p.nextToken.Text))
	}

	// Make progress
	p.next()
}

// error emits a plain parse error.
func (p *Parser) error(err error) {
	p.errs = append(p.errs, err)
}

// errorf emits a formatted parse error.
func (p *Parser) errorf(format string, args ...any) {
	p.errs = append(p.errs, fmt.Errorf(format, args...))
}

// Parse is the top level parsing method, and will parse an entire Lox
// source file to completion.
func (p *Parser) Parse() (ast.Program, error) {
	prog := ast.Program{}
	var statement ast.Statement
	for !p.currentToken.Is(token.EOF) {
		switch p.currentToken.Kind {
		case token.Var:
			statement = p.ParseVarDecl()
		default:
			return prog, fmt.Errorf("todo handle %s", p.currentToken.Kind)
		}

		prog.Statements = append(prog.Statements, statement)
	}

	// TODO(@FollowTheProcess): Better error handling, we need a structured error
	// with source locations and stuff so we can hand that off to a handler which
	// shows a nice representation to the user
	return prog, errors.Join(p.errs...)
}

// Errors returns any parsing errors that have been collected during parsing.
func (p *Parser) Errors() []error {
	return p.errs
}

// ParseVarDecl parses a `var <ident> = <expr>` statement.
func (p *Parser) ParseVarDecl() ast.Statement {
	var decl ast.VarDeclaration
	p.expect(token.Ident)
	decl.Ident = ast.Ident{Tok: p.currentToken}

	p.expect(token.Equal)

	// TODO(@FollowTheProcess): Parse the expression, currently just skip
	// until ';' or EOF
	for !p.currentToken.Is(token.SemiColon) {
		p.next()
		if p.currentToken.Is(token.EOF) {
			// TODO(@FollowTheProcess): Sort this out, we need a structured error type
			p.error(ErrUnexpectedEOF)
			return nil
		}
	}

	return decl
}
