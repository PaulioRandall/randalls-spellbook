package wizzard

import (
	"database/sql"
	ref "reflect"
	"strings"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// Create creates the table the model represents within the
// passed database. It assumes the database is open. If a
// table with the same name already exists then nothing
// happens and no error is returned.
func Create(m Model, db *sql.DB) error {
	query := nidoking.Given(`
		CREATE TABLE IF NOT EXISTS {{table.SqlName}} (
			{{cols.SqlName}} {{cols.SqlType}} NOT NULL DEFAULT {{cols.SqlDefault}},
		  PRIMARY KEY ({{pk_col.SqlName}})
		)
	`).
		InlineMap("table", m).
		ListMap("cols", "", m.Props...).
		InlineMap("pk_col", m.KeyProp()).
		String()

	_, e := db.Exec(query)
	if e != nil {
		return ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Upsert inserts the object if it doesn't exist within the
// database, else it updates the row. It assumes the
// database is open and the table exists. An error is
// returned if the object's type does not match the model's
// GoType.
//
// Inserting via an object that only represents part of a
// table will cause the unset columns to default to
// their zero value. When updating, the key property is
// used to target the row but is not updated.
func Upsert(m Model, db *sql.DB, object any) error {
	if !m.Represents(object) {
		return ErrForType.
			Fmt(m.GoName).
			Wrap(ErrWrongObjectType)
	}

	pkCol := m.KeyProp()
	nonPkCols := m.NonKeyProps()

	if pkCol == (Property{}) {
		return ErrForType.
			Fmt(m.GoName).
			Wrap(ErrMissingIdField)
	}

	query := nidoking.Given(`
		INSERT INTO {{table.SqlName}} (
		  {{col.SqlName}}
		)
		VALUES (
			{{q_marks}}
		)
		ON CONFLICT (
			{{pk_col.SqlName}}
		)
		DO UPDATE SET
			{{non_pk_col.SqlName}} = excluded.{{non_pk_col.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		InlineMap("table", m).
		ListMap("col", ",", m.Props...).
		ListRepeat("q_marks", ",", "?", len(m.Props)).
		InlineMap("pk_col", pkCol).
		ListMap("non_pk_col", ",", nonPkCols...).
		String()

	// Because PK is in the WHERE clause
	cols := append(m.Props, pkCol)
	values := extractFieldValues(cols, object)

	_, e := db.Exec(query, values...)
	if e != nil {
		return ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Insert inserts the object into the database. It assumes
// the database is open and the table exists. An error is
// returned if the object's type does not match the model's
// GoType. Inserting via an object that only represents
// part of a table will cause the unset columns to default
// to their zero value.
func Insert(m Model, db *sql.DB, object any) error {
	if !m.Represents(object) {
		return ErrForType.
			Fmt(m.GoName).
			Wrap(ErrWrongObjectType)
	}

	if m.KeyProp() == (Property{}) {
		return ErrForType.
			Fmt(m.GoName).
			Wrap(ErrMissingIdField)
	}

	query := nidoking.Given(`
		INSERT INTO {{table.SqlName}} (
		  {{cols.SqlName}}
		)
		VALUES (
			{{q_marks}}
		)
	`).
		InlineMap("table", m).
		ListMap("cols", ",", m.Props...).
		ListRepeat("q_marks", ",", "?", len(m.Props)).
		String()

	values := extractFieldValues(m.Props, object)
	_, e := db.Exec(query, values...)
	if e != nil {
		return ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Update updates the object within the database. It
// assumes the database is open and the table exists. An
// error is returned if the object's type does not match
// the model's GoType. The key property is used to target
// the row but is not updated.
func Update(m Model, db *sql.DB, object any) error {
	if !m.Represents(object) {
		return ErrForType.
			Fmt(m.GoName).
			Wrap(ErrWrongObjectType)
	}

	pkCol := m.KeyProp()
	nonPkCols := m.NonKeyProps()

	if pkCol == (Property{}) {
		return ErrForType.
			Fmt(m.GoName).
			Wrap(ErrMissingIdField)
	}

	query := nidoking.Given(`
		UPDATE
			{{table.SqlName}}
		SET
			{{non_pk_col.SqlName}} = ?
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		InlineMap("table", m).
		ListMap("non_pk_col", ",", nonPkCols...).
		InlineMap("pk_col", pkCol).
		String()

	// Because PK is in the WHERE clause
	cols := append(nonPkCols, pkCol)
	values := extractFieldValues(cols, object)

	_, e := db.Exec(query, values...)
	if e != nil {
		return ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Select returns all rows matching the passed where clause
// from the table the model represents within passed
// database. It assumes the database is open and the table
// exists.
func Select[T any](m Model, db *sql.DB, where string, args ...any) ([]T, error) {
	var o T
	if !m.Represents(o) {
		return nil, ErrForType.
			Fmt(m.GoName).
			Wrap(ErrWrongParameterType)
	}

	where = strings.TrimSpace(where)
	if where != "" {
		where = "WHERE " + where
	}

	query := nidoking.Given(`
		SELECT
			{{col.SqlName}}
		FROM
			{{table.SqlName}}
		{{where}}
	`).
		InlineMap("table", m).
		ListMap("col", ",", m.Props...).
		Inline("where", where).
		String()

	rows, e := db.Query(query, args...)
	if e != nil {
		return nil, ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return ScanRows[T](m, rows)
}

// SelectFirst returns the first row matching the passed
// where clause from the table the model represents within
// the passed database. It assumes the database is open and
// the table exists. The first return value will contain
// the row data and second will be true if the row exists,
// else zero value and false are returned.
func SelectFirst[T any](m Model, db *sql.DB, where string, args ...any) (T, bool, error) {
	var zero T

	if !m.Represents(zero) {
		return zero, false, ErrForType.
			Fmt(m.GoName).
			Wrap(ErrWrongParameterType)
	}

	where = strings.TrimSpace(where)
	if where != "" {
		where = "WHERE " + where
	}

	query := nidoking.Given(`
		SELECT
			{{cols.SqlName}}
		FROM
			{{table.SqlName}}
		{{where}}
		LIMIT 1
	`).
		ListMap("cols", ",", m.Props...).
		InlineMap("table", m).
		InlineMap("pk_col", m.KeyProp()).
		Inline("where", where).
		String()

	rows, e := db.Query(query, args...)
	if e != nil {
		return zero, false, ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return ScanFirstRow[T](m, rows)
}

// Delete removes all table rows matching the given where
// clause and arguments from the table the model represents
// within the passed database. It assumes the database is
// open and the table exists. If no matching rows are found
// then nothing happens and no error is returned.
func Delete(m Model, db *sql.DB, where string, args ...any) error {
	where = strings.TrimSpace(where)
	if where != "" {
		where = "WHERE " + where
	}

	query := nidoking.Given(`
		DELETE FROM
			{{table.SqlName}}
		{{where}}
	`).
		InlineMap("table", m).
		InlineMap("pk_col", m.KeyProp()).
		Inline("where", where).
		String()

	_, e := db.Exec(query, args...)
	if e != nil {
		return ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Drop removes the table the model represents from the
// passed database. It assumes the database is open. If the
// table doesn't exist then nothing happens and no error is
// returned.
func Drop(m Model, db *sql.DB) error {
	query := nidoking.Given(`
		DROP TABLE IF EXISTS {{table.SqlName}}
	`).
		InlineMap("table", m).
		String()

	_, e := db.Exec(query)
	if e != nil {
		return ErrForType.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// ScanRows scans all SQL rows returned from a query into
// objects of type T. If T's type does not match the
// model's GoType then an error is returned. The rows
// object will then be closed preventing further scanning.
func ScanRows[T any](m Model, rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	var zero T
	if !m.Represents(zero) {
		return nil, ErrWrongParameterType
	}

	values, valuePtrs := newValueContainers(m)
	var results []T

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrRowScan.Fmt(i).Wrap(e)
		}

		var item T
		populate(m, &item, values)
		results = append(results, item)
	}

	e := rows.Err()
	if e != nil {
		return nil, e
	}

	return results, nil
}

// ScanFirstRow scans the first result within the passed
// SQL rows into into an object of type T. If there are no
// more rows then the zero value of T is retuned. If T's
// type does not match the model's GoType then an error is
// returned. The rows object will then be closed preventing
// further scanning.
func ScanFirstRow[T any](m Model, rows *sql.Rows) (T, bool, error) {
	defer rows.Close()

	var zero T
	if !m.Represents(zero) {
		return zero, false, ErrWrongParameterType
	}

	if !rows.Next() {
		return zero, false, nil
	}

	values, valuePtrs := newValueContainers(m)
	e := rows.Scan(valuePtrs...)
	if e != nil {
		return zero, false, ErrRowScan.Fmt(0).Wrap(e)
	}

	e = rows.Err()
	if e != nil {
		return zero, false, e
	}

	var item T
	populate(m, &item, values)
	return item, true, nil
}

func extractFieldValues(
	props []Property,
	object any,
) []any {
	value := ref.ValueOf(object)
	result := make([]any, len(props))

	for i := 0; i < len(props); i++ {
		field := value.FieldByName(props[i].GoName)
		result[i] = field.Interface()
	}

	return result
}

func newValueContainers(m Model) ([]any, []any) {
	colCount := len(m.Props)
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

	for i, col := range m.Props {
		values[i] = col.New[any]()
		valuePtrs[i] = &values[i]
	}

	return values, valuePtrs
}

func populate[T any](m Model, item *T, values []any) {
	objVal := ref.ValueOf(item).Elem()

	for valueIdx, col := range m.Props {
		v := ref.ValueOf(values[valueIdx])
		fieldVal := objVal.Field(col.GoIndex)

		if v.CanConvert(fieldVal.Type()) {
			v = v.Convert(fieldVal.Type())
		} else {
			// TODO: Return as error
			panic("Can't convert from " + v.Type().Name() + " to " + fieldVal.Type().Name())
		}

		fieldVal.Set(v)
	}
}
