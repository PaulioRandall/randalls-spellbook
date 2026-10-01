package storm

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/modtab"
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

	mapper := modtab.CachedMapper{}

	for _, id := range ids {
		table, found, e := mapper.Map(st.db, model)
		if e != nil {
			return st.errForObject(model, e, id)
		}

		if !found {
			continue
		}

		e = table.DeleteById(st.db, id)
		if e != nil {
			return st.errForObject(model, e, id)
		}
	}

	return nil
}
