// Package parser implements the glox parser.
package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// TODO(@FollowTheProcess): A sort of integration test where there is a single txtar archive file per test case
// containing 3 files:
// 	- src.lox: 		Raw Lox source code
//	- tokens.txt:	The output from tokenising src.lox
//  - ast.txt:		Some serialised representation of the AST from parsing tokens.txt
//
// Would thoroughly test all 3 components without coupling them together too much

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

// expect asserts that the next token is of a particular kind, appending a syntax
// error to the parser if it's not.
//
// If the next token is as expected, expect advances the parser onto that token so
// that it is now p.current.
func (p *Parser) expect(kind token.Kind) {
	if !p.next.Is(kind) {
		if p.next.Is(token.EOF) {
			// If it's EOF, don't show the empty string
			p.syntaxError(
				"expected %q, got %s",
				kind.Lexeme(),
				p.next.Kind,
			)
		} else {
			p.syntaxError(
				"expected %q, got %s: %q",
				kind.Lexeme(),
				p.next.Kind,
				p.src[p.next.Start:p.next.End],
			)
		}
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
		if statement != nil && len(p.errs) == 0 {
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
	statement.Ident = ast.Ident{Tok: p.current, Name: p.src[p.current.Start:p.current.End]}

	p.expect(token.Eq)
	p.advance()

	statement.Value = p.parseExpression(token.PrecedenceMin)

	p.expect(token.SemiColon)

	return statement
}

// parseReturnStatement parses a `return <expr>;` statement.
func (p *Parser) parseReturnStatement() ast.Statement {
	statement := ast.ReturnStatement{Tok: p.current}

	p.advance()
	statement.Value = p.parseExpression(token.PrecedenceMin)

	p.expect(token.SemiColon)

	return statement
}

// parsePrintStatement parses a `print <expr>;` statement.
func (p *Parser) parsePrintStatement() ast.Statement {
	statement := ast.PrintStatement{Tok: p.current}

	p.advance()
	statement.Value = p.parseExpression(token.PrecedenceMin)

	p.expect(token.SemiColon)

	return statement
}

// parseExpressionStatement parses a generic expression statement
// i.e. `<expr>;`.
func (p *Parser) parseExpressionStatement() ast.Statement {
	statement := ast.ExpressionStatement{Tok: p.current}
	statement.Value = p.parseExpression(token.PrecedenceMin)

	if p.next.Is(token.SemiColon) {
		p.advance()
	}

	return statement
}

// parseExpression is the top level parse function for precedence based
// expression parsing.
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var expression ast.Expression

	switch p.current.Kind {
	case token.Ident:
		expression = p.parseIdentifierExpression()
	case token.Number:
		expression = p.parseNumberLiteralExpression()
	case token.Bang, token.Minus:
		expression = p.parseUnaryExpression()
	case token.True, token.False:
		expression = p.parseBoolLiteralExpression()
	case token.OpenParen:
		expression = p.parseGroupedExpression()
	default:
		p.syntaxError("Unhandled token in parseExpression (unary switch): %s", p.current.Kind)
		return nil
	}

	for !p.next.Is(token.SemiColon) && precedence < p.next.Precedence() {
		p.advance()
		switch p.current.Kind {
		case token.Or,
			token.And,
			token.DoubleEq,
			token.BangEq,
			token.LessThan,
			token.LessThanEq,
			token.GreaterThan,
			token.GreaterThanEq,
			token.Plus,
			token.Minus,
			token.Star,
			token.ForwardSlash:
			expression = p.parseBinaryExpression(expression)
		default:
			p.syntaxError("Unhandled token in parseExpression (binary switch): %s", p.current.Kind)
			return nil
		}
	}

	return expression
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return ast.Ident{Tok: p.current, Name: p.src[p.current.Start:p.current.End]}
}

func (p *Parser) parseNumberLiteralExpression() ast.Expression {
	src := p.src[p.current.Start:p.current.End]
	n, err := strconv.ParseFloat(src, 64)
	if err != nil {
		p.syntaxError("invalid number literal %q: %v", src, err)
		return nil
	}

	return ast.Number{Value: n, Tok: p.current}
}

func (p *Parser) parseBoolLiteralExpression() ast.Expression {
	src := p.src[p.current.Start:p.current.End]
	v, err := strconv.ParseBool(src)
	if err != nil {
		p.syntaxError("invalid boolean literal %q: %v", src, err)
	}

	return ast.Bool{Value: v, Tok: p.current}
}

// parseUnaryExpression parses a unary expression
// i.e. `!true`.
func (p *Parser) parseUnaryExpression() ast.Expression {
	expression := ast.UnaryExpression{Tok: p.current}

	p.advance()

	expression.Value = p.parseExpression(token.PrecedenceUnary)

	return expression
}

// parseBinaryExpression parses a binary expression
// i.e. `x != y` or `5 + 5`.
func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expression := ast.BinaryExpression{Left: left, Op: p.current}

	precedence := p.current.Precedence()
	p.advance()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseGroupedExpression parses a parenthesised expression
// i.e. `(x + y)`.
func (p *Parser) parseGroupedExpression() ast.Expression {
	expression := ast.GroupedExpression{LParen: p.current}
	p.advance()
	inner := p.parseExpression(token.PrecedenceMin)

	p.expect(token.CloseParen)

	expression.Value = inner
	expression.RParen = p.current

	return expression
}
