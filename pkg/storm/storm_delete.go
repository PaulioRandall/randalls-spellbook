package storm

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoking"
)

// Delete removes the records with the given IDs from
// the table associated with the passed model. If no
// record is found then nothing happens. The model's
// type must match a registered type or an error is
// returned.
func (st *Storm) Delete[T, ID any](
	model T,
	ids ...ID,
) (e error) {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	for _, id := range ids {
		e = st.deleteForId(model, id)
		if e != nil {
			return st.errForObject(model, e, id)
		}
	}

	return nil
}

func (st *Storm) deleteForId(model any, id any) error {
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
