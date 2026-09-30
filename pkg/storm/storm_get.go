package storm

import (
	"database/sql"
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

// Select returns all records for the table associated
// with the passed model.
func (st *Storm) Select[T any](model T) (result []T, e error) {
	var query string
	var rows *sql.Rows

	if !st.IsOpen() {
		return nil, ErrNotOpen
	}

	table, found, e := st.getTable(model)
	if e != nil {
		return nil, e
	}

	if !found {
		return nil, nil
	}

	query = nidoking.Given(`
		SELECT
			{{col.SqlName}}
		FROM
			{{table.SqlName}}
	`).
		InlineMap("table", table).
		ListMap("col", ",", table.Columns...).
		String()

	rows, e = st.db.Query(query)
	if e != nil {
		goto Err
	}

	result, e = st.scanSelectedRows[T](table, rows)
	if e != nil {
		goto Err
	}

	return result, nil

Err:
	return nil, ErrTableRequest.Fmt(typeName(model)).Wrap(e)
}

// SelectById returns the record with the given id from
// the table associated with the passed model. If no
// record is found then an error is returned. The model's
// type must match a registered type or an error is
// returned.
//
//	object, err := SelectById(Model{}, 123)
func (st *Storm) SelectById[T, ID any](
	model T,
	id ID,
) (result T, e error) {
	var empty T
	var query string
	var rows *sql.Rows
	var resultSet []T
	var ok bool

	if !st.IsOpen() {
		return empty, ErrNotOpen
	}

	table, found, e := st.getTable(model)
	if e != nil {
		goto Err
	}

	if !found {
		e = ErrObjectNotFound.Fmt(table.GoName, id)
		goto Err
	}

	query = nidoking.Given(`
		SELECT
			{{col.SqlName}}
		FROM
			{{table.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		ListMap("col", ",", table.Columns...).
		InlineMap("table", table).
		InlineMap("pk_col", table.PrimaryKeyColumn()).
		String()

	rows, e = st.db.Query(query, id)
	if e != nil {
		goto Err
	}

	resultSet, e = st.scanSelectedRows[T](table, rows)
	if e != nil {
		goto Err
	}

	result, ok = getFirstItemIfArray[T](resultSet)
	if !ok {
		e = ErrObjectNotFound.Fmt(table.GoName, id)
		goto Err
	}

	return result, nil

Err:
	return empty, sin.Fmt("For object with ID '%v'", id).
		Wrap(e).
		WrapIn(ErrTableRequest).
		Fmt(typeName(model))
}

func (st *Storm) scanSelectedRows[T any](
	table ModelTable,
	rows *sql.Rows,
) ([]T, error) {
	values, valuePtrs := createValueContainers(table)
	var result []T

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrScanningRow.Fmt(i).Wrap(e)
		}

		object := constructObject[T](table, values)
		result = append(result, object)
	}

	e := rows.Err()
	if e != nil {
		return nil, e
	}

	return result, nil
}

func createValueContainers(table ModelTable) ([]any, []any) {
	colCount := len(table.Columns)
	values := make([]any, colCount)
	valuePtrs := make([]any, colCount)

	for i, col := range table.Columns {
		values[i] = col.New[any]()
		valuePtrs[i] = &values[i]
	}

	return values, valuePtrs
}

func constructObject[T any](
	table ModelTable,
	values []any,
) T {
	var result T

	objVal := reflect.ValueOf(&result).Elem()

	for valueIdx, col := range table.Columns {
		v := reflect.ValueOf(values[valueIdx])
		fieldVal := objVal.Field(col.GoIndex)

		if v.CanConvert(fieldVal.Type()) {
			v = v.Convert(fieldVal.Type())
		} else {
			// TODO: Rethink and tidy
			panic("Can't convert from " + v.Type().Name() + " to " + fieldVal.Type().Name())
		}

		fieldVal.Set(v)
	}

	return result
}

func getFirstItemIfArray[T any](v any) (T, bool) {
	var empty T
	rv := reflect.ValueOf(v)

	isArray := rv.Kind() == reflect.Array
	isSlice := rv.Kind() == reflect.Slice

	if !isArray && !isSlice {
		return empty, false
	}

	if rv.Len() == 0 {
		return empty, false
	}

	return rv.Index(0).Interface().(T), true
}
