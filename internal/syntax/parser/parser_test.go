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

func TestParseVarDecl(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		errs    []error
		want    ast.VarDeclaration
		wantErr bool
	}{
		{
			name: "valid",
			src:  "var something = 2;",
			want: ast.VarDeclaration{
				Ident: ast.Ident{
					Name: "something",
					Tok:  token.Token{Kind: token.Ident, Start: 4, End: 9},
				},
			},
			wantErr: false,
		},
		{
			name:    "missing semicolon",
			src:     "var something = 2", // <- no semicolon at the end
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
			p := parser.New(t.Name(), tt.src)

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

			test.Equal(t, decl.Ident.Name, tt.want.Ident.Name)
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

	if a.Token != b.Token {
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
