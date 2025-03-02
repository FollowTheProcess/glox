package types_test

import (
	"testing"
	"testing/quick"

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

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b types.Type // Types under test
		name string     // Name of the test case
		want bool       // Expected return value
	}{
		{
			name: "strings equal",
			a:    types.String{Value: "yes"},
			b:    types.String{Value: "yes"},
			want: true,
		},
		{
			name: "strings not equal",
			a:    types.String{Value: "yes"},
			b:    types.String{Value: "no"},
			want: false,
		},
		{
			name: "bool equal",
			a:    types.True,
			b:    types.True,
			want: true,
		},
		{
			name: "bool not equal",
			a:    types.True,
			b:    types.False,
			want: false,
		},
		{
			name: "number equal",
			a:    types.Number{Value: 42},
			b:    types.Number{Value: 42},
			want: true,
		},
		{
			name: "number not equal",
			a:    types.Number{Value: 42},
			b:    types.Number{Value: 3.14159},
			want: false,
		},
		{
			name: "different types not equal",
			a:    types.String{Value: "42"},
			b:    types.Number{Value: 42},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := types.Equal(tt.a, tt.b)

			test.Equal(t, got, tt.want, test.Context("types.Equal[a %T, b %T](%s, %s) mismatch", tt.a, tt.b, tt.a, tt.b))
		})
	}
}

// TODO(@FollowTheProcess): Turn this into Fuzz? Or maybe implement [quick.Generator]
// for each type? That could be fun

func TestProperties(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		f := func(a, b types.String) bool {
			// Property: Two types with different kinds are never equal
			if a.Kind() != b.Kind() {
				eq := types.Equal(a, b)
				if eq {
					return false // Fail test by returning false
				}
			}

			// Property: If String() is different, equal must return false
			if a.String() != b.String() {
				eq := types.Equal(a, b)
				if eq {
					return false // Fail test by returning false
				}
			}

			return true
		}
		if err := quick.Check(f, nil); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("number", func(t *testing.T) {
		f := func(a, b types.Number) bool {
			// Property: Two types with different kinds are never equal
			if a.Kind() != b.Kind() {
				eq := types.Equal(a, b)
				if eq {
					return false // Fail test by returning false
				}
			}

			// Property: If String() is different, equal must return false
			if a.String() != b.String() {
				eq := types.Equal(a, b)
				if eq {
					return false // Fail test by returning false
				}
			}

			return true
		}
		if err := quick.Check(f, nil); err != nil {
			t.Fatal(err)
		}
	})
}
