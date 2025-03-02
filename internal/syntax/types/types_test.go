package types_test

import (
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/types"
	"github.com/FollowTheProcess/test"
)

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		typ  types.Type // The type under test
		name string     // Name of the test case
		want bool       // Expected return value
	}{
		{
			name: "True",
			typ:  types.True,
			want: true,
		},
		{
			name: "False",
			typ:  types.False,
			want: false,
		},
		{
			name: "zero",
			typ:  types.Number{Value: 0},
			want: false,
		},
		{
			name: "non zero",
			typ:  types.Number{Value: 42},
			want: true,
		},
		{
			name: "empty string",
			typ:  types.String{Value: ""},
			want: false,
		},
		{
			name: "non empty string",
			typ:  types.String{Value: "stuff"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, types.IsTruthy(tt.typ), tt.want, test.Context("IsTruthy(%T) mismatch", tt.typ))
		})
	}
}
