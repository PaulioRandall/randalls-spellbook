package storm

import (
	"fmt"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

var (
	ErrCreatingTable        = curse.Err("Failed to create table")
	ErrDroppingTable        = curse.Err("Failed to drop table")
	ErrInsertingObject      = curse.Err("Failed to insert object")
	ErrUpdatingObject       = curse.Err("Failed to update object")
	ErrDeletingObject       = curse.Err("Failed to delete object")
	ErrSelectingAllObjects  = curse.Err("Failed to select all objects")
	ErrSelectingObjectsById = curse.Err("Failed to select objects by ID")

	ErrExeSqlCreate = curse.Err("Executing CREATE TABLE")
	ErrExeSqlDrop   = curse.Err("Executing DROP TABLE")
	ErrExeSqlInsert = curse.Err("Executing INSERT INTO")
	ErrExeSqlUpdate = curse.Err("Executing UPDATE")
	ErrExeSqlDelete = curse.Err("Executing DELETE")
	ErrExeSqlSelect = curse.Err("Executing SELECT")

	ErrScanningRows = curse.Err("Scanning selected rows")
	ErrScanningRow  = curse.Err("Scanning row")

	ErrObjectNotFound = curse.Err("Object not found")

	// ErrDatabaseFile is returned when an error occurs with
	// or while opening or closing the database.
	ErrDatabaseFile = curse.Err("Database IO error")

	// ErrNotStruct is returned when attempting to use a
	// model type with a non-struct kind.
	ErrNotStruct = curse.Err("Model must be a struct")

	// ErrBadFieldKind is returned when a model's type
	// contains an unsupported kind for one of its exported
	// fields.
	ErrBadFieldKind = curse.Err(
		"Model struct has unsupported field kind",
	)

	// ErrNoExportedFields is returned when a model's type
	// has no exported fields. Every table must have at
	// least one column.
	ErrNoExportedFields = curse.Err(
		"Model struct must have at least one exported field",
	)

	// ErrNoSqlType is returned when there is no mapping from
	// a Go type to an database specific Sql type..
	ErrNoSqlType = curse.Err(
		"No type mapping from Go to SQL",
	)

	// ErrTooFewRows is returned when passing a non-positive
	// number to SqlGenerator functions that can generate
	// queries working with multiple records, e.g. INSERT.
	ErrTooFewRows = curse.Err(
		"Not enough rows to generate insert",
	)

	// ErrNoSuchTable is returned when an object is passed
	// to a function which does not have a registered table
	// for its type.
	ErrNoSuchTable = curse.Err(
		"No matching table for object type",
	)

	// ErrBadIdType is returned when an ID passed to a
	// function, e.g. SelectById, is not of the same type as
	// the ID field of the associated model type. This may
	// be returned even for compatible types like int when
	// int64 is expected.
	ErrBadIdType = curse.Err(
		"ID type must match model's ID field type",
	)
)

func tableInfo(table string) string {
	return "Table/Struct: " + table
}

func columnInfo(column string) string {
	return "Column/Field: " + column
}

func idInfo(id any) string {
	return fmt.Sprintf("Object ID: %v", id)
}

func rowInfo(idx int) string {
	return fmt.Sprintf("Row index: %d", idx)
}
