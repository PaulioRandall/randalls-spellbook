package wizzard

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

var (
	// ErrForTable occurs when an operation on a known table
	// fails.
	ErrForTable = sin.Template(
		"Regarding table '%s'",
	)

	// ErrForType occurs when an operation on a known model
	// fails.
	ErrForType = sin.Template(
		"Regarding model '%s'",
	)

	// ErrForField occurs when an operation on a known
	// property within a model fails.
	ErrForField = sin.Template(
		"Regarding property '%s'",
	)

	// ErrNotStruct occurs when attempting the parse a
	// non-struct type as a model.
	ErrNotStruct = sin.Template(
		"Model must be a struct, got '%s'",
	)

	// ErrUnsupportedType occurs when attempting to parse
	// a struct into a model but one of the struct's fields
	// has an unsupported type.
	ErrUnsupportedType = sin.Template(
		"Unsupported Go type used for exported field '%s'",
	)

	// ErrWrongObjectType occurs when passing an object with
	// a type that is not the model's GoType.
	ErrWrongObjectType = sin.Template(
		"This model does not handle the passed object type",
	)

	// ErrWrongParameterType occurs when using a type that is
	// not the model's GoType with one of the model's generic
	// methods.
	ErrWrongParameterType = sin.Template(
		"This model does not handle the parameter type used (generics)",
	)

	// ErrTypeMismatch occurs when mapping a struct type to
	// an existing table but a field has a type not
	// compatible with its corresponding table column
	// type.
	ErrTypeMismatch = sin.Template(
		"Model's field type '%s' is not compatible with existing column type '%s'",
	)

	// ErrMissingIdField occurs when attempting to insert,
	// update, or upsert a row but the model used does not
	// contain the required ID field.
	ErrMissingIdField = sin.Err(
		"Model missing ID field",
	)

	// ErrRowScan is returned when an error occurs scanning
	// database rows.
	ErrRowScan = sin.Template(
		"When scanning row '%d'",
	)
)
