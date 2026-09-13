package storm

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrCreatingTable        = stormy("Failed to create table")
	ErrDroppingTable        = stormy("Failed to drop table")
	ErrInsertingObject      = stormy("Failed to insert object")
	ErrUpdatingObject       = stormy("Failed to update object")
	ErrDeletingObject       = stormy("Failed to delete object")
	ErrSelectingAllObjects  = stormy("Failed to select all objects")
	ErrSelectingObjectsById = stormy("Failed to select objects by ID")

	ErrExeSqlCreate = stormy("Executing CREATE TABLE")
	ErrExeSqlDrop   = stormy("Executing DROP TABLE")
	ErrExeSqlInsert = stormy("Executing INSERT INTO")
	ErrExeSqlUpdate = stormy("Executing UPDATE")
	ErrExeSqlDelete = stormy("Executing DELETE")
	ErrExeSqlSelect = stormy("Executing SELECT")

	ErrScanningRows = stormy("Scanning selected rows")
	ErrScanningRow  = stormy("Scanning row")

	ErrObjectNotFound = stormy("Object not found")

	// ErrNotStruct is returned when attempting to use a
	// model type with a non-struct kind.
	ErrNotStruct = stormy("Model must be a struct")

	// ErrBadFieldKind is returned when a model's type
	// contains an unsupported kind for one of its exported
	// fields.
	ErrBadFieldKind = stormy(
		"Model struct has unsupported field kind",
	)

	// ErrNoExportedFields is returned when a model's type
	// has no exported fields. Every table must have at
	// least one column.
	ErrNoExportedFields = stormy(
		"Model struct must have at least one exported field",
	)

	// ErrNoSqlType is returned when there is no mapping from
	// a Go type to an database specific Sql type..
	ErrNoSqlType = stormy(
		"No type mapping from Go to SQL",
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
)

// StormError is returned for all errors. It implements
// functions needed to work with the standard errors
// package.
type StormError struct {
	errId  string
	msg    string
	cause  error
	table  string
	column string
	id     any
	row    int
}

func stormy(msg string, args ...any) StormError {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return StormError{
		errId: uuid.New().String(),
		msg:   msg,
		row:   -1,
	}
}

// Wrap wraps the passed error as the cause of receiver's
// error.
func (se StormError) Wrap(e error) StormError {
	type wrapper interface {
		Unwrap() error
	}

	for se.cause = e; ; {
		seCause, ok := e.(StormError)
		if ok {
			copyUnfilledErrInfo(&se, &seCause)
		}

		wrapped, ok := e.(wrapper)
		if !ok {
			break
		}

		e = wrapped.Unwrap()
	}

	return se
}

func copyUnfilledErrInfo(target, se *StormError) {
	if target.table == "" {
		target.table = se.table
	}

	if target.column == "" {
		target.column = se.column
	}

	if target.id == nil {
		target.id = se.id
	}

	if target.row == -1 {
		target.row = se.row
	}
}

// Unwrap returns the cause of the error, i.e. the wrapped
// error.
func (se StormError) Unwrap() error {
	return se.cause
}

// Table sets the name of the table associated  with the
// error.
func (se StormError) Table(table string) StormError {
	se.table = table
	return se
}

// Column sets the name of the column associated with the
// error.
func (se StormError) Column(column string) StormError {
	se.column = column
	return se
}

// Id sets the object ID associated with the error.
func (se StormError) Id(id any) StormError {
	se.id = id
	return se
}

// Row sets the row index associated with the error.
func (se StormError) Row(row int) StormError {
	se.row = row
	return se
}

// Error returns the full error message. It satisfies Go's
// error interface.
func (se StormError) Error() string {
	s := se.msg

	if se.row > -1 {
		s += newErrLine("Row index: %d", se.row)
	}

	if se.table != "" {
		s += newErrLine("Table/struct: %s", se.table)
	}

	if se.column != "" {
		s += newErrLine("Column/field: %s", se.column)
	}

	if se.id != nil {
		s += newErrLine("Object ID: %v", se.id)
	}

	if se.cause != nil {
		s += ":\n[Cause] " + se.cause.Error()
	}

	return s
}

// Is returns true if target is a StormError and has the
// same errId as the receiving StormError.
func (se StormError) Is(target error) bool {
	if se2, ok := target.(StormError); ok {
		return se2.errId == se.errId
	}
	return false
}

func newErrLine(msg string, args ...any) string {
	return ":\n\t+ " + fmt.Sprintf(msg, args...)
}

var _ error = StormError{}
