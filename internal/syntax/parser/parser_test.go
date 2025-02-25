package parser_test

import (
	"slices"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

// parseTest is a table driven test where we parse a full program and make assertions
// about the AST that gets produced.
type parseTest struct {
	name string      // name of the test case
	src  string      // source code to lex and parse
	errs []error     // Expected parse errors, empty signals no errors expected
	want ast.Program // The full AST we expect to see
}

func TestParseVarStatement(t *testing.T) {
	tests := []parseTest{
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
						Value: ast.NumberLiteral{
							Value: 2,
							Tok:   token.Token{Kind: token.Number, Start: 16, End: 17},
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
				File:  "TestParseVarStatement/missing_semicolon",
				Msg:   `expected ";", got EOF: ""`, // TODO(@FollowTheProcess): If it's EOF, don't show the (empty) value
				Token: token.Token{Kind: token.Number, Start: 16, End: 17},
				Line:  1,
				Col:   16,
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

			testParse(t, got, tt.want)
		})
	}
}

func TestParseReturnStatement(t *testing.T) {
	tests := []parseTest{
		{
			name: "valid",
			src:  "return 3;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ReturnStatement{
						Tok: token.Token{Kind: token.Return, Start: 0, End: 6},
						Value: ast.NumberLiteral{
							Value: 3,
							Tok:   token.Token{Kind: token.Number, Start: 7, End: 8},
						},
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

			testParse(t, got, tt.want)
		})
	}
}

func TestParsePrintStatement(t *testing.T) {
	tests := []parseTest{
		{
			name: "valid",
			src:  "print 3.14159;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.PrintStatement{
						Tok: token.Token{Kind: token.Print, Start: 0, End: 5},
						Value: ast.NumberLiteral{
							Value: 3.14159,
							Tok:   token.Token{Kind: token.Number, Start: 6, End: 13},
						},
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

			testParse(t, got, tt.want)
		})
	}
}

func TestParseIdentifierExpression(t *testing.T) {
	tests := []parseTest{
		{
			name: "valid",
			src:  "foo;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.IdentExpression{
							Name: "foo",
							Tok:  token.Token{Kind: token.Ident, Start: 0, End: 3},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 3},
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

			testParse(t, got, tt.want)
		})
	}
}

func TestParseNumberLiteral(t *testing.T) {
	tests := []parseTest{
		{
			name: "integer",
			src:  "5",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.NumberLiteral{
							Value: 5,
							Tok:   token.Token{Kind: token.Number, Start: 0, End: 1},
						},
						Tok: token.Token{Kind: token.Number, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "bigger integer",
			src:  "9463",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.NumberLiteral{
							Value: 9463,
							Tok:   token.Token{Kind: token.Number, Start: 0, End: 4},
						},
						Tok: token.Token{Kind: token.Number, Start: 0, End: 4},
					},
				},
			},
		},
		{
			name: "float",
			src:  "3.14159",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.NumberLiteral{
							Value: 3.14159,
							Tok:   token.Token{Kind: token.Number, Start: 0, End: 7},
						},
						Tok: token.Token{Kind: token.Number, Start: 0, End: 7},
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

			testParse(t, got, tt.want)
		})
	}
}

// testParse tests two ast.Programs are identical, failing the test if not.
func testParse(tb testing.TB, got, want ast.Program) {
	tb.Helper()

	test.Equal(tb, len(got.Statements), len(want.Statements), test.Context("mismatch in number of statements"))

	for index, wantStatement := range want.Statements {
		gotStatement := got.Statements[index]

		test.NotEqual(tb, gotStatement, nil, test.Context("testParse gotStatement was nil"))
		test.NotEqual(tb, wantStatement, nil, test.Context("testParse wantStatement was nil"))

		switch wantStatement := wantStatement.(type) {
		case ast.VarStatement:
			testVarStatement(tb, gotStatement, wantStatement)
		case ast.ReturnStatement:
			testReturnStatement(tb, gotStatement, wantStatement)
		case ast.PrintStatement:
			testPrintStatement(tb, gotStatement, wantStatement)
		case ast.ExpressionStatement:
			testExpressionStatement(tb, gotStatement, wantStatement)
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

	test.Equal(tb, got.Ident.Name, want.Ident.Name, test.Context("ident name mismatch"))
	test.Equal(tb, got.Ident.Token(), want.Ident.Token(), test.Context("ident token mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testReturnStatement tests two [ast.ReturnStatement] nodes for equality, failing the test if they
// are not identical.
func testReturnStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.ReturnStatement)
	test.True(tb, ok, test.Context("expected got to be ast.ReturnStatement, got %T: %[1]#v", statement))

	want, ok := expected.(ast.ReturnStatement)
	test.True(tb, ok, test.Context("expected want to be ast.ReturnStatement, got %T: %[1]#v", expected))

	test.Equal(tb, got.Tok, want.Tok, test.Context("ReturnStatement token mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testReturnStatement tests two [ast.PrintStatement] nodes for equality, failing the test if they
// are not identical.
func testPrintStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.PrintStatement)
	test.True(tb, ok, test.Context("expected got to be ast.PrintStatement, got %T: %[1]#v", statement))

	want, ok := expected.(ast.PrintStatement)
	test.True(tb, ok, test.Context("expected want to be ast.PrintStatement, got %T: %[1]#v", expected))

	test.Equal(tb, got.Tok, want.Tok, test.Context("PrintStatement token mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testExpression tests two [ast.Expression] nodes for equality, failing the test
// if they are not identical.
func testExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	test.NotEqual(tb, expression, nil, test.Context("testExpression expression was nil"))
	test.NotEqual(tb, expected, nil, test.Context("testExpression expected was nil"))

	switch expected.(type) {
	case ast.NumberLiteral:
		testNumberLiteralExpression(tb, expression, expected)
	default:
		tb.Fatalf("unhandled ast Node in testExpression: %T", expected)
	}
}

// testExpressionStatement tests two [ast.ExpressionStatement] nodes for equality, failing the test
// if they are not identical.
func testExpressionStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.ExpressionStatement)
	test.True(tb, ok, test.Context("expected got to be ast.ExpressionStatement, got %T: %[1]#v", statement))

	want, ok := expected.(ast.ExpressionStatement)
	test.True(tb, ok, test.Context("expected want to be ast.ExpressionStatement, got %T: %[1]#v", expected))

	switch want.Value.(type) {
	case ast.IdentExpression:
		testIdentExpression(tb, got, want)
	case ast.NumberLiteral:
		testNumberLiteralStatement(tb, got, want)
	default:
		tb.Fatalf("unhandled ast Node in testExpressionStatement: %T", want.Value)
	}
}

// testIdentExpression tests two [ast.IdentExpression] nodes for equality, failing the test
// if they are not identical.
func testIdentExpression(tb testing.TB, statement, expected ast.ExpressionStatement) {
	tb.Helper()

	test.Equal(tb, statement.Tok, expected.Tok, test.Context("ExpressionStatement token mismatch"))

	got, ok := statement.Value.(ast.IdentExpression)
	test.True(tb, ok, test.Context("expected got to be ast.IdentExpression, got %T: %[1]#v", statement.Value))

	want, ok := expected.Value.(ast.IdentExpression)
	test.True(tb, ok, test.Context("expected want to be ast.IdentExpression, got %T: %[1]#v", expected.Value))

	test.Equal(tb, got, want, test.Context("IdentExpression mismatch"))
}

// testNumberLiteralStatement tests two [ast.NumberLiteral] nodes for equality, failing the test
// if they are not identical, used in the context where the number is an expression statement as
// in `5;`.
func testNumberLiteralStatement(tb testing.TB, statement, expected ast.ExpressionStatement) {
	tb.Helper()

	test.Equal(tb, statement.Token(), expected.Token(), test.Context("Expression token mismatch"))

	got, ok := statement.Value.(ast.NumberLiteral)
	test.True(tb, ok, test.Context("expected got to be ast.NumberLiteral, got %T: %[1]#v", statement.Value))

	want, ok := expected.Value.(ast.NumberLiteral)
	test.True(tb, ok, test.Context("expected want to be ast.NumberLiteral, got %T: %[1]#v", expected.Value))

	test.Equal(tb, got, want, test.Context("NumberLiteral mismatch"))
}

// testNumberLiteralExpression tests two [ast.NumberLiteral] nodes for equality, failing the test
// if they are not identical, used in the context where the number is an expression, as in
// `var x = 5;`.
func testNumberLiteralExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	test.Equal(tb, expression.Token(), expected.Token(), test.Context("Expression token mismatch"))

	got, ok := expression.(ast.NumberLiteral)
	test.True(tb, ok, test.Context("expected got to be ast.NumberLiteral, got %T: %[1]#v", expression))

	want, ok := expected.(ast.NumberLiteral)
	test.True(tb, ok, test.Context("expected want to be ast.NumberLiteral, got %T: %[1]#v", expected))

	test.Equal(tb, got, want, test.Context("NumberLiteral mismatch"))
}
