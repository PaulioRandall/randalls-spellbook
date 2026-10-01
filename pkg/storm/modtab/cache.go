package modtab

import (
	"database/sql"
	ref "reflect"
)

// CacheEntry is the entry type for [CachedMapper].
type CacheEntry struct {
	// Exists is true if the table exsits within the
	// database.
	Exists bool

	// Table is the parsed model being cached. It is the
	// value returned by [Parse] or [Map].
	Table ModelTable
}

// CachedMapper is a cache for [ModelTable]. The purpose of
// the cache is to minimise database calls accessing table
// information during bulk operations such as inserts,
// updates, and deletes.
//
// Because changes to table structure can be made at any
// time, independent of this library, it's not possible
// to be certain about the integrity of any [ModelTable] at
// anytime. As a minimum, we have to assume a table's
// structure will not change (or table deleted) for the
// duration of a request. But the user programmer (you)
// will usually know what can change when, thus, the user
// programmer must create and manage this cache.
//
// Because the cache is a map, the Go len, delete, and
// clear functions work directly on instances of the
// mapper.
type CachedMapper map[ref.Type]CacheEntry

// Map checks the cache for a parsed model first. If not
// found a call to [Map] is made and the result cached only
// if the table currently exists within the database.
func (tm CachedMapper) Map(
	db *sql.DB,
	m any,
) (ModelTable, bool, error) {
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

// DeleteModel removes the entry for the specified model,
// if it exists.
func (tm CachedMapper) DeleteModel(model any) {
	t := ref.TypeOf(model)
	delete(tm, t)
}

// DeleteTable removes all entries associated with the
// specified table name.
func (tm CachedMapper) DeleteTable(name string) {
	for t, entry := range tm {
		if entry.Table.SqlName == name {
			delete(tm, t)
		}
	}
}
