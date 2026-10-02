package modtab

import (
	"database/sql"
	"reflect"
)

// CachedMapper adapts the [Map] function to cache
// results. The purpose of the cache is to minimise
// database calls accessing table information during bulk
// operations such as inserts, updates, and deletes.
//
// Because changes to table structure can be modified at
// any time, independent of this library, it's not possible
// to be certain about the integrity of any [Model]
// and the state of the database table at anytime. As a
// minimum, we have to assume a table's structure will not
// change (or table deleted) for the duration of a request.
// However, the user programmer (you) will usually have a
// clear idea about what can change and when, thus, the
// user programmer is charged with managing the cache.
//
// Because the cache is a map, Go's len, delete, and
// clear functions work directly on instances of it.
type CachedMapper map[reflect.Type]Model

// Map checks the cache for an existing model and returns
// if found, else a call to [Map] is made and the result
// cached only if the table exists.
func (tm CachedMapper) Map(
	db *sql.DB,
	object any,
) (Model, bool, error) {
	t := reflect.TypeOf(object)

	cachedModel, ok := tm[t]
	if ok {
		return cachedModel, true, nil
	}

	model, exists, e := Map(db, object)
	if e != nil {
		return Model{}, false, e
	}

	if exists {
		tm[t] = model
	}

	return model, exists, nil
}

// Clear removes all entries from the cache. This may also
// be done using Go's clear function, e.g.:
//
//	cache := CachedMapper{}
//	clear(cache)
func (cm CachedMapper) Clear() {
	clear(cm)
}

// ClearType removes the model represented by the object,
// if it exists.
func (cm CachedMapper) ClearType(object any) {
	t := reflect.TypeOf(object)
	delete(cm, t)
}

// ClearTable removes all models associated with the
// specified table name.
func (cm CachedMapper) ClearTable(name string) {
	for t, model := range cm {
		if model.SqlName == name {
			delete(cm, t)
		}
	}
}
