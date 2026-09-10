package storm

import (
	"fmt"
)

var (
	ErrCreatingTable       = stormy("Failed to create table")
	ErrDroppingTable       = stormy("Failed to drop table")
	ErrInsertingObject     = stormy("Failed to insert object")
	ErrUpdatingObject      = stormy("Failed to update object")
	ErrDeletingObject      = stormy("Failed to delete object")
	ErrSelectingAllObjects = stormy(
		"Failed to select all objects",
	)
	ErrSelectingObjectsById = stormy(
		"Failed to select objects by ID",
	)

	ErrExeSqlCreate = stormy("Executing CREATE TABLE")
	ErrExeSqlDrop   = stormy("Executing DROP TABLE")
	ErrExeSqlInsert = stormy("Executing INSERT INTO")
	ErrExeSqlUpdate = stormy("Executing UPDATE")
	ErrExeSqlDelete = stormy("Executing DELETE")
	ErrExeSqlSelect = stormy("Executing SELECT")

	ErrScanningRows = stormy("Scanning selected rows")
	ErrScanningRow  = stormy("Scanning row")

	ErrInitObject = stormy("Failed to intiialise object")

	// ErrNotStruct is returned when attempting to use a
	// model type with a non-struct kind.
	ErrNotStruct = stormy("Model must be a struct")

	// ErrBadFieldKind is returned when a model's type
	// contains an unsupported kind for one of its exported
	// fields.
	ErrBadFieldKind = stormy(
		"Model struct has unsupported field kind",
	)

	// ErrMissFields is returned when a model's type has no/
	// public (exported) fields. Every table must have at
	// least one column.
	ErrMissFields = stormy(
		"Model struct must have at least one exported field",
	)

	// ErrNoSqlType is returned when there is no mapping from
	// a Go type to an database specific Sql type, e.g, for
	// SQLite 'int => INTEGER' but there is no mapping for
	// 'error => _'.
	ErrNoSqlType = stormy(
		"No mapping for Go to SQL type",
	)

	// ErrTooFewRows is returned when passing a non-positive
	// number to SqlGenerator functions that can generate
	// queries working with multiple records, e.g. INSERT.
	ErrTooFewRows = stormy(
		"Not enough rows to generate insert",
	)

	// ErrNoSuchTable is returned when an object is passed
	// to a function which does not have a registered table
	// for its type.
	ErrNoSuchTable = stormy(
		"No matching table for object type",
	)

	// ErrBadIdType is returned when an ID passed to a
	// function, e.g. SelectById, is not of the same type as
	// the ID field of the associated model type. This may
	// be returned even for compatible types like int when
	// int64 is expected.
	ErrBadIdType = stormy(
		"ID type must match model's ID field type",
	)

	// ErrObjectNotFound is returned when an object could not
	// be found with a given ID.
	ErrObjectNotFound = stormy("Object not found")
)

// StormError is returned for all errors. It implements
// functions needed to work with the standard errors
// package.
type StormError struct {
	original error
	wrapped  error
	tableMsg string
	idMsg    string
	rowMsg   string
	msg      string
}

func stormy(msg string, args ...any) StormError {
	se := StormError{
		msg: msg,
	}

	if len(args) > 0 {
		se.msg = "[Stormy] " + fmt.Sprintf(msg, args...)
	}

	se.original = se
	return se
}

func (se StormError) Wrap(e error) StormError {
	se.wrapped = e
	return se
}

func (se StormError) Table(tableName string) StormError {
	se.tableMsg = newErrLine("For table/struct: %s", tableName)
	return se
}

func (se StormError) Id(id any) StormError {
	se.idMsg = newErrLine("With ID: %s", id)
	return se
}

func (se StormError) Row(rowIdx int) StormError {
	se.rowMsg = newErrLine("At row index: %d", rowIdx)
	return se
}

func (se StormError) Error() string {
	s := se.msg +
		se.rowMsg +
		se.tableMsg +
		se.idMsg

	if se.wrapped != nil {
		s += ":\n" + se.wrapped.Error()
	}

	return s
}

func (se StormError) Is(target error) bool {
	if se2, ok := target.(StormError); ok {
		return se2.original == se.original
	}
	return false
}

func (se StormError) Unwrap() error {
	return se.wrapped
}

func newErrLine(msg string, args ...any) string {
	return ":\n\t+ " + fmt.Sprintf(msg, args...)
}

var _ error = StormError{}
