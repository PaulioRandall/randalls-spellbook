package storm

import (
	ref "reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// TODO: Returns error when model doesn't contain primary
//       key.

// Put inserts or updates (upserts) every object in the
// passed list. Inserting via an object that only
// represents part of a table will cause the unfilled
// columns to default to their zero value. When updating,
// the primary key (ID) field is used to identify the row
// but the primary key column is not updated.
func (st *Storm) Put[T any](objects ...T) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	mapper := newModelMapper(st)

	for _, o := range objects {
		table, e := mapper.getOrCreate(o)
		if e == nil {
			e = st.upsert(o, table)
		}

		if e != nil {
			return st.errForModel(o, e)
		}
	}

	return nil
}

func (st *Storm) upsert(object any, table ModelTable) error {
	pkCol := table.PkCol()
	nonPkCols := table.NonPkCols()

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
		InlineMap("table", table).
		ListMap("col", ",", table.Columns...).
		ListRepeat("q_marks", ",", "?", len(table.Columns)).
		InlineMap("pk_col", pkCol).
		ListMap("non_pk_col", ",", nonPkCols...).
		String()

	// Because PK is in the WHERE clause
	cols := append(table.Columns, pkCol)
	values := extractColumnValues(cols, object)

	_, e := st.db.Exec(query, values...)
	if e == nil {
		return nil
	}

	// We know primary key value is at the end.
	id := cols[len(cols)-1]
	return ErrForObject.Fmt(id).Wrap(e)
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
