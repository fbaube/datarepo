package utils

import (
	"time"
)

// These are in app model packages but not here I guess.
type DAO struct{}

// Row is a type constraint for types that
//   - represent a single database row, and
//   - therefore correspond to a Go struct.
//
// According to this definition, a Row can
//   - cough up a set of ptrs suitable for
//     setting every field/column, and
//   - have its own address be taken (and
//     passed around), making it writable.
type Row[T any] interface {
	// PtrFields returns a slice of
	// ptrs, one per struct field,
	// for use with [Row.Scan].
	// Implement it for - and call
	// it with - ptr receivers ONLY.
	PtrFields() []any // TODO Try []*any
	// *T means we can create instances of type T as
	// ptrs. This is necessary to avoid call-by-value,
	// thus permitting the actual return of values. See:
	// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#pointer-method-example.
	*T
}

// Field is a type constraint for types that represent
//   - a single database column, or
//   - any other non-struct datum,
//     such as an int64 primary index.
type Field interface {
	~byte | ~int16 | ~int32 | ~int64 | ~float64 |
		~string | ~bool | time.Time
}

// Sorry, no dice. Compiler barfs:
// ../db/access.go:38:11: cannot use type Field outside
//  a type constraint: interface contains type constraints
//
// type Ptr *Field
// The kind of thing that works:
// func Nil[T ~*E, E any](ptr T) bool {
// So try
/* OOPS
type Ptr interface {
	/* T * / ~*E
	// E any
}
* /
type FieldPtr interface {
	~*Field
}
*/
// Maybe a Ptr should be a union of ptrs to all the basic
// Go datatypes that are mapped to SQL field types.
// In other words, more or less, just duplicate the list
// defined for "Field", but with asterisks.
