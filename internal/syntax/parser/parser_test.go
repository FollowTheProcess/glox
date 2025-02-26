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

func TestParseUnaryExpression(t *testing.T) {
	tests := []parseTest{
		{
			name: "minus five",
			src:  "-5;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.UnaryExpression{
							Tok: token.Token{Kind: token.Minus, Start: 0, End: 1},
							Value: ast.NumberLiteral{
								Tok:   token.Token{Kind: token.Number, Start: 1, End: 2},
								Value: 5,
							},
						},
						Tok: token.Token{Kind: token.Minus, Start: 0, End: 1},
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

// TODO(@FollowTheProcess): Lots more tests, all the ops and op precedence
// as well as syntax errors and stuff

func TestParseBinaryExpression(t *testing.T) {
	tests := []parseTest{
		{
			name: "add",
			src:  "x + y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.IdentExpression{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Plus, Start: 2, End: 3},
							Right: ast.IdentExpression{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "minus",
			src:  "x - y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.IdentExpression{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Minus, Start: 2, End: 3},
							Right: ast.IdentExpression{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "or",
			src:  "x or y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.IdentExpression{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Or, Start: 2, End: 4},
							Right: ast.IdentExpression{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 5, End: 6}},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
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

func TestOperatorPrecedence(t *testing.T) {
	tests := []parseTest{
		{
			name: "unary beats multiplication",
			src:  "-a * b", // ((-a) * b)
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.UnaryExpression{
								Value: ast.IdentExpression{Name: "a", Tok: token.Token{Kind: token.Ident, Start: 1, End: 2}},
								Tok:   token.Token{Kind: token.Minus, Start: 0, End: 1},
							},
							Right: ast.IdentExpression{Name: "b", Tok: token.Token{Kind: token.Ident, Start: 5, End: 6}},
							Op:    token.Token{Kind: token.Star, Start: 3, End: 4},
						},
						Tok: token.Token{Kind: token.Minus, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "three adds",
			src:  "a + b + c", // ((a + b) + c)
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.BinaryExpression{
								Left:  ast.IdentExpression{Name: "a", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
								Right: ast.IdentExpression{Name: "b", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
								Op:    token.Token{Kind: token.Plus, Start: 2, End: 3},
							},
							Right: ast.IdentExpression{Name: "c", Tok: token.Token{Kind: token.Ident, Start: 8, End: 9}},
							Op:    token.Token{Kind: token.Plus, Start: 6, End: 7},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "two adds and subtract",
			src:  "a + b - c", // ((a + b) - c)
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.BinaryExpression{
								Left:  ast.IdentExpression{Name: "a", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
								Right: ast.IdentExpression{Name: "b", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
								Op:    token.Token{Kind: token.Plus, Start: 2, End: 3},
							},
							Right: ast.IdentExpression{Name: "c", Tok: token.Token{Kind: token.Ident, Start: 8, End: 9}},
							Op:    token.Token{Kind: token.Minus, Start: 6, End: 7},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "three multiplies",
			src:  "a * b * c", // ((a * b) * c)
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.BinaryExpression{
								Left:  ast.IdentExpression{Name: "a", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
								Right: ast.IdentExpression{Name: "b", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
								Op:    token.Token{Kind: token.Star, Start: 2, End: 3},
							},
							Right: ast.IdentExpression{Name: "c", Tok: token.Token{Kind: token.Ident, Start: 8, End: 9}},
							Op:    token.Token{Kind: token.Star, Start: 6, End: 7},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "multiply then divide",
			src:  "a * b / c", // ((a * b) / c)
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.BinaryExpression{
								Left:  ast.IdentExpression{Name: "a", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
								Right: ast.IdentExpression{Name: "b", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
								Op:    token.Token{Kind: token.Star, Start: 2, End: 3},
							},
							Right: ast.IdentExpression{Name: "c", Tok: token.Token{Kind: token.Ident, Start: 8, End: 9}},
							Op:    token.Token{Kind: token.ForwardSlash, Start: 6, End: 7},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "divide beats plus",
			src:  "a + b / c", // (a + (b / c))
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.IdentExpression{Name: "a", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Right: ast.BinaryExpression{
								Left:  ast.IdentExpression{Name: "b", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
								Right: ast.IdentExpression{Name: "c", Tok: token.Token{Kind: token.Ident, Start: 8, End: 9}},
								Op:    token.Token{Kind: token.ForwardSlash, Start: 6, End: 7},
							},
							Op: token.Token{Kind: token.Plus, Start: 2, End: 3},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "lots of things",
			src:  "a + b * c + d / e - f", // (((a + (b * c)) + (d / e)) - f)
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.BinaryExpression{
								Left: ast.BinaryExpression{
									Left: ast.IdentExpression{Name: "a", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
									Right: ast.BinaryExpression{
										Left: ast.IdentExpression{
											Name: "b",
											Tok:  token.Token{Kind: token.Ident, Start: 4, End: 5},
										},
										Right: ast.IdentExpression{
											Name: "c",
											Tok:  token.Token{Kind: token.Ident, Start: 8, End: 9},
										},
										Op: token.Token{Kind: token.Star, Start: 6, End: 7},
									},
									Op: token.Token{Kind: token.Plus, Start: 2, End: 3},
								},
								Right: ast.BinaryExpression{
									Left: ast.IdentExpression{
										Name: "d",
										Tok:  token.Token{Kind: token.Ident, Start: 12, End: 13},
									},
									Right: ast.IdentExpression{
										Name: "e",
										Tok:  token.Token{Kind: token.Ident, Start: 16, End: 17},
									},
									Op: token.Token{Kind: token.ForwardSlash, Start: 14, End: 15},
								},
								Op: token.Token{Kind: token.Plus, Start: 10, End: 11},
							},
							Right: ast.IdentExpression{Name: "f", Tok: token.Token{Kind: token.Ident, Start: 20, End: 21}},
							Op:    token.Token{Kind: token.Minus, Start: 18, End: 19},
						},
						Tok: token.Token{Kind: token.Ident, Start: 0, End: 1},
					},
				},
			},
		},
		{
			name: "two statements",
			src:  "3 + 4; -5 * 5", // (3 + 4)((-5) * 5)
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.NumberLiteral{Value: 3, Tok: token.Token{Kind: token.Number, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Plus, Start: 2, End: 3},
							Right: ast.NumberLiteral{Value: 4, Tok: token.Token{Kind: token.Number, Start: 4, End: 5}},
						},
						Tok: token.Token{Kind: token.Number, Start: 0, End: 1},
					},
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left: ast.UnaryExpression{
								Value: ast.NumberLiteral{Value: 5, Tok: token.Token{Kind: token.Number, Start: 8, End: 9}},
								Tok:   token.Token{Kind: token.Minus, Start: 7, End: 8},
							},
							Op:    token.Token{Kind: token.Star, Start: 10, End: 11},
							Right: ast.NumberLiteral{Value: 5, Tok: token.Token{Kind: token.Number, Start: 12, End: 13}},
						},
						Tok: token.Token{Kind: token.Minus, Start: 7, End: 8},
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
	test.True(tb, ok, test.Context("expected got to be ast.VarStatement, got %T: %#v", statement, statement))

	want, ok := expected.(ast.VarStatement)
	test.True(tb, ok, test.Context("expected want to be ast.VarStatement, got %T: %#v", expected, expected))

	test.Equal(tb, got.Ident.Name, want.Ident.Name, test.Context("ident name mismatch"))
	test.Equal(tb, got.Ident.Token(), want.Ident.Token(), test.Context("ident token mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testReturnStatement tests two [ast.ReturnStatement] nodes for equality, failing the test if they
// are not identical.
func testReturnStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.ReturnStatement)
	test.True(tb, ok, test.Context("expected got to be ast.ReturnStatement, got %T: %#v", statement, statement))

	want, ok := expected.(ast.ReturnStatement)
	test.True(tb, ok, test.Context("expected want to be ast.ReturnStatement, got %T: %#v", expected, expected))

	test.Equal(tb, got.Tok, want.Tok, test.Context("ReturnStatement token mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testPrintStatement tests two [ast.PrintStatement] nodes for equality, failing the test if they
// are not identical.
func testPrintStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.PrintStatement)
	test.True(tb, ok, test.Context("expected got to be ast.PrintStatement, got %T: %#v", statement, statement))

	want, ok := expected.(ast.PrintStatement)
	test.True(tb, ok, test.Context("expected want to be ast.PrintStatement, got %T: %#v", expected, expected))

	test.Equal(tb, got.Tok, want.Tok, test.Context("PrintStatement token mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testExpression tests two [ast.Expression] nodes for equality, failing the test
// if they are not identical.
func testExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	test.NotEqual(tb, expression, nil, test.Context("testExpression expression was nil"))
	test.NotEqual(tb, expected, nil, test.Context("testExpression expected was nil"))

	test.Equal(tb, expression.Token(), expected.Token(), test.Context("Expression token mismatch"))

	switch expected.(type) {
	case ast.NumberLiteral:
		testNumberLiteralExpression(tb, expression, expected)
	case ast.IdentExpression:
		testIdentExpression(tb, expression, expected)
	case ast.UnaryExpression:
		testUnaryExpression(tb, expression, expected)
	case ast.BinaryExpression:
		testBinaryExpression(tb, expression, expected)
	default:
		tb.Fatalf("unhandled ast Expression in testExpression: %T", expected)
	}
}

// testExpressionStatement tests two [ast.ExpressionStatement] nodes for equality, failing the test
// if they are not identical.
func testExpressionStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.ExpressionStatement)
	test.True(tb, ok, test.Context("expected got to be ast.ExpressionStatement, got %T: %#v", statement, statement))

	want, ok := expected.(ast.ExpressionStatement)
	test.True(tb, ok, test.Context("expected want to be ast.ExpressionStatement, got %T: %#v", expected, expected))

	test.Equal(tb, got.Tok, want.Tok, test.Context("ExpressionStatement.Tok mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testIdentExpression tests two [ast.IdentExpression] nodes for equality, failing the test
// if they are not identical.
func testIdentExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.IdentExpression)
	test.True(tb, ok, test.Context("expected got to be ast.IdentExpression, got %T: %#v", expression, expression))

	want, ok := expected.(ast.IdentExpression)
	test.True(tb, ok, test.Context("expected want to be ast.IdentExpression, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("IdentExpression mismatch"))
}

// testNumberLiteralExpression tests two [ast.NumberLiteral] nodes for equality, failing the test
// if they are not identical, used in the context where the number is an expression, as in
// `var x = 5;`.
func testNumberLiteralExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.NumberLiteral)
	test.True(tb, ok, test.Context("expected got to be ast.NumberLiteral, got %T: %#v", expression, expression))

	want, ok := expected.(ast.NumberLiteral)
	test.True(tb, ok, test.Context("expected want to be ast.NumberLiteral, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("NumberLiteral mismatch"))
}

// testUnaryExpression tests two [ast.UnaryExpression] nodes for equality, failing the test
// if they are not identical.
func testUnaryExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.UnaryExpression)
	test.True(tb, ok, test.Context("expected got to be ast.UnaryExpression, got %T: %#v", expression, expression))

	want, ok := expected.(ast.UnaryExpression)
	test.True(tb, ok, test.Context("expected want to be ast.UnaryExpression, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("UnaryExpression mismatch"))
}

// testBinaryExpression tests two [ast.BinaryExpression] nodes for equality, failing the test
// if they are not identical.
func testBinaryExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.BinaryExpression)
	test.True(tb, ok, test.Context("expected got to be ast.BinaryExpression, got %T: %#v", expression, expression))

	want, ok := expected.(ast.BinaryExpression)
	test.True(tb, ok, test.Context("expected want to be ast.BinaryExpression, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("BinaryExpression mismatch"))
}
