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
	KindString
)

// TODO(@FollowTheProcess): Could we use generics here? Have a Value() T method in the interface?

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

	String struct {
		Value string
	}
)

// [Type] implementations

func (n Number) Kind() Kind { return KindNumber }
func (b Bool) Kind() Kind   { return KindBool }
func (s String) Kind() Kind { return KindString }

func (n Number) String() string { return strconv.FormatFloat(n.Value, 'g', -1, 64) }
func (b Bool) String() string   { return strconv.FormatBool(b.Value) }
func (s String) String() string { return strconv.Quote(s.Value) }

// IsTruthy reports whether the type is considered truthy or falsey.
func IsTruthy(t Type) bool {
	switch typ := t.(type) {
	case *Bool:
		return typ.Value
	case Number:
		return typ.Value != 0
	case String:
		// Like python, an empty string is falsey, note this differs from Lox in the book
		// but I don't like the implementation there where only nil and false are falsey
		return len(typ.Value) != 0
	default:
		// TODO(@FollowTheProcess): Don't love the panic, but it's a good way of making sure
		// I don't miss one while developing. Remove this when all the types are handled and just
		// return false
		panic(fmt.Sprintf("Unhandled type in IsTruthy: %T", t))
	}
}

// Equal reports whether two types should be considered equal.
func Equal(a, b Type) bool {
	// If one is nil and the other isn't, not equal
	if (a == nil) != (b == nil) {
		return false
	}

	// Now they are either both nil, or both actual types

	if a == nil && b == nil {
		return true // Two nils are equal, like python None
	}

	// By here they *must* be non-nil

	// Two types of different kinds are *never* equal
	if a.Kind() != b.Kind() {
		return false
	}

	// TODO(@FollowTheProcess): Not sure if we can get away with this long term
	// but just comparing the string representations seems fine enough? We know
	// by this point that they are the same kind, and String() simply formats
	// the underlying value in a kind-specific way so... seems legit
	return a.String() == b.String()
}
