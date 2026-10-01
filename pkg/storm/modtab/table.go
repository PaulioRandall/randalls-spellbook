package modtab

import (
	"database/sql"
	ref "reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// ModelTable maps a Go struct, or part of it, to a
// database table structure, or part of it. A ModelTable is
// always derived from a Go type of kind 'relfect.Struct'.
type ModelTable struct {
	// GoType is the struct type as returned by
	// reflect.TypeOf.
	GoType ref.Type

	// GoName is the struct's name, i.e. GoType.Name().
	GoName string

	// SqlName is the name of the table the struct represents
	// within a database. Currently, this is always the same
	// as GoName.
	SqlName string

	// Columns is an ordered list of exported field to column
	// mappings.
	Columns []ModelColumn
}

// ModelColumn maps a Go struct's field to a specific
// database table column.
type ModelColumn struct {
	// GoType is the field type as returned by reflect.TypeOf.
	GoType ref.Type

	// GoIndex is the field index, i.e. it's 0-based position
	// within the struct definition
	GoIndex int

	// GoName is the field's name, i.e. GoType.Name().
	GoName string

	// SqlType is the SQLite type that maps to the GoType.
	// Note that several Go types may map to a single SQLite
	// typee.
	SqlType string

	// SqlName is the name of the column the field represents
	// within a database. Currently, this is always the same
	// as GoName.
	SqlName string

	// SqlDefault is the default value given to the column if
	// a value is not provided during row insertion.
	// Currently, this is always the zero value of the
	// GoType.
	SqlDefault any

	// PrimaryKey is true if the column represents the
	// primary key. By default this is the first column,
	// however, this may differ if adjusting a ModelColumn to
	// align with an existing table column.
	PrimaryKey bool
}

// Parse creates a full [ModelTable] representation of the
// passed model's type. One should not assume the result
// will map directly to a same named table within a
// database; use [Map] to create a representation that is
// usable with a specific database. Type mappings:
//
//	INTEGER:
//		int, int8, int16, int32, int64,
//		uint, uint8, uint16, uint32, uint64
//	REAL:
//		float32, float64
//	TEXT:
//		string
//
// You are responsible for choosing appropriate data types
// for your structures. If you choose to use an unsigned
// integer type (unit etc) or shorter integer type
// (int8 etc) then you are responsible for ensuring or
// managing number polarity and overflow. I recommend
// sticking to int, int64, float64, and string as these are
// the least likely to face issues.
func Parse(model any) (table ModelTable, e error) {
	modelType := derefModelType(model)

	if modelType.Kind() != ref.Struct {
		e = ErrNotStruct.Fmt(modelType.Kind().String())
		goto Err
	}

	table.GoType = modelType
	table.GoName = modelType.Name()
	table.SqlName = modelType.Name()

	table.Columns, e = parseColumns(modelType)
	if e != nil {
		goto Err
	}

	return table, nil

Err:
	return ModelTable{}, ErrForModel.
		Fmt(modelType.Name()).
		Wrap(e)
}

func derefModelType(model any) ref.Type {
	modelType := ref.TypeOf(model)

	for modelType.Kind() == ref.Ptr {
		modelType = modelType.Elem()
	}

	return modelType
}

func parseColumns(modelType ref.Type) ([]ModelColumn, error) {
	isIdField := true

	var cols []ModelColumn

	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)

		if !f.IsExported() {
			continue
		}

		col, e := parseColumn(f, i, isIdField)
		if e != nil {
			return nil, ErrForField.Fmt(f.Name).Wrap(e)
		}

		cols = append(cols, col)
		isIdField = false
	}

	return cols, nil
}

func parseColumn(
	f ref.StructField,
	i int,
	isIdField bool,
) (ModelColumn, error) {
	var col ModelColumn
	var defaultValue any

	if f.Type == ref.TypeOf("") {
		defaultValue = "''"
	} else {
		defaultValue = ref.Zero(f.Type).Interface()
	}

	sqlType, e := mapGoToSqlType(f.Type.Kind())
	if e != nil {
		return ModelColumn{}, e
	}

	col.GoType = f.Type
	col.GoIndex = i
	col.GoName = f.Name
	col.SqlType = sqlType
	col.SqlName = f.Name
	col.SqlDefault = defaultValue
	col.PrimaryKey = isIdField

	return col, nil
}

func mapGoToSqlType(fieldKind ref.Kind) (string, error) {
	switch fieldKind {
	case ref.Int, ref.Int8, ref.Int16, ref.Int32, ref.Int64:
		fallthrough
	case ref.Uint, ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64:
		return "INTEGER", nil
	case ref.Float32, ref.Float64:
		return "REAL", nil
	case ref.String:
		return "TEXT", nil
	default:
		return "", ErrUnsupportedType.Fmt(fieldKind.String())
	}
}

// PkCol returns the column representing the table's
// primary key. Zero value is returned if not found.
func (t ModelTable) PkCol() ModelColumn {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return ModelColumn{}
}

// NonPkCols returns all columns except the primary key
// column.
func (t ModelTable) NonPkCols() []ModelColumn {
	result := make([]ModelColumn, 0, len(t.Columns))

	for _, col := range t.Columns {
		if !col.PrimaryKey {
			result = append(result, col)
		}
	}

	return result
}

// Create creates the table within the passed database. It
// assumes the database is open. If a table with the same
// name already exists then nothing happens and no error is
// returned.
func (t ModelTable) Create(db *sql.DB) error {
	query := nidoking.Given(`
		CREATE TABLE IF NOT EXISTS {{table.SqlName}} (
			{{col.SqlName}} {{col.SqlType}} NOT NULL DEFAULT {{col.SqlDefault}},
		  PRIMARY KEY ({{pk_col.SqlName}})
		)
	`).
		InlineMap("table", t).
		ListMap("col", "", t.Columns...).
		InlineMap("pk_col", t.PkCol()).
		String()

	_, e := db.Exec(query)
	if e != nil {
		return ErrForModel.Fmt(t.GoName).Wrap(e)
	}

	return nil
}

// Drop removes the table from the passed database. It
// assumes the database is open. If the table doesn't exist
// then nothing happens and no error is returned.
func (t ModelTable) Drop(db *sql.DB) error {
	query := nidoking.Given(`
		DROP TABLE IF EXISTS {{table.SqlName}}
	`).
		InlineMap("table", t).
		String()

	_, e := db.Exec(query)
	if e != nil {
		return ErrForModel.Fmt(t.GoName).Wrap(e)
	}

	return nil
}

// Upsert inserts the object if it doesn't exist within the
// database, else it updates the object. It assumes the
// database is open and the table exists.
//
// Inserting via an object that only represents part of a
// table will cause the unset columns to default to
// their zero value. When updating, the ID (primary key)
// field is used to target the row but is not updated.
func (t ModelTable) Upsert(db *sql.DB, object any) error {
	pkCol := t.PkCol()
	nonPkCols := t.NonPkCols()

	if pkCol == (ModelColumn{}) {
		return ErrForModel.
			Fmt(t.GoName).
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
		InlineMap("table", t).
		ListMap("col", ",", t.Columns...).
		ListRepeat("q_marks", ",", "?", len(t.Columns)).
		InlineMap("pk_col", pkCol).
		ListMap("non_pk_col", ",", nonPkCols...).
		String()

	// Because PK is in the WHERE clause
	cols := append(t.Columns, pkCol)
	values := extractColumnValues(cols, object)

	_, e := db.Exec(query, values...)
	if e != nil {
		return ErrForModel.Fmt(t.GoName).Wrap(e)
	}

	return nil
}

// Select returns all table rows within the passed
// database. It assumes the database is open and the table
// exists.
func (t ModelTable) Select[T any](db *sql.DB) ([]T, error) {
	query := nidoking.Given(`
		SELECT
			{{col.SqlName}}
		FROM
			{{table.SqlName}}
	`).
		InlineMap("table", t).
		ListMap("col", ",", t.Columns...).
		String()

	rows, e := db.Query(query)
	if e != nil {
		return nil, ErrForModel.Fmt(t.GoName).Wrap(e)
	}

	return t.scanRows[T](rows)
}

// SelectById returns the table row with the passed ID
// (primary key) from the passed database. It assumes the
// database is open and the table exists. The first return
// value will contain the row data and second will be true
// if the row exists, else zero value and false.
func (t ModelTable) SelectById[T any](db *sql.DB, id any) (T, bool, error) {
	query := nidoking.Given(`
		SELECT
			{{col.SqlName}}
		FROM
			{{table.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		ListMap("col", ",", t.Columns...).
		InlineMap("table", t).
		InlineMap("pk_col", t.PkCol()).
		String()

	rows, e := db.Query(query, id)
	if e != nil {
		var zero T
		return zero, false, ErrForModel.Fmt(t.GoName).Wrap(e)
	}

	return t.scanFirstRow[T](rows)
}

// DeleteById removes the table row with the given ID
// (primary key) from the passed database. It assumes the
// database is open and the table exists. If no matching
// row is found then nothing happens and no error is
// returned.
func (t ModelTable) DeleteById(db *sql.DB, id any) error {
	query := nidoking.Given(`
		DELETE FROM
			{{table.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		InlineMap("table", t).
		InlineMap("pk_col", t.PkCol()).
		String()

	_, e := db.Exec(query, id)
	if e != nil {
		return ErrForModel.Fmt(t.GoName).Wrap(e)
	}

	return nil
}

func extractColumnValues(
	columns []ModelColumn,
	object any,
) []any {
	value := ref.ValueOf(object)
	result := make([]any, len(columns))

	for i := 0; i < len(columns); i++ {
		field := value.FieldByName(columns[i].GoName)
		result[i] = field.Interface()
	}

	return result
}

func (t ModelTable) scanRows[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	values, valuePtrs := t.newValueContainers()
	var results []T

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrRowScan.Fmt(i).Wrap(e)
		}

		var item T
		t.populate(&item, values)
		results = append(results, item)
	}

	e := rows.Err()
	if e != nil {
		return nil, e
	}

	return results, nil
}

func (t ModelTable) scanFirstRow[T any](rows *sql.Rows) (T, bool, error) {
	defer rows.Close()
	var zero T

	if !rows.Next() {
		return zero, false, nil
	}

	values, valuePtrs := t.newValueContainers()
	e := rows.Scan(valuePtrs...)
	if e != nil {
		return zero, false, ErrRowScan.Fmt(0).Wrap(e)
	}

	e = rows.Err()
	if e != nil {
		return zero, false, e
	}

	var item T
	t.populate(&item, values)
	return item, true, nil
}

func (t ModelTable) newValueContainers() ([]any, []any) {
	colCount := len(t.Columns)
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

	for i, col := range t.Columns {
		values[i] = col.New[any]()
		valuePtrs[i] = &values[i]
	}

	return values, valuePtrs
}

func (t ModelTable) populate[T any](item *T, values []any) {
	objVal := ref.ValueOf(item).Elem()

	for valueIdx, col := range t.Columns {
		v := ref.ValueOf(values[valueIdx])
		fieldVal := objVal.Field(col.GoIndex)

		if v.CanConvert(fieldVal.Type()) {
			v = v.Convert(fieldVal.Type())
		} else {
			// TODO: Rethink and tidy
			panic("Can't convert from " + v.Type().Name() + " to " + fieldVal.Type().Name())
		}

		fieldVal.Set(v)
	}
}

// New creates a instance of the columns GoType. It will
// contain the type's zero value.
func (c ModelColumn) New[T any]() T {
	return ref.New(c.GoType).Interface().(T)
}
