package modtab

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

var (
	// ErrForModel occurs when an operation on a known model
	// or table fails.
	ErrForModel = sin.Template(
		"Regarding model '%s'",
	)

	// ErrForField occurs when an operation on a known field
	// within a model fails.
	ErrForField = sin.Template(
		"Regarding field '%s'",
	)

	// ErrNotStruct occurs when a non-struct type is used as
	// a model.
	ErrNotStruct = sin.Template(
		"Model must be a struct kind, got kind '%s'",
	)

	// ErrUnsupportedType occurs when the Go kind for model
	// field's type is not supported.
	ErrUnsupportedType = sin.Template(
		"Unsupported Go kind used for exported field '%s'",
	)

	// ErrTypeMismatch occurs when mapping a struct type to
	// an existing table but one of the field's has a type
	// not compatible with its corrisponding table column
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
