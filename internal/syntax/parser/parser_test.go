package parser_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

// testLexer is a lexer that implements the Tokeniser interface simply
// by returning tokens from a list.
//
// It allows us to decouple the parser from the implementation of the lexer.
type testLexer struct {
	tokens []token.Token
}

// NextToken implements the Tokeniser interface for our testLexer.
func (t *testLexer) NextToken() token.Token {
	if len(t.tokens) == 0 {
		return token.Token{Kind: token.EOF}
	}

	// Get the first token in the stream, "consume" it and return it
	tok := t.tokens[0]
	t.tokens = t.tokens[1:]

	return tok
}

// newTestParser returns a Parser configured with a testLexer emitting a given stream of tokens.
func newTestParser(t *testing.T, tokens []token.Token) *parser.Parser {
	t.Helper()
	lexer := &testLexer{tokens: tokens}
	return parser.New(t.Name(), nil, lexer)
}

func TestParseVarDecl(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []token.Token
		errs    []error
		want    ast.VarDeclaration
		wantErr bool
	}{
		{
			name: "valid",
			tokens: []token.Token{
				{Kind: token.Var, Start: 0, End: 3},
				{Kind: token.Ident, Start: 4, End: 9},
				{Kind: token.Eq, Start: 14, End: 15},
				{Kind: token.Number, Start: 16, End: 17},
				{Kind: token.SemiColon, Start: 17, End: 18},
				{Kind: token.EOF, Start: 18, End: 18},
			},
			want: ast.VarDeclaration{
				Ident: ast.Ident{
					Tok: token.Token{Kind: token.Ident, Start: 4, End: 9},
				},
			},
			wantErr: false,
		},
		{
			name: "missing semicolon",
			tokens: []token.Token{
				{Kind: token.Var, Start: 0, End: 3},
				{Kind: token.Ident, Start: 4, End: 13},
				{Kind: token.Eq, Start: 14, End: 15},
				{Kind: token.Number, Start: 16, End: 17},
				// {Kind: token.SemiColon, Start: 17, End: 18}, // <- This should be here but isn't
				{Kind: token.EOF, Start: 17, End: 17}, // So the EOF occurs at pos 17
			},
			want:    ast.VarDeclaration{},
			wantErr: true,
			errs: []error{parser.SyntaxError{
				File:  "TestParseVarDecl/missing_semicolon",
				Msg:   "unexpected EOF",
				Token: token.Token{Kind: token.EOF, Start: 17, End: 17},
				Line:  1,
				Col:   17,
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParser(t, tt.tokens)

			statement := p.ParseVarDecl()

			if tt.wantErr {
				// If we wanted an error, the statement should be nil, and our errs list
				// should contain the right parse errors
				test.Equal(t, statement, nil)
				test.EqualFunc(t, p.Errors(), tt.errs, syntaxErrorsEqual)
				return
			}

			// We didn't want an error
			test.EqualFunc(t, p.Errors(), nil, slices.Equal, test.Context("found unexpected parser errors"))

			decl, ok := statement.(ast.VarDeclaration)
			test.True(t, ok, test.Context("expected ast.VarDeclaration, got %T", statement))

			test.Equal(t, decl.Ident.Name(), tt.want.Ident.Name())
		})
	}
}

func syntaxErrorEqual(a, b parser.SyntaxError) bool {
	if a.File != b.File {
		return false
	}

	if a.Msg != b.Msg {
		return false
	}

	if a.Line != b.Line {
		return false
	}

	if a.Col != b.Col {
		return false
	}

	if !token.Equal(a.Token, b.Token) {
		return false
	}

	return true
}

func syntaxErrorsEqual(a, b []error) bool {
	if len(a) != len(b) {
		return false
	}

	for i, errA := range a {
		errB := b[i]

		var syntaxErrA parser.SyntaxError
		if !errors.As(errA, &syntaxErrA) {
			return false
		}

		var syntaxErrB parser.SyntaxError
		if !errors.As(errB, &syntaxErrB) {
			return false
		}

		if !syntaxErrorEqual(syntaxErrA, syntaxErrB) {
			return false
		}
	}

	return true
}
