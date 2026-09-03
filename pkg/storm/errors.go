package storm

import (
	"errors"
)

var (
	// ErrNotStruct is returned when attempting to use a
	// model type with a non-struct kind.
	ErrNotStruct = errors.New("Model must be a struct")

	// ErrNotPublic is returned when a model type is not
	// public (exported).
	ErrNotPublic = errors.New("Model struct must be public")

	// ErrBadFieldKind is returned when a model's type
	// contains an unsupported kind for one of its exported
	// fields.
	ErrBadFieldKind = errors.New(
		"Model struct has unsupported field kind",
	)

	// ErrMissFields is returned when a model's type has no/
	// public (exported) fields. Every table must have at
	// least one column.
	ErrMissFields = errors.New(
		"Model struct must have at least one exported field",
	)

	// ErrNoSqlType is returned when there is no mapping from
	// a Go type to an database specific Sql type, e.g, for
	// SQLite 'int => INTEGER' but there is no mapping for
	// 'error => _'.
	ErrNoSqlType = errors.New(
		"No mapping for Go to SQL type",
	)

	// ErrTooFewRows is returned when passing a non-positive
	// number to SqlGenerator functions that can generate
	// queries working with multiple records, e.g. INSERT.
	ErrTooFewRows = errors.New(
		"Not enough rows to generate insert",
	)

	// ErrNoTableForType is returned when an object is passed
	// to a function which does not have a registered table
	// for its type.
	ErrNoTableForType = errors.New(
		"No matching table for object type",
	)

	// ErrBadIdType is returned when an ID passed to a
	// function, e.g. SelectById, is not of the same type as
	// the ID field of the associated model type. This may
	// be returned even for compatible types like int when
	// int64 is expected.
	ErrBadIdType = errors.New(
		"ID type must match model's ID field type",
	)

	// ErrRecordNotFound is returned when an record could not
	// be found with a given ID, e.g. SelectById.
	ErrRecordNotFound = errors.New("Record not found")
)
