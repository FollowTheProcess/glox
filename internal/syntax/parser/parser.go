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
	name      string       // The name of the source file or "stdin" if REPL
	tokeniser *lexer.Lexer // The lexer
	src       string       // The raw source code
	errs      []error      // List of [SyntaxError], collected while parsing
	current   token.Token  // The current token the parser is sat on
	next      token.Token  // The next token in the stream
}

// New returns a new Parser.
func New(name, src string) *Parser {
	parser := &Parser{
		name:      name,
		src:       src,
		tokeniser: lexer.New(src),
	}

	// Read 2 tokens so current and next are set
	parser.advance()
	parser.advance()

	return parser
}

// advance advances the parser by a single token.
func (p *Parser) advance() {
	p.current = p.next
	p.next = p.tokeniser.NextToken()
}

// expect advances the parser by a single token, and asserts that the token
// is of a particular kind. If not, a parsing error will be produced and
// appended to the error list.
func (p *Parser) expect(kind token.Kind) {
	if !p.next.Is(kind) {
		p.syntaxError(
			"expected %s, got %s: %q",
			kind,
			p.next.Kind,
			p.src[p.next.Start:p.next.End],
		)
	}

	// Make progress, so that p.next above is now p.current
	p.advance()
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

		if index > p.current.Start {
			break
		}
	}

	// The column is therefore the number of bytes from the position of the most recent newline
	// encountered before the token, and the offset of the token itself
	col := p.current.Start - lastNewLineOffset

	err := SyntaxError{
		File:  p.name,
		Msg:   fmt.Sprintf(format, args...),
		Token: p.current,
		Line:  line,
		Col:   col,
	}

	p.errs = append(p.errs, err)
}

// Parse is the top level parsing method, and will parse an entire Lox
// source file to completion.
func (p *Parser) Parse() (ast.Program, error) {
	prog := ast.Program{}
	for !p.current.Is(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			prog.Statements = append(prog.Statements, statement)
		}
		p.advance()
	}

	return prog, errors.Join(p.errs...)
}

// parseStatement parses statements of all kinds.
func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Kind {
	case token.Var:
		return p.parseVarDecl()
	case token.Return:
		return p.parseReturnStatement()
	case token.Print:
		return p.parsePrintStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Errors returns any parsing errors that have been collected during parsing.
func (p *Parser) Errors() []error {
	return p.errs
}

// parseVarDecl parses a `var <ident> = <expr>` statement.
func (p *Parser) parseVarDecl() ast.Statement {
	var statement ast.VarStatement
	p.expect(token.Ident)
	statement.Ident = ast.IdentExpression{Tok: p.current, Name: p.src[p.current.Start:p.current.End]}

	p.expect(token.Eq)

	// TODO(@FollowTheProcess): Parse the expression, currently just skip
	// until ';' or EOF
	for !p.current.Is(token.SemiColon) {
		p.advance()
		if p.current.Is(token.EOF) {
			p.syntaxError("unexpected EOF")
			return nil
		}
	}

	return statement
}

// parseReturnStatement parses a `return <expr>;` statement.
func (p *Parser) parseReturnStatement() ast.Statement {
	statement := ast.ReturnStatement{Tok: p.current}

	p.advance()

	// TODO(@FollowTheProcess): Parse the expression, currently just skip
	// until ';' or EOF
	for !p.current.Is(token.SemiColon) {
		p.advance()
		if p.current.Is(token.EOF) {
			p.syntaxError("unexpected EOF")
			return nil
		}
	}

	return statement
}

// parsePrintStatement parses a `print <expr>;` statement.
func (p *Parser) parsePrintStatement() ast.Statement {
	statement := ast.PrintStatement{Tok: p.current}

	p.advance()

	// TODO(@FollowTheProcess): Parse the expression, currently just skip
	// until ';' or EOF
	for !p.current.Is(token.SemiColon) {
		p.advance()
		if p.current.Is(token.EOF) {
			p.syntaxError("unexpected EOF")
			return nil
		}
	}

	return statement
}

// parseExpressionStatement parses a generic expression statement
// i.e. `<expr>;`.
func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := ast.ExpressionStatement{Tok: p.current}

	// TODO(@FollowTheProcess): Precedence

	statement.Value = p.parseExpression()

	if p.next.Is(token.SemiColon) {
		p.advance()
	}

	return statement
}

// parseExpression is the top level parse function for precedence based
// expression parsing.
func (p *Parser) parseExpression() ast.Expression {
	switch p.current.Kind {
	case token.Ident:
		return p.parseIdentifierExpression()
	default:
		p.syntaxError("TODO: handle %s in parseExpression", p.current.Kind)
		return nil
	}
}

// parseIdentifierExpression parses a single ident expression.
func (p *Parser) parseIdentifierExpression() ast.Expression {
	return ast.IdentExpression{Tok: p.current, Name: p.src[p.current.Start:p.current.End]}
}
