package eval_test

import (
	"testing"

	"github.com/FollowTheProcess/glox/internal/eval"
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eval.Eval(tt.node)
			test.Ok(t, err)

			test.Equal(t, got, tt.want)
		})
	}
}
