package interpreter_test

import (
	"testing"

	"github.com/FollowTheProcess/glox/internal/interpreter"
	"github.com/FollowTheProcess/glox/internal/syntax/ast"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/glox/internal/syntax/types"
	"github.com/FollowTheProcess/test"
)

func TestEval(t *testing.T) {
	tests := []struct {
		node ast.Node   // The AST node to evaluate
		want types.Type // The expected evaluated type
		name string     // Name of the test case
	}{
		{
			name: "number int",
			node: ast.Number{Value: 5},
			want: types.Number{Value: 5},
		},
		{
			name: "number float",
			node: ast.Number{Value: 3.14159},
			want: types.Number{Value: 3.14159},
		},
		{
			name: "bool true",
			node: ast.Bool{Value: true},
			want: types.True,
		},
		{
			name: "bool false",
			node: ast.Bool{Value: false},
			want: types.False,
		},
		{
			name: "string literal",
			node: ast.String{Value: "a string"},
			want: types.String{Value: "a string"},
		},
		{
			name: "grouped number",
			node: ast.GroupedExpression{Value: ast.Number{Value: 1}},
			want: types.Number{Value: 1},
		},
		{
			name: "unary negated number",
			node: ast.UnaryExpression{
				Value: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number}, // Note: position information not important here
				},
				Tok: token.Token{Kind: token.Minus},
			},
			want: types.Number{Value: -5},
		},
		{
			name: "unary not true",
			node: ast.UnaryExpression{
				Value: ast.Bool{
					Value: true,
					Tok:   token.Token{Kind: token.True},
				},
				Tok: token.Token{Kind: token.Bang},
			},
			want: types.Bool{Value: false},
		},
		{
			name: "binary add",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 2,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.Plus},
			},
			want: types.Number{Value: 7},
		},
		{
			name: "binary subtract",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 2,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.Minus},
			},
			want: types.Number{Value: 3},
		},
		{
			name: "binary division",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 10,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 2,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.ForwardSlash},
			},
			want: types.Number{Value: 5},
		},
		{
			name: "binary multiply",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 10,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 2,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.Star},
			},
			want: types.Number{Value: 20},
		},
		{
			name: "greater than true",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 3,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.GreaterThan},
			},
			want: types.True,
		},
		{
			name: "greater than false",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 3,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.GreaterThan},
			},
			want: types.False,
		},
		{
			name: "greater than eq true greater",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 3,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.GreaterThanEq},
			},
			want: types.True,
		},
		{
			name: "greater than eq true equal",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.GreaterThanEq},
			},
			want: types.True,
		},
		{
			name: "greater than eq false",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 3,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.GreaterThanEq},
			},
			want: types.False,
		},
		{
			name: "less than true",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 1,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 3,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.LessThan},
			},
			want: types.True,
		},
		{
			name: "less than false",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 3,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 1,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.LessThan},
			},
			want: types.False,
		},
		{
			name: "less than eq true less",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 1,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 3,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.LessThanEq},
			},
			want: types.True,
		},
		{
			name: "less than eq true equal",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 5,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.LessThanEq},
			},
			want: types.True,
		},
		{
			name: "less than eq false",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 7,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 4,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.LessThanEq},
			},
			want: types.False,
		},
		{
			name: "equal numbers true",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 42,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 42,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.DoubleEq},
			},
			want: types.True,
		},
		{
			name: "equal numbers false",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 69,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 42,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.DoubleEq},
			},
			want: types.False,
		},
		{
			name: "equal strings true",
			node: ast.BinaryExpression{
				Left: ast.String{
					Value: "yes",
					Tok:   token.Token{Kind: token.String},
				},
				Right: ast.String{
					Value: "yes",
					Tok:   token.Token{Kind: token.String},
				},
				Op: token.Token{Kind: token.DoubleEq},
			},
			want: types.True,
		},
		{
			name: "equal strings false",
			node: ast.BinaryExpression{
				Left: ast.String{
					Value: "yes",
					Tok:   token.Token{Kind: token.String},
				},
				Right: ast.String{
					Value: "no",
					Tok:   token.Token{Kind: token.String},
				},
				Op: token.Token{Kind: token.DoubleEq},
			},
			want: types.False,
		},
		{
			name: "not equal numbers true",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 42,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 69,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.BangEq},
			},
			want: types.True,
		},
		{
			name: "not equal numbers false",
			node: ast.BinaryExpression{
				Left: ast.Number{
					Value: 69,
					Tok:   token.Token{Kind: token.Number},
				},
				Right: ast.Number{
					Value: 69,
					Tok:   token.Token{Kind: token.Number},
				},
				Op: token.Token{Kind: token.BangEq},
			},
			want: types.False,
		},
		{
			name: "not equal strings true",
			node: ast.BinaryExpression{
				Left: ast.String{
					Value: "yes",
					Tok:   token.Token{Kind: token.String},
				},
				Right: ast.String{
					Value: "no",
					Tok:   token.Token{Kind: token.String},
				},
				Op: token.Token{Kind: token.BangEq},
			},
			want: types.True,
		},
		{
			name: "not equal strings false",
			node: ast.BinaryExpression{
				Left: ast.String{
					Value: "yes",
					Tok:   token.Token{Kind: token.String},
				},
				Right: ast.String{
					Value: "yes",
					Tok:   token.Token{Kind: token.String},
				},
				Op: token.Token{Kind: token.BangEq},
			},
			want: types.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interp := interpreter.New()
			got, err := interp.Eval(tt.node)
			test.Ok(t, err)

			test.Equal(t, got, tt.want, test.Context("eval.Eval(T%) mismatch", tt.node))
		})
	}
}
