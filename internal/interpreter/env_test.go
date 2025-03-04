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

func TestDefine(t *testing.T) {
	tests := []struct {
		value   types.Type            // Variable value to define
		values  map[string]types.Type // Create test environment with these variables
		name    string                // Name of the test case
		define  string                // Variable name to define
		errMsg  string                // The expected error message, if any
		wantErr bool                  // Whether Define should return an error
	}{
		{
			name: "already exists",
			values: map[string]types.Type{
				"something": types.True,
			},
			define:  "something",
			value:   types.False,
			wantErr: true,
			errMsg:  `variable "something" already defined in this scope (test): false`,
		},
		{
			name:    "valid with nil map",
			values:  nil,
			define:  "something",
			value:   types.True,
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := interpreter.NewEnvironment("test", nil) // Define never looks at parent scope
			if len(tt.values) != 0 {
				for name, value := range tt.values {
					err := env.Define(name, value)
					test.Ok(t, err, test.Context("could not pre-define variable %s", name))
				}
			}

			err := env.Define(tt.define, tt.value)
			test.WantErr(t, err, tt.wantErr)
			if err != nil {
				test.Equal(t, err.Error(), tt.errMsg)
			} else {
				// If there wasn't an error inserting it, we should be able to get
				// it back again with no problem
				got, ok := env.Get(tt.define)
				test.True(t, ok)
				test.EqualFunc(t, got, tt.value, types.Equal)
			}
		})
	}
}

func TestAssign(t *testing.T) {
	tests := []struct {
		value   types.Type            // Variable to assign
		parent  map[string]types.Type // Create parent environment with these variables
		values  map[string]types.Type // Create test environment with these variables
		name    string                // Name of the test case
		assign  string                // Variable name to assign
		errMsg  string                // If wantErr, what should the message say
		wantErr bool                  // Whether assign should return an error
	}{
		{
			name:   "undefined",
			parent: nil,
			values: map[string]types.Type{
				"a_thing": types.String{Value: "yes"},
			},
			assign:  "something_else",
			value:   types.String{Value: "hello"},
			wantErr: true,
			errMsg:  `assignment to undefined variable "something_else"`,
		},
		{
			name:   "exists in local scope",
			parent: nil,
			values: map[string]types.Type{
				"something": types.True,
			},
			assign:  "something",
			value:   types.False,
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "exists in parent scope",
			parent: map[string]types.Type{
				"pi": types.Number{Value: 3.14159},
			},
			values:  nil,
			assign:  "pi",
			value:   types.Number{Value: 3}, // A crime!
			wantErr: false,
			errMsg:  "",
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

			err := env.Assign(tt.assign, tt.value)
			test.WantErr(t, err, tt.wantErr)
			if err != nil {
				test.Equal(t, err.Error(), tt.errMsg)
			} else {
				// If there wasn't an error assigning it, we should be able to
				// get it back
				got, ok := env.Get(tt.assign)
				test.True(t, ok)
				test.EqualFunc(t, got, tt.value, types.Equal)
			}
		})
	}
}
