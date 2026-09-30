package storm

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// DeleteById removes the records with the given ids from
// the table associated with the passed model. If no
// record is found then nothing happens. The model's
// type must match a registered type or an error is
// returned.
func (st *Storm) DeleteById[T, ID any](
	model T,
	ids ...ID,
) (e error) {
	if !st.IsOpen() {
		return ErrNotOpen
	}

	for _, id := range ids {
		e = st.deleteById(model, id)
		if e != nil {
			return ErrTableRequest.Fmt(typeName(model)).Wrap(e)
		}
	}

	return nil
}

func (st *Storm) deleteById(model any, id any) error {
	table, found, e := st.getTable(model)
	if e != nil {
		return e
	}

	if !found {
		return nil
	}

	query := nidoking.Given(`
		DELETE FROM
			{{table.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		InlineMap("table", table).
		InlineMap("pk_col", table.PrimaryKeyColumn()).
		String()

	_, e = st.db.Exec(query, id)
	return e
}
