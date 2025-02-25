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
					ast.VarStatement{
						Ident: ast.IdentExpression{
							Name: "something",
							Tok:  token.Token{Kind: token.Ident, Start: 4, End: 13},
						},
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

			parseTest(t, got, tt.want)
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

			parseTest(t, got, tt.want)
		})
	}
}

func TestParsePrintStatement(t *testing.T) {
	tests := []struct {
		name string      // Name of the test case
		src  string      // Source code
		errs []error     // Expected parse errors
		want ast.Program // Expected program
	}{
		{
			name: "valid",
			src:  "print x == y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.PrintStatement{
						Tok: token.Token{Kind: token.Print, Start: 0, End: 5},
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

			parseTest(t, got, tt.want)
		})
	}
}

// parseTest tests two ast.Programs are identical, failing the test if not.
func parseTest(tb testing.TB, got, want ast.Program) {
	tb.Helper()

	test.Equal(tb, len(got.Statements), len(want.Statements), test.Context("mismatch in number of statements"))

	for index, wantStatement := range want.Statements {
		gotStatement := got.Statements[index]

		switch wantStatement := wantStatement.(type) {
		case ast.VarStatement:
			testVarStatement(tb, gotStatement, wantStatement)
		case ast.ReturnStatement:
			testReturnStatement(tb, gotStatement, wantStatement)
		case ast.PrintStatement:
			testPrintStatement(tb, gotStatement, wantStatement)
		default:
			tb.Fatalf("unhandled ast Node in parseTest: %T", wantStatement)
		}
	}
}

// testVarStatement tests two [ast.VarStatement] nodes for equality, failing the test if they
// are not identical.
func testVarStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.VarStatement)
	test.True(tb, ok, test.Context("expected got to be ast.VarStatement, got %T: %[1]#v", statement))

	want, ok := expected.(ast.VarStatement)
	test.True(tb, ok, test.Context("expected want to be ast.VarStatement, got %T: %[1]#v", expected))

	// TODO(@FollowTheProcess): Test Value once expressions are implemented
	test.Equal(tb, got.Ident.Name, want.Ident.Name, test.Context("ident name mismatch"))
	test.Equal(tb, got.Ident.Token(), want.Ident.Token(), test.Context("ident token mismatch"))
}

// testReturnStatement test two [ast.ReturnStatement] nodes for equality, failing the test if they
// are not identical.
func testReturnStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.ReturnStatement)
	test.True(tb, ok, test.Context("expected got to be ast.ReturnStatement, got %T: %[1]#v", statement))

	want, ok := expected.(ast.ReturnStatement)
	test.True(tb, ok, test.Context("expected want to be ast.ReturnStatement, got %T: %[1]#v", expected))

	// TODO(@FollowTheProcess): Test Value once expressions are implemented
	test.Equal(tb, got.Tok, want.Tok, test.Context("ReturnStatement token mismatch"))
}

// testReturnStatement test two [ast.PrintStatement] nodes for equality, failing the test if they
// are not identical.
func testPrintStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.PrintStatement)
	test.True(tb, ok, test.Context("expected got to be ast.PrintStatement, got %T: %[1]#v", statement))

	want, ok := expected.(ast.PrintStatement)
	test.True(tb, ok, test.Context("expected want to be ast.PrintStatement, got %T: %[1]#v", expected))

	// TODO(@FollowTheProcess): Test Value once expressions are implemented
	test.Equal(tb, got.Tok, want.Tok, test.Context("PrintStatement token mismatch"))
}
