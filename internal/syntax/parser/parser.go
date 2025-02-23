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
	name         string       // The name of the source file or "stdin" if REPL
	tokeniser    *lexer.Lexer // The lexer
	src          string       // The raw source code
	errs         []error      // List of [SyntaxError], collected while parsing
	currentToken token.Token  // The current token the parser is sat on
	nextToken    token.Token  // The next token in the stream
}

// New returns a new Parser.
func New(name, src string) *Parser {
	parser := &Parser{
		name:      name,
		src:       src,
		tokeniser: lexer.New(src),
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
		p.syntaxError(
			"expected %s, got %s: %q",
			kind,
			p.nextToken.Kind,
			p.src[p.nextToken.Start:p.nextToken.End],
		)
	}

	// Make progress
	p.next()
}

// syntaxError emits a [SyntaxError], populating it with line/col info using
// the parser's current state.
func (p *Parser) syntaxError(format string, args ...any) {
	// Calculate line and col based on the offending token's offset
	line := 1              // Line counter
	lastNewLineOffset := 0 // The byte offset of the last newline seen
	for index, byt := range p.src {
		if byt == '\n' {
			line++
			lastNewLineOffset = index
		}

		if index > p.currentToken.Start {
			break
		}
	}

	// The column is therefore the number of bytes from the position of the most recent newline
	// encountered before the token, and the offset of the token itself
	col := p.currentToken.Start - lastNewLineOffset

	err := SyntaxError{
		File:  p.name,
		Msg:   fmt.Sprintf(format, args...),
		Token: p.currentToken,
		Line:  line,
		Col:   col,
	}

	p.errs = append(p.errs, err)
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
			p.syntaxError("TODO: handle %s", p.currentToken.Kind)
		}

		prog.Statements = append(prog.Statements, statement)
	}

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
	decl.Ident = ast.Ident{Tok: p.currentToken, Name: p.src[p.currentToken.Start:p.currentToken.End]}

	p.expect(token.Eq)

	// TODO(@FollowTheProcess): Parse the expression, currently just skip
	// until ';' or EOF
	for !p.currentToken.Is(token.SemiColon) {
		p.next()
		if p.currentToken.Is(token.EOF) {
			p.syntaxError("unexpected EOF")
			return nil
		}
	}

	return decl
}
