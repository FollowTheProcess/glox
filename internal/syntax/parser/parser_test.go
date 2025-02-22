package parser_test

import (
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
func newTestParser(tokens []token.Token) *parser.Parser {
	lexer := &testLexer{tokens: tokens}
	return parser.New(nil, lexer)
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
				{Kind: token.Var, Text: []byte("var"), Offset: 0, Width: 3},
				{Kind: token.Ident, Text: []byte("something"), Offset: 4, Width: 9},
				{Kind: token.Equal, Text: []byte("="), Offset: 14, Width: 1},
				{Kind: token.Number, Text: []byte("2"), Offset: 16, Width: 1},
				{Kind: token.SemiColon, Text: []byte(";"), Offset: 17, Width: 1},
				{Kind: token.EOF},
			},
			want: ast.VarDeclaration{
				Ident: ast.Ident{
					Tok: token.Token{
						Kind:   token.Ident,
						Text:   []byte("something"),
						Offset: 4,
						Width:  9,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing semicolon",
			tokens: []token.Token{
				{Kind: token.Var, Text: []byte("var"), Offset: 0, Width: 3},
				{Kind: token.Ident, Text: []byte("something"), Offset: 4, Width: 9},
				{Kind: token.Equal, Text: []byte("="), Offset: 14, Width: 1},
				{Kind: token.Number, Text: []byte("2"), Offset: 16, Width: 1},
				{Kind: token.EOF},
			},
			want:    ast.VarDeclaration{},
			wantErr: true,
			errs:    []error{parser.ErrUnexpectedEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParser(tt.tokens)

			statement := p.ParseVarDecl()

			if tt.wantErr {
				// If we wanted an error, the statement should be nil, and our errs list
				// should contain the right parse errors
				test.Equal(t, statement, nil)
				test.EqualFunc(t, p.Errors(), tt.errs, slices.Equal)
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
