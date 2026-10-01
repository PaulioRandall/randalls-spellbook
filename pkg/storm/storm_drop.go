package storm

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/modtab"
)

// Drop removes a table from the database. If the target
// table doesn't exist then nothing happens and no error is
// returned. Related model cache entries are also removed.
//
// All data is deleted in the process and there's no way to
// restore it. To protect data, create backups of the
// database file.
func (st *Storm) Drop(models ...any) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	mapper := modtab.CachedMapper{}

	for _, m := range models {
		table, found, e := mapper.Map(st.db, m)
		if e != nil {
			return st.errForModel(m, e)
		}

		if !found {
			continue
		}

		e = table.Drop(st.db)
		if e != nil {
			return st.errForModel(m, e)
		}
	}

	return nil
}
