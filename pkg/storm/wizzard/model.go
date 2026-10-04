package wizzard

import (
	"database/sql"
	"fmt"
	ref "reflect"
	"strings"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// Model maps a Go struct, or part of it, to a database
// table structure, or part of it. A Model is always
// derived from a struct, so there's a one-to-one
// relationship between a Go type and a model but
// one-to-many relationship between tables and models and
// Go types. However, it's possible for some models to
// differ by SqlName only.
type Model struct {
	// GoType is the struct type as returned by
	// reflect.TypeOf.
	GoType ref.Type

	// GoName is the struct's name, i.e. GoType.Name().
	GoName string

	// SqlName is the name of the table the struct represents
	// within a database.
	SqlName string

	// Props is an ordered list of exported field-to-column
	// mappings.
	Props []Property
}

// Property maps a Go struct's field to a database table
// column.
type Property struct {
	// GoType is the field type as returned by
	// reflect.TypeOf.
	GoType ref.Type

	// GoIndex is the field index, i.e. it's 0-based position
	// within the struct definition.
	GoIndex int

	// GoName is the field's name, i.e. GoType.Name().
	GoName string

	// SqlType is the SQLite type that maps to the GoType.
	// Note that several Go types may map to a single SQLite
	// type.
	SqlType string

	// SqlName is the name of the column within a database.
	SqlName string

	// SqlDefault is the default value given to the column if
	// a value is not provided during row insertion.
	// This is always the zero value of the GoType.
	SqlDefault any

	// IsKey is true if the property represents an ID or
	// primary key for the model. By default this is the
	// first property in the model, however, this may differ
	// if adjusting a property to align with an existing
	// table column as with the [Map] function.
	IsKey bool
}

// Parse creates a complete [Model] of the passed object's
// type. One should not assume the result will map directly
// to a table within a database; use [Map] to create a
// model that is usable with a specific database. Type
// mappings:
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
// for your structs. If you choose to use an unsigned
// integer type (unit etc) or shorter integer type
// (int8 etc) then you are responsible for ensuring or
// managing number polarity and overflow. I recommend
// sticking to int, int64, float64, and string as these are
// the least likely to cause problems.
func Parse(object any) (Model, error) {
	objectType := derefObjectType(object)
	return parseTypeAsModel(objectType.Name(), objectType)
}

// ParseAs is the same as [Parse] except the table name is
// provided explicitly.
func ParseAs(table string, object any) (Model, error) {
	objectType := derefObjectType(object)
	return parseTypeAsModel(table, objectType)
}

func parseTypeAsModel(
	table string,
	objectType ref.Type,
) (model Model, e error) {
	if objectType.Kind() != ref.Struct {
		e = ErrNotStruct.Fmt(objectType.Kind().String())
		goto Err
	}

	model.GoType = objectType
	model.GoName = objectType.Name()
	model.SqlName = table

	model.Props, e = parseProps(objectType)
	if e != nil {
		goto Err
	}

	return model, nil

Err:
	return Model{}, ErrForModel.
		Fmt(objectType.Name()).
		Wrap(e).
		WrapIn(ErrForTable).
		Fmt(table)
}

func derefObjectType(object any) ref.Type {
	typ := ref.TypeOf(object)

	for typ.Kind() == ref.Ptr {
		typ = typ.Elem()
	}

	return typ
}

func parseProps(objectType ref.Type) ([]Property, error) {
	isIdField := true

	var props []Property

	for i := 0; i < objectType.NumField(); i++ {
		f := objectType.Field(i)

		if !f.IsExported() {
			continue
		}

		prop, e := parseProp(f, i, isIdField)
		if e != nil {
			return nil, ErrForField.Fmt(f.Name).Wrap(e)
		}

		props = append(props, prop)
		isIdField = false
	}

	return props, nil
}

func parseProp(
	f ref.StructField,
	i int,
	isIdField bool,
) (Property, error) {
	var prop Property
	var defaultValue any

	if f.Type == ref.TypeOf("") {
		defaultValue = "''"
	} else {
		defaultValue = ref.Zero(f.Type).Interface()
	}

	sqlType, e := mapGoToSqlType(f.Type.Kind())
	if e != nil {
		return Property{}, e
	}

	prop.GoType = f.Type
	prop.GoIndex = i
	prop.GoName = f.Name
	prop.SqlType = sqlType
	prop.SqlName = f.Name
	prop.SqlDefault = defaultValue
	prop.IsKey = isIdField

	return prop, nil
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

// String returns a developer friendly representation of
// the model.
func (m Model) String() string {
	return m.GoTypeString() + "\n" + m.SqlTableString()
}

// GoTypeString returns a string representation of the
// model's Go type.
func (m Model) GoTypeString() string {
	sb := strings.Builder{}

	sb.WriteString("Go Type: ")
	sb.WriteString(m.GoName)

	for _, p := range m.Props {
		s := fmt.Sprintf(
			"\n\t[%d] %s %s",
			p.GoIndex,
			p.GoName,
			p.GoType.Name(),
		)
		sb.WriteString(s)
	}

	return sb.String()
}

// SqlTableString returns a string representation of the
// model's table details.
func (m Model) SqlTableString() string {
	sb := strings.Builder{}

	sb.WriteString("Sql Table: ")
	sb.WriteString(m.SqlName)

	for _, p := range m.Props {
		s := fmt.Sprintf(
			"\n\t%s %s",
			p.SqlName,
			p.SqlType,
		)
		sb.WriteString(s)
	}

	return sb.String()
}

// Represents returns true if the model represents the
// type of the passed object, and thus can be used with
// the model's methods without causing a type mismatch
// error.
func (m Model) Represents(object any) bool {
	return m.GoType == ref.TypeOf(object)
}

// KeyProp returns the property representing the model's
// key. Zero value is returned if not present.
func (m Model) KeyProp() Property {
	for _, prop := range m.Props {
		if prop.IsKey {
			return prop
		}
	}
	return Property{}
}

// NonKeyProps returns all properties except the key
// property.
func (m Model) NonKeyProps() []Property {
	result := make([]Property, 0, len(m.Props))

	for _, prop := range m.Props {
		if !prop.IsKey {
			result = append(result, prop)
		}
	}

	return result
}

// Create creates the table the model represents within the
// passed database. It assumes the database is open. If a
// table with the same name already exists then nothing
// happens and no error is returned.
func (m Model) Create(db *sql.DB) error {
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
		return ErrForModel.Fmt(m.GoName).Wrap(e)
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
func (m Model) Upsert(db *sql.DB, object any) error {
	if !m.Represents(object) {
		return ErrForModel.
			Fmt(m.GoName).
			Wrap(ErrWrongObjectType)
	}

	pkCol := m.KeyProp()
	nonPkCols := m.NonKeyProps()

	if pkCol == (Property{}) {
		return ErrForModel.
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
		return ErrForModel.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Insert inserts the object into the database. It assumes
// the database is open and the table exists. An error is
// returned if the object's type does not match the model's
// GoType. Inserting via an object that only represents
// part of a table will cause the unset columns to default
// to their zero value.
func (m Model) Insert(db *sql.DB, object any) error {
	if !m.Represents(object) {
		return ErrForModel.
			Fmt(m.GoName).
			Wrap(ErrWrongObjectType)
	}

	if m.KeyProp() == (Property{}) {
		return ErrForModel.
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
		return ErrForModel.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Update updates the object within the database. It
// assumes the database is open and the table exists. An
// error is returned if the object's type does not match
// the model's GoType. The key property is used to target
// the row but is not updated.
func (m Model) Update(db *sql.DB, object any) error {
	if !m.Represents(object) {
		return ErrForModel.
			Fmt(m.GoName).
			Wrap(ErrWrongObjectType)
	}

	pkCol := m.KeyProp()
	nonPkCols := m.NonKeyProps()

	if pkCol == (Property{}) {
		return ErrForModel.
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
		return ErrForModel.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Select returns all rows matching the passed where clause
// from the table the model represents within passed
// database. It assumes the database is open and the table
// exists.
func (m Model) Select[T any](db *sql.DB, where string, args ...any) ([]T, error) {
	var o T
	if !m.Represents(o) {
		return nil, ErrForModel.
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
		return nil, ErrForModel.Fmt(m.GoName).Wrap(e)
	}

	return m.scanRows[T](rows)
}

// SelectFirst returns the first row matching the passed
// where clause from the table the model represents within
// the passed database. It assumes the database is open and
// the table exists. The first return value will contain
// the row data and second will be true if the row exists,
// else zero value and false are returned.
func (m Model) SelectFirst[T any](db *sql.DB, where string, args ...any) (T, bool, error) {
	var zero T

	if !m.Represents(zero) {
		return zero, false, ErrForModel.
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
		return zero, false, ErrForModel.Fmt(m.GoName).Wrap(e)
	}

	return m.scanFirstRow[T](rows)
}

// Delete removes all table rows matching the given where
// clause and arguments from the table the model represents
// within the passed database. It assumes the database is
// open and the table exists. If no matching rows are found
// then nothing happens and no error is returned.
func (m Model) Delete(db *sql.DB, where string, args ...any) error {
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
		return ErrForModel.Fmt(m.GoName).Wrap(e)
	}

	return nil
}

// Drop removes the table the model represents from the
// passed database. It assumes the database is open. If the
// table doesn't exist then nothing happens and no error is
// returned.
func (m Model) Drop(db *sql.DB) error {
	query := nidoking.Given(`
		DROP TABLE IF EXISTS {{table.SqlName}}
	`).
		InlineMap("table", m).
		String()

	_, e := db.Exec(query)
	if e != nil {
		return ErrForModel.Fmt(m.GoName).Wrap(e)
	}

	return nil
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

func (m Model) scanRows[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	values, valuePtrs := m.newValueContainers()
	var results []T

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrRowScan.Fmt(i).Wrap(e)
		}

		var item T
		m.populate(&item, values)
		results = append(results, item)
	}

	e := rows.Err()
	if e != nil {
		return nil, e
	}

	return results, nil
}

func (m Model) scanFirstRow[T any](rows *sql.Rows) (T, bool, error) {
	defer rows.Close()
	var zero T

	if !rows.Next() {
		return zero, false, nil
	}

	values, valuePtrs := m.newValueContainers()
	e := rows.Scan(valuePtrs...)
	if e != nil {
		return zero, false, ErrRowScan.Fmt(0).Wrap(e)
	}

	e = rows.Err()
	if e != nil {
		return zero, false, e
	}

	var item T
	m.populate(&item, values)
	return item, true, nil
}

func (m Model) newValueContainers() ([]any, []any) {
	colCount := len(m.Props)
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

	for i, col := range m.Props {
		values[i] = col.New[any]()
		valuePtrs[i] = &values[i]
	}

	return values, valuePtrs
}

func (m Model) populate[T any](item *T, values []any) {
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

// New creates a instance of the columns GoType. It will
// contain the type's zero value. The value returned is
// explicitly cast to T.
func (p Property) New[T any]() T {
	return ref.New(p.GoType).Interface().(T)
}
