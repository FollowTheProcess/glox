package parser_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

func TestParseVarDecl(t *testing.T) {
	tests := []struct {
		name string      // Name of the test case
		src  string      // Source code
		errs []error     // Expected parse errors
		want ast.Program // Expected program
	}{
		{
			name: "valid",
			src:  "var something = 2;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.VarDeclaration{
						Ident: ast.Ident{
							Name: "something",
							Tok:  token.Token{Kind: token.Ident, Start: 4, End: 13},
						},
						// TODO(@FollowTheProcess): Test Value once expressions are implemented
					},
				},
			},
		},
		{
			name: "missing semicolon",
			src:  "var something = 2", // <- no semicolon at the end
			want: ast.Program{},
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
			got, err := p.Parse()

			// Whether or not we wanted an error is encoded in the length of tt.errs:
			// 	0:	No, any error is unexpected and should fail the test
			// 	>0:	Yes, we wanted very specific errors and should test for them
			wantedError := len(tt.errs) != 0
			test.WantErr(t, err, wantedError)

			if wantedError {
				// If we wanted an error, the Program should be empty, and our errs list
				// should contain the right parse errors
				test.Equal(t, len(got.Statements), 0, test.Context("expected empty program"))
				test.EqualFunc(t, p.Errors(), tt.errs, slices.Equal, test.Context("syntax errors did not match"))
				return
			}

			// We didn't want an error
			test.EqualFunc(t, p.Errors(), nil, slices.Equal, test.Context("found unexpected parser errors"))

			test.Equal(
				t,
				len(got.Statements),
				len(tt.want.Statements),
				test.Context("Mismatch in number of statements"),
			)

			for i, wantStatement := range tt.want.Statements {
				gotStatement := got.Statements[i]

				gotDecl, ok := gotStatement.(ast.VarDeclaration)
				test.True(
					t,
					ok,
					test.Context("expected got statement %d to be ast.VarDeclaration got %T: %[2]#v", i, gotStatement))

				wantDecl, ok := wantStatement.(ast.VarDeclaration)
				test.True(
					t,
					ok,
					test.Context("expected want statement %d to be ast.VarDeclaration got %T: %[2]#v", i, wantStatement),
				)

				test.Equal(t, gotDecl.Ident.Name, wantDecl.Ident.Name, test.Context("ident name mismatch"))
				test.Equal(t, gotDecl.Ident.Token(), wantDecl.Ident.Token(), test.Context("ident token mismatch"))
			}
		})
	}
}

func TestParseReturnStatement(t *testing.T) {
	tests := []struct {
		name string      // Name of the test case
		src  string      // Source code
		errs []error     // Expected parse errors
		want ast.Program // Expected program
	}{
		{
			name: "valid",
			src:  "return 3;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ReturnStatement{
						Tok: token.Token{Kind: token.Return, Start: 0, End: 6},
						// TODO(@FollowTheProcess): Test Value once expressions are implemented
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src)
			got, err := p.Parse()

			// Whether or not we wanted an error is encoded in the length of tt.errs:
			// 	0:	No, any error is unexpected and should fail the test
			// 	>0:	Yes, we wanted very specific errors and should test for them
			wantedError := len(tt.errs) != 0
			test.WantErr(t, err, wantedError)

			if wantedError {
				// If we wanted an error, the Program should be empty, and our errs list
				// should contain the right parse errors
				test.Equal(t, len(got.Statements), 0, test.Context("expected empty program"))
				test.EqualFunc(t, p.Errors(), tt.errs, slices.Equal, test.Context("syntax errors did not match"))
				return
			}

			// We didn't want an error
			test.EqualFunc(t, p.Errors(), nil, slices.Equal, test.Context("found unexpected parser errors"))

			test.Equal(
				t,
				len(got.Statements),
				len(tt.want.Statements),
				test.Context("Mismatch in number of statements"),
			)

			for i, wantStatement := range tt.want.Statements {
				gotStatement := got.Statements[i]

				got, ok := gotStatement.(ast.ReturnStatement)
				test.True(
					t,
					ok,
					test.Context("expected got statement %d to be ast.ReturnStatement got %T: %[2]#v", i, gotStatement))

				want, ok := wantStatement.(ast.ReturnStatement)
				test.True(
					t,
					ok,
					test.Context("expected want statement %d to be ast.ReturnStatement got %T: %[2]#v", i, wantStatement),
				)

				test.Equal(t, got.Token(), want.Token(), test.Context("Token() mismatch in ast.ReturnStatement"))
			}
		})
	}
}
