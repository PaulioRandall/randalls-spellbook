package storm

import (
	"database/sql"
	ref "reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// TODO: REfactor & tidy

// List returns all records for the table associated
// with the passed model.
func (st *Storm) List[T any](model T) (result []T, e error) {
	var query string
	var rows *sql.Rows

	if !st.IsOpen() {
		return nil, st.errNotOpen()
	}

	table, found, e := MapModel(st.db, model)
	if e != nil {
		goto Err
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

	result, e = scanSelectedRows[T](table, rows)
	if e != nil {
		goto Err
	}

	return result, nil

Err:
	return nil, st.errForModel(model, e)
}

// Get returns the record with the given id from
// the table associated with the passed model. If no
// record is found then an error is returned. The model's
// type must match a registered type or an error is
// returned.
func (st *Storm) Get[T, ID any](
	model T,
	id ID,
) (result T, e error) {
	var empty T
	var query string
	var rows *sql.Rows
	var resultSet []T
	var ok bool

	if !st.IsOpen() {
		return empty, st.errNotOpen()
	}

	table, found, e := MapModel(st.db, model)
	if e != nil {
		goto Err
	}

	if !found {
		e = ErrObjectNotFound
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
		InlineMap("pk_col", table.PkCol()).
		String()

	rows, e = st.db.Query(query, id)
	if e != nil {
		goto Err
	}

	resultSet, e = scanSelectedRows[T](table, rows)
	if e != nil {
		goto Err
	}

	result, ok = getFirstItemIfArray[T](resultSet)
	if !ok {
		e = ErrObjectNotFound
		goto Err
	}

	return result, nil

Err:
	return empty, st.errForObject(model, e, id)
}

func scanSelectedRows[T any](
	table ModelTable,
	rows *sql.Rows,
) ([]T, error) {
	values, valuePtrs := createValueContainers(table)
	var result []T

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(valuePtrs...)
		if e != nil {
			return nil, ErrRowScan.Fmt(i).Wrap(e)
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

	objVal := ref.ValueOf(&result).Elem()

	for valueIdx, col := range table.Columns {
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

	return result
}

func getFirstItemIfArray[T any](v any) (T, bool) {
	var empty T
	rv := ref.ValueOf(v)

	isArray := rv.Kind() == ref.Array
	isSlice := rv.Kind() == ref.Slice

	if !isArray && !isSlice {
		return empty, false
	}

	if rv.Len() == 0 {
		return empty, false
	}

	return rv.Index(0).Interface().(T), true
}
