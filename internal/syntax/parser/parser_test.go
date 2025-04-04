package parser_test

import (
	"flag"
	"os"
	"slices"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

var debug = flag.Bool("debug", false, "Emit parse traces to stderr during tests")

// parseTest is a table driven test where we parse a full program and make assertions
// about the AST that gets produced.
type parseTest struct {
	name string      // name of the test case
	src  string      // source code to lex and parse
	errs []error     // Expected parse errors, empty signals no errors expected
	want ast.Program // The full AST we expect to see
}

func TestParseVarDeclaration(t *testing.T) {
	tests := []parseTest{
		{
			name: "valid",
			src:  "var something = 2;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.DeclarationStatement{
						Value: ast.VarDeclaration{
							Ident: ast.Ident{
								Name: "something",
								Tok:  token.Token{Kind: token.Ident, Start: 4, End: 13},
							},
							Value: ast.Number{
								Value: 2,
								Tok:   token.Token{Kind: token.Number, Start: 16, End: 17},
							},
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
				File:  "TestParseVarDeclaration/missing_semicolon",
				Msg:   `expected ";", got EOF`,
				Token: token.Token{Kind: token.Number, Start: 16, End: 17},
				Line:  1,
				Col:   18, // Columns in editors start at 1, so end offset + 1
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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
						Value: ast.Number{
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
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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
						Value: ast.Number{
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
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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

func TestParseIdent(t *testing.T) {
	tests := []parseTest{
		{
			name: "valid",
			src:  "foo;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.Ident{
							Name: "foo",
							Tok:  token.Token{Kind: token.Ident, Start: 0, End: 3},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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

func TestParseNumber(t *testing.T) {
	tests := []parseTest{
		{
			name: "integer",
			src:  "5;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.Number{
							Value: 5,
							Tok:   token.Token{Kind: token.Number, Start: 0, End: 1},
						},
					},
				},
			},
		},
		{
			name: "bigger integer",
			src:  "9463;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.Number{
							Value: 9463,
							Tok:   token.Token{Kind: token.Number, Start: 0, End: 4},
						},
					},
				},
			},
		},
		{
			name: "float",
			src:  "3.14159;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.Number{
							Value: 3.14159,
							Tok:   token.Token{Kind: token.Number, Start: 0, End: 7},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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

func TestParseBool(t *testing.T) {
	tests := []parseTest{
		{
			name: "true",
			src:  "true;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.Bool{
							Value: true,
							Tok:   token.Token{Kind: token.True, Start: 0, End: 4},
						},
					},
				},
			},
		},
		{
			name: "false",
			src:  "false;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.Bool{
							Value: false,
							Tok:   token.Token{Kind: token.False, Start: 0, End: 5},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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

func TestParseString(t *testing.T) {
	tests := []parseTest{
		{
			name: "valid",
			src:  `"I'm a string literal";`,
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.String{
							Value: "I'm a string literal",
							Tok:   token.Token{Kind: token.String, Start: 0, End: 22},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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
							Value: ast.Number{
								Tok:   token.Token{Kind: token.Number, Start: 1, End: 2},
								Value: 5,
							},
						},
					},
				},
			},
		},
		{
			name: "not true",
			src:  "!true;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.UnaryExpression{
							Value: ast.Bool{
								Value: true,
								Tok:   token.Token{Kind: token.True, Start: 1, End: 5},
							},
							Tok: token.Token{Kind: token.Bang, Start: 0, End: 1},
						},
					},
				},
			},
		},
		{
			name: "not false",
			src:  "!false;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.UnaryExpression{
							Value: ast.Bool{
								Value: false,
								Tok:   token.Token{Kind: token.False, Start: 1, End: 6},
							},
							Tok: token.Token{Kind: token.Bang, Start: 0, End: 1},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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

func TestParseBinaryExpression(t *testing.T) {
	tests := []parseTest{
		{
			name: "add",
			src:  "x + y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Plus, Start: 2, End: 3},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
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
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Minus, Start: 2, End: 3},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
					},
				},
			},
		},
		{
			name: "multiply",
			src:  "x * y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Star, Start: 2, End: 3},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
					},
				},
			},
		},
		{
			name: "divide",
			src:  "x / y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.ForwardSlash, Start: 2, End: 3},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
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
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.Or, Start: 2, End: 4},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 5, End: 6}},
						},
					},
				},
			},
		},
		{
			name: "and",
			src:  "x and y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.And, Start: 2, End: 5},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 6, End: 7}},
						},
					},
				},
			},
		},
		{
			name: "greater",
			src:  "x > y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.GreaterThan, Start: 2, End: 3},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
					},
				},
			},
		},
		{
			name: "less",
			src:  "x < y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.LessThan, Start: 2, End: 3},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 4, End: 5}},
						},
					},
				},
			},
		},
		{
			name: "greater or equal",
			src:  "x >= y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.GreaterThanEq, Start: 2, End: 4},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 5, End: 6}},
						},
					},
				},
			},
		},
		{
			name: "less or equal",
			src:  "x <= y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.LessThanEq, Start: 2, End: 4},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 5, End: 6}},
						},
					},
				},
			},
		},
		{
			name: "equal",
			src:  "x == y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.DoubleEq, Start: 2, End: 4},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 5, End: 6}},
						},
					},
				},
			},
		},
		{
			name: "not equal",
			src:  "x != y;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Ident{Name: "x", Tok: token.Token{Kind: token.Ident, Start: 0, End: 1}},
							Op:    token.Token{Kind: token.BangEq, Start: 2, End: 4},
							Right: ast.Ident{Name: "y", Tok: token.Token{Kind: token.Ident, Start: 5, End: 6}},
						},
					},
				},
			},
		},
		{
			name: "bool equal true",
			src:  "true == true;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Bool{Value: true, Tok: token.Token{Kind: token.True, Start: 0, End: 4}},
							Op:    token.Token{Kind: token.DoubleEq, Start: 5, End: 7},
							Right: ast.Bool{Value: true, Tok: token.Token{Kind: token.True, Start: 8, End: 12}},
						},
					},
				},
			},
		},
		{
			name: "bool not equal",
			src:  "true != false;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Bool{Value: true, Tok: token.Token{Kind: token.True, Start: 0, End: 4}},
							Op:    token.Token{Kind: token.BangEq, Start: 5, End: 7},
							Right: ast.Bool{Value: false, Tok: token.Token{Kind: token.False, Start: 8, End: 13}},
						},
					},
				},
			},
		},
		{
			name: "bool equal true",
			src:  "false == false;",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.BinaryExpression{
							Left:  ast.Bool{Value: false, Tok: token.Token{Kind: token.False, Start: 0, End: 5}},
							Op:    token.Token{Kind: token.DoubleEq, Start: 6, End: 8},
							Right: ast.Bool{Value: false, Tok: token.Token{Kind: token.False, Start: 9, End: 14}},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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

func TestParseGroupedExpression(t *testing.T) {
	tests := []parseTest{
		{
			name: "unary",
			src:  "(-5);",
			want: ast.Program{
				Statements: []ast.Statement{
					ast.ExpressionStatement{
						Value: ast.GroupedExpression{
							LParen: token.Token{Kind: token.OpenParen, Start: 0, End: 1},
							RParen: token.Token{Kind: token.CloseParen, Start: 3, End: 4},
							Value: ast.UnaryExpression{
								Value: ast.Number{
									Value: 5,
									Tok:   token.Token{Kind: token.Number, Start: 2, End: 3},
								},
								Tok: token.Token{Kind: token.Minus, Start: 1, End: 2},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
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
	tests := []struct {
		name string // Name of the test case
		src  string // Source code
		want string // Expected precedence string
	}{
		{
			name: "unary beats multiply",
			src:  "-a * b;",
			want: "((-a) * b)",
		},
		{
			name: "unary beats negation",
			src:  "!-a;",
			want: "(!(-a))",
		},
		{
			name: "three adds",
			src:  "a + b + c;",
			want: "((a + b) + c)",
		},
		{
			name: "add then subtract",
			src:  "a + b - c;",
			want: "((a + b) - c)",
		},
		{
			name: "three multiplies",
			src:  "a * b * c;",
			want: "((a * b) * c)",
		},
		{
			name: "multiply beats divide",
			src:  "a * b / c;",
			want: "((a * b) / c)",
		},
		{
			name: "divide beats add",
			src:  "a + b / c;",
			want: "(a + (b / c))",
		},
		{
			name: "lots of stuff",
			src:  "a + b * c + d / e - f;",
			want: "(((a + (b * c)) + (d / e)) - f)",
		},
		{
			name: "binary comparison",
			src:  "5 > 4 == 3 < 4;",
			want: "((5 > 4) == (3 < 4))",
		},
		{
			name: "binary comparison 2",
			src:  "5 < 4 != 3 > 4;",
			want: "((5 < 4) != (3 > 4))",
		},
		{
			name: "two complex expressions equal",
			src:  "3 + 4 * 5 == 3 * 1 + 4 * 5;",
			want: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			name: "bool true",
			src:  "true;",
			want: "true",
		},
		{
			name: "bool false",
			src:  "false;",
			want: "false",
		},
		{
			name: "bool comparison false",
			src:  "3 > 5 == false;",
			want: "((3 > 5) == false)",
		},
		{
			name: "bool comparison true",
			src:  "3 < 5 == true;",
			want: "((3 < 5) == true)",
		},
		{
			name: "grouped add",
			src:  "1 + (2 + 3) + 4;",
			want: "((1 + (2 + 3)) + 4)",
		},
		{
			name: "grouped add beats multiply",
			src:  "(5 + 5) * 2;",
			want: "((5 + 5) * 2)",
		},
		{
			name: "grouped add beats divide",
			src:  "2 / (5 + 5);",
			want: "(2 / (5 + 5))",
		},
		{
			name: "group beats unary negate",
			src:  "-(5 + 5);",
			want: "(-(5 + 5))",
		},
		{
			name: "group beats unary not",
			src:  "!(true == true);",
			want: "(!(true == true))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(t.Name(), tt.src, *debug, os.Stderr)
			prog, err := p.Parse()
			test.Ok(t, err, test.Context("unexpected parse error"))

			for _, statement := range prog.Statements {
				expressionStatement, ok := statement.(ast.ExpressionStatement)
				test.True(t, ok, test.Context("non ExpressionStatement node found: %T", statement))

				test.Equal(t, expressionStatement.Value.Precedence(), tt.want)
			}
		})
	}
}

// testParse tests two ast.Programs are identical, failing the test if not.
//
// It recursively descends the grammar (much like the actual parser!) asserting the nodes
// and tokens/positions are identical at every stage.
func testParse(tb testing.TB, got, want ast.Program) {
	tb.Helper()

	test.Equal(tb, len(got.Statements), len(want.Statements), test.Context("mismatch in number of statements"))

	for index, wantStatement := range want.Statements {
		gotStatement := got.Statements[index]

		test.NotEqual(tb, gotStatement, nil, test.Context("testParse gotStatement was nil"))
		test.NotEqual(tb, wantStatement, nil, test.Context("testParse wantStatement was nil"))

		switch wantStatement := wantStatement.(type) {
		case ast.ReturnStatement:
			testReturnStatement(tb, gotStatement, wantStatement)
		case ast.PrintStatement:
			testPrintStatement(tb, gotStatement, wantStatement)
		case ast.ExpressionStatement:
			testExpressionStatement(tb, gotStatement, wantStatement)
		case ast.DeclarationStatement:
			testDeclarationStatement(tb, gotStatement, wantStatement)
		default:
			tb.Fatalf("unhandled ast Node in parseTest: %T", wantStatement)
		}
	}
}

// testVarDeclaration tests two [ast.VarDeclaration] nodes for equality, failing the test if they
// are not identical.
func testVarDeclaration(tb testing.TB, declaration, expected ast.Declaration) {
	tb.Helper()

	got, ok := declaration.(ast.VarDeclaration)
	test.True(tb, ok, test.Context("expected got to be ast.VarDeclaration, got %T: %#v", declaration, declaration))

	want, ok := expected.(ast.VarDeclaration)
	test.True(tb, ok, test.Context("expected want to be ast.VarDeclaration, got %T: %#v", expected, expected))

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

// testExpressionStatement tests two [ast.ExpressionStatement] nodes for equality, failing the test
// if they are not identical.
func testExpressionStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.ExpressionStatement)
	test.True(tb, ok, test.Context("expected got to be ast.ExpressionStatement, got %T: %#v", statement, statement))

	want, ok := expected.(ast.ExpressionStatement)
	test.True(tb, ok, test.Context("expected want to be ast.ExpressionStatement, got %T: %#v", expected, expected))

	test.Equal(tb, got.Token(), want.Token(), test.Context("ExpressionStatement.Token() mismatch"))

	testExpression(tb, got.Value, want.Value)
}

// testDeclarationStatement tests two [ast.DeclarationStatement] nodes for equality, failing the test
// if they are not identical.
func testDeclarationStatement(tb testing.TB, statement, expected ast.Statement) {
	tb.Helper()

	got, ok := statement.(ast.DeclarationStatement)
	test.True(tb, ok, test.Context("expected got to be ast.DeclarationStatement, got %T: %#v", statement, statement))

	want, ok := expected.(ast.DeclarationStatement)
	test.True(tb, ok, test.Context("expected want to be ast.DeclarationStatement, got %T: %#v", expected, expected))

	test.Equal(tb, got.Token(), want.Token(), test.Context("DeclarationStatement.Token() mismatch"))

	testDeclaration(tb, got.Value, want.Value)
}

// testExpression tests two [ast.Expression] nodes for equality, failing the test
// if they are not identical.
func testExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	test.NotEqual(tb, expression, nil, test.Context("testExpression expression was nil"))
	test.NotEqual(tb, expected, nil, test.Context("testExpression expected was nil"))

	test.Equal(tb, expression.Token(), expected.Token(), test.Context("Expression token mismatch"))

	switch expected.(type) {
	case ast.Number:
		testNumber(tb, expression, expected)
	case ast.Bool:
		testBool(tb, expression, expected)
	case ast.Ident:
		testIdent(tb, expression, expected)
	case ast.String:
		testString(tb, expression, expected)
	case ast.UnaryExpression:
		testUnaryExpression(tb, expression, expected)
	case ast.BinaryExpression:
		testBinaryExpression(tb, expression, expected)
	case ast.GroupedExpression:
		testGroupedExpression(tb, expression, expected)
	default:
		tb.Fatalf("unhandled ast Expression in testExpression: %T", expected)
	}
}

// testDeclaration tests two [ast.Declaration] nodes for equality, failing the test
// if they are not identical.
func testDeclaration(tb testing.TB, declaration, expected ast.Declaration) {
	tb.Helper()

	test.NotEqual(tb, declaration, nil, test.Context("testDeclaration declaration was nil"))
	test.NotEqual(tb, expected, nil, test.Context("testDeclaration expected was nil"))

	test.Equal(tb, declaration.Token(), expected.Token(), test.Context("Declaration token mismatch"))

	switch expected.(type) {
	case ast.VarDeclaration:
		testVarDeclaration(tb, declaration, expected)
	default:
		tb.Fatalf("unhandled ast Declaration in testDeclaration: %T", expected)
	}
}

// testIdent tests two [ast.Ident] nodes for equality, failing the test
// if they are not identical.
func testIdent(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.Ident)
	test.True(tb, ok, test.Context("expected got to be ast.Ident, got %T: %#v", expression, expression))

	want, ok := expected.(ast.Ident)
	test.True(tb, ok, test.Context("expected want to be ast.Ident, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("Ident mismatch"))
}

// testNumber tests two [ast.Number] nodes for equality, failing the test
// if they are not identical.
func testNumber(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.Number)
	test.True(tb, ok, test.Context("expected got to be ast.Number, got %T: %#v", expression, expression))

	want, ok := expected.(ast.Number)
	test.True(tb, ok, test.Context("expected want to be ast.Number, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("Number mismatch"))
}

// testBool tests two [ast.Bool] nodes for equality, failing the test
// if they are not identical.
func testBool(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.Bool)
	test.True(tb, ok, test.Context("expected got to be ast.Bool, got %T: %#v", expression, expression))

	want, ok := expected.(ast.Bool)
	test.True(tb, ok, test.Context("expected want to be ast.Bool, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("Bool mismatch"))
}

// testString tests two [ast.String] nodes for equality, failing the test
// if they are not identical.
func testString(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.String)
	test.True(tb, ok, test.Context("expected got to be ast.String, got %T: %#v", expression, expression))

	want, ok := expected.(ast.String)
	test.True(tb, ok, test.Context("expected want to be ast.String, got %T: %#v", expected, expected))

	test.Equal(tb, got, want, test.Context("String mismatch"))
}

// testUnaryExpression tests two [ast.Bool] nodes for equality, failing the test
// if they are not identical.
func testUnaryExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.UnaryExpression)
	test.True(tb, ok, test.Context("expected got to be ast.UnaryExpression, got %T: %#v", expression, expression))

	want, ok := expected.(ast.UnaryExpression)
	test.True(tb, ok, test.Context("expected want to be ast.UnaryExpression, got %T: %#v", expected, expected))

	test.Equal(tb, got.Tok, want.Tok, test.Context("UnaryExpression operator mismatch"))
	testExpression(tb, got.Value, want.Value)
}

// testBinaryExpression tests two [ast.BinaryExpression] nodes for equality, failing the test
// if they are not identical.
func testBinaryExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.BinaryExpression)
	test.True(tb, ok, test.Context("expected got to be ast.BinaryExpression, got %T: %#v", expression, expression))

	want, ok := expected.(ast.BinaryExpression)
	test.True(tb, ok, test.Context("expected want to be ast.BinaryExpression, got %T: %#v", expected, expected))

	testExpression(tb, got.Left, want.Left)
	test.Equal(tb, got.Op, want.Op, test.Context("BinaryExpression operator mismatch"))
	testExpression(tb, got.Right, want.Right)
}

// testGroupedExpression tests two [ast.GroupedExpression] nodes for equality, failing the test
// if they are not identical.
func testGroupedExpression(tb testing.TB, expression, expected ast.Expression) {
	tb.Helper()

	got, ok := expression.(ast.GroupedExpression)
	test.True(tb, ok, test.Context("expected got to be ast.GroupedExpression, got %T: %#v", expression, expression))

	want, ok := expected.(ast.GroupedExpression)
	test.True(tb, ok, test.Context("expected want to be ast.GroupedExpression, got %T: %#v", expected, expected))

	test.Equal(tb, got.LParen, want.LParen, test.Context("GroupedExpression LParen mismatch"))
	testExpression(tb, got.Value, want.Value)
	test.Equal(tb, got.RParen, want.RParen, test.Context("GroupedExpression RParen mismatch"))
}
