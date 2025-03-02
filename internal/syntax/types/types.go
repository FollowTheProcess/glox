// Package type implements the type/object system for glox.
//
// The original Java implementation fell back to the Java Object class, we obviously
// can't do that in Go as there is no Object, so we make our own!
package types

import (
	"fmt"
	"strconv"
)

// Kind is the kind of Lox type e.g. bool, number, string etc.
type Kind int

const (
	KindNumber Kind = iota
	KindBool
)

// Type is a Lox object type.
type Type interface {
	fmt.Stringer

	// Kind returns the kind of the Type.
	Kind() Kind
}

// Only one canonical true and false.
var (
	True  = &Bool{Value: true}
	False = &Bool{Value: false}
)

// Concrete types.
type (
	Number struct {
		Value float64
	}

	Bool struct {
		Value bool
	}
)

// [Type] implementations

func (n Number) Kind() Kind { return KindNumber }
func (b Bool) Kind() Kind   { return KindBool }

func (n Number) String() string { return strconv.FormatFloat(n.Value, 'g', -1, 64) }
func (b Bool) String() string   { return strconv.FormatBool(b.Value) }

// IsTruthy reports whether the type is considered truthy or falsey.
func IsTruthy(t Type) bool {
	switch typ := t.(type) {
	case *Bool:
		return typ.Value
	case Number:
		return typ.Value != 0
	default:
		// TODO(@FollowTheProcess): Don't love the panic, but it's a good way of making sure
		// I don't miss one while developing. Remove this when all the types are handled and just
		// return false
		panic(fmt.Sprintf("Unhandled type in IsTruthy: %T", t))
	}
}
