package storm

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/modtab"
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

	mapper := modtab.CachedMapper{}

	for _, o := range objects {
		table, exists, e := mapper.Map(st.db, o)
		if e != nil {
			return st.errForModel(o, e)
		}

		if !exists {
			e = table.Create(st.db)
			if e != nil {
				return st.errForModel(o, e)
			}
		}

		if e == nil {
			e = table.Upsert(st.db, o)
		}

		if e != nil {
			return st.errForModel(o, e)
		}
	}

	return nil
}
