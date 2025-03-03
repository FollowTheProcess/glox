package interpreter_test

import (
	"testing"

	"github.com/FollowTheProcess/glox/internal/interpreter"
	"github.com/FollowTheProcess/glox/internal/syntax/types"
	"github.com/FollowTheProcess/test"
)

func TestGet(t *testing.T) {
	tests := []struct {
		want   types.Type            // Expected return value
		parent map[string]types.Type // Create parent environment with these variables
		values map[string]types.Type // Create test environment with these variables
		name   string                // Name of the test case
		get    string                // Variable name to lookup
		ok     bool                  // Expected ok return value
	}{
		{
			name:   "empty string",
			parent: nil,
			values: nil,
			get:    "",
			want:   nil,
			ok:     false,
		},
		{
			name:   "no variables",
			parent: nil,
			values: nil,
			get:    "something",
			want:   nil,
			ok:     false,
		},
		{
			name:   "bad ident whitespace",
			parent: nil,
			values: nil,
			get:    "not an ident",
			want:   nil,
			ok:     false,
		},
		{
			name:   "bad ident first char",
			parent: nil,
			values: nil,
			get:    "1variableplz",
			want:   nil,
			ok:     false,
		},
		{
			name: "global variable",
			parent: map[string]types.Type{
				"is_global": types.True,
			},
			values: nil,
			get:    "is_global",
			want:   types.True,
			ok:     true,
		},
		{
			name:   "local variable",
			parent: nil,
			values: map[string]types.Type{
				"is_local": types.True,
			},
			get:  "is_local",
			want: types.True,
			ok:   true,
		},
		{
			name: "local preferred over global",
			parent: map[string]types.Type{
				"meaning_of_life": types.Number{Value: 42},
			},
			values: map[string]types.Type{
				"meaning_of_life": types.String{Value: "who knows"},
			},
			get:  "meaning_of_life",
			want: types.String{Value: "who knows"},
			ok:   true,
		},
		{
			name:   "undefined",
			parent: nil,
			values: map[string]types.Type{
				"is_local": types.True,
			},
			get:  "something_else",
			want: nil,
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parent *interpreter.Environment
			if len(tt.parent) != 0 {
				parent = interpreter.NewEnvironment("parent", nil)
				for name, value := range tt.parent {
					err := parent.Define(name, value)
					test.Ok(t, err, test.Context("could not define variable in parent env"))
				}
			}

			env := interpreter.NewEnvironment("test", parent)
			if len(tt.values) != 0 {
				for name, value := range tt.values {
					err := env.Define(name, value)
					test.Ok(t, err, test.Context("could not define variable in test env"))
				}
			}

			got, ok := env.Get(tt.get)
			test.Equal(t, ok, tt.ok, test.Context("Get(%s) ok did not match expected", tt.get))
			test.EqualFunc(t, got, tt.want, types.Equal, test.Context("Get(%s) value did not match expected", tt.get))
		})
	}
}
