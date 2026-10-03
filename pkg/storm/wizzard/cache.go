package wizzard

import (
	"database/sql"
	"reflect"
)

// CachedMapper is an adapter for the [Map] function to
// add result caching. The purpose of the cache is to
// minimise database calls accessing table information
// during bulk operations such as inserts, updates, and
// deletes.
//
// Because changes to table structure can be made at
// anytime, independent of this library, it's not possible
// to be certain about the integrity of any [Model]
// and the state of the database table at anytime. As a
// minimum, we have to assume a table's structure will not
// change (or table deleted) for the duration of an
// operation on a [Model]. However, the user programmer
// (you) will usually have a clear idea about what can
// change and when, thus, th user programmer is charged
// with managing the cache.
//
// Because the cache is a map, Go's len, delete, and
// clear functions work directly on instances of it.
type CachedMapper map[reflect.Type]Model

// Map checks the cache for an existing model of the
// object's type and returns it if found, else a call to
// [Map] is made and the result cached only if the database
// table exists.
func (cm CachedMapper) Map(
	db *sql.DB,
	object any,
) (Model, bool, error) {
	t := reflect.TypeOf(object)

	cachedModel, ok := cm[t]
	if ok {
		return cachedModel, true, nil
	}

	model, exists, e := Map(db, object)
	if e != nil {
		return Model{}, false, e
	}

	if exists {
		cm[t] = model
	}

	return model, exists, nil
}

// Clear removes all entries from the cache. This may also
// be done using Go's clear function.
func (cm CachedMapper) Clear() {
	clear(cm)
}

// ClearType removes the model represented by the object,
// if it exists. This may also be done using Go's delete
// function.
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
