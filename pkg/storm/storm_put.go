package storm

import (
	ref "reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// Put inserts or updates all passed objects into the
// database. If a table doesn't exist for an object, it
// will be created as if passed to [Storm.Create].
func (st *Storm) Put[T any](objects ...T) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, o := range objects {
		e := st.insertObject(o)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(o)).Wrap(e)
		}
	}

	return nil
}

func (st *Storm) insertObject(object any) error {
	table, e := st.getOrCreateTable(object)
	if e != nil {
		return e
	}

	query := nidoking.Given(`
		INSERT INTO {{table.SqlName}} (
		  {{col.SqlName}}
		)
		VALUES (
			{{q_marks}}
		)
	`).
		InlineMap("table", table).
		ListMap("col", ",", table.Columns...).
		ListRepeat("q_marks", ",", "?", len(table.Columns)).
		String()

	values := extractColumnValues(table.Columns, object)
	_, e = st.db.Exec(query, values...)
	return e
}

// Update updates the objects within the database. Each
// object's type must match a registered type or an error
// is returned. All fields are updated except the ID
// field, which is used to determine which record to
// update.
//
//	alice := Player{
//		Id: 69,
//		Name: "Alice",
//		RoleId: 5,
//	}
//	err := db.Insert(alice)
//	// YUDO: Handle error.
//
//	alice.Name = "Alicia"
//	err = db.Update(alice)
func (st *Storm) Update[T any](objects ...T) error {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, obj := range objects {
		e := st.updateObject(obj)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(obj)).Wrap(e)
		}
	}

	return nil
}

func (st *Storm) updateObject(object any) error {
	table, found, e := st.getTable(object)
	if e != nil {
		return e
	}

	if !found {
		return ErrNoSuchTable.Fmt(typeName(object))
	}

	pkCol := table.PrimaryKeyColumn()
	nonPkCols := table.NonPrimaryKeyColumns()

	query := nidoking.Given(`
			UPDATE
				{{table.SqlName}}
			SET
				{{non_pk_col.SqlName}} = ?
			WHERE
				{{pk_col.SqlName}} = ?
		`).
		InlineMap("table", table).
		ListMap("non_pk_col", ",", nonPkCols...).
		InlineMap("pk_col", pkCol).
		String()

	// Because PK is the last SQL parameter.
	cols := append(nonPkCols, pkCol)
	values := extractColumnValues(cols, object)
	_, e = st.db.Exec(query, values...)
	return e
}

func extractColumnValues(
	columns []ModelColumn,
	object any,
) []any {
	value := ref.ValueOf(object)
	result := make([]any, len(columns))

	for i := 0; i < len(columns); i++ {
		col := columns[i]
		field := value.FieldByName(col.GoName)
		result[i] = field.Interface()
	}

	return result
}
