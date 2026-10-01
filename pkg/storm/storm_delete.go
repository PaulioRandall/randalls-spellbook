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

	mapper := newModelMapper(st)

	for _, id := range ids {
		table, found, e := mapper.get(model)
		if e != nil {
			return st.errForObject(model, e, id)
		}

		if !found {
			continue
		}

		e = st.deleteForId(table, id)
		if e != nil {
			return st.errForObject(model, e, id)
		}
	}

	return nil
}

func (st *Storm) deleteForId(table ModelTable, id any) error {
	query := nidoking.Given(`
		DELETE FROM
			{{table.SqlName}}
		WHERE
			{{pk_col.SqlName}} = ?
	`).
		InlineMap("table", table).
		InlineMap("pk_col", table.PkCol()).
		String()

	_, e := st.db.Exec(query, id)
	return e
}
