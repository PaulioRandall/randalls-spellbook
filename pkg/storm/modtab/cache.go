package modtab

import (
	"database/sql"
	ref "reflect"
)

type CacheEntry struct {
	Exists bool
	Table  ModelTable
}

type CachedMapper map[ref.Type]CacheEntry

func (tm CachedMapper) Map(db *sql.DB, m any) (ModelTable, bool, error) {
	t := ref.TypeOf(m)

	entry, ok := tm[t]
	if ok {
		return entry.Table, entry.Exists, nil
	}

	table, exists, e := Map(db, m)
	if e != nil {
		return ModelTable{}, false, e
	}

	tm[t] = CacheEntry{
		Exists: exists,
		Table:  table,
	}

	return table, exists, nil
}
