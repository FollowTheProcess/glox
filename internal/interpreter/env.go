package interpreter

import (
	"fmt"

	"github.com/FollowTheProcess/glox/internal/syntax/types"
)

// Environment is the environment/scope of the Lox interpreter.
type Environment struct {
	parent *Environment
	values map[string]types.Type
	name   string
}

// NewEnvironment returns a new environment, with an optional parent.
//
// To define the global environment use NewEnvironment(nil).
func NewEnvironment(name string, parent *Environment) *Environment {
	return &Environment{
		name:   name,
		parent: parent,
		values: make(map[string]types.Type),
	}
}

// Define defines a variable in the current scope, if the variable is already defined
// an error is returned.
func (e *Environment) Define(name string, variable types.Type) error {
	if _, exists := e.values[name]; exists {
		return fmt.Errorf("variable %q already defined in this scope (%s): %s", name, e.name, variable)
	}

	e.values[name] = variable
	return nil
}

// Get looks up a name in the environment, if it's not found it traverses the tree of enclosing
// scopes until either the variable is found, or we reach the outermost scope.
func (e *Environment) Get(name string) (v types.Type, ok bool) {
	if v, exists := e.values[name]; exists {
		return v, true
	}

	if e.parent != nil {
		return e.parent.Get(name)
	}

	return nil, false
}

// Assign assigns a type to an already existing variable, if the variable is not found
// in the inner scope, it traverses up the scopes to look for it. If it is not found
// anywhere, an error is returned.
func (e *Environment) Assign(name string, variable types.Type) error {
	if _, exists := e.values[name]; exists {
		e.values[name] = variable
		return nil
	}

	if e.parent != nil {
		return e.parent.Assign(name, variable)
	}

	return fmt.Errorf("assignment to undefined variable %q", name)
}
