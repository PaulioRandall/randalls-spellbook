package wizzard

import (
	"database/sql"
	"reflect"
)

// ModelCache stores mappings between models and the GoType
// they model. ModelCache ignores the table name so it's
// recommended to always use [TableCache] instead. The
// purpose of the cache is to minimise database calls
// accessing table information.
//
// Because the cache is a map, Go's len, delete, and
// clear functions work directly on instances of it.
type ModelCache map[reflect.Type]Model

// Map checks the cache for an existing model of the
// object's type and returns it if found, else a call to
// [Map] is made and the result cached only if the database
// table exists.
func (mc ModelCache) Map(
	db *sql.DB,
	object any,
) (Model, bool, error) {
	objectType := derefObjectType(object)
	return mc.mapTypeAsModel(
		db,
		objectType.Name(),
		objectType,
	)
}

func (mc ModelCache) mapTypeAsModel(
	db *sql.DB,
	table string,
	objectType reflect.Type,
) (Model, bool, error) {
	cachedModel, ok := mc[objectType]
	if ok {
		return cachedModel, true, nil
	}

	model, exists, e := mapTypeAsModel(db, table, objectType)
	if e != nil {
		return Model{}, false, e
	}

	if exists {
		mc[objectType] = model
	}

	return model, exists, nil
}

// ClearType removes the model represented by the object,
// if it exists.
func (mc ModelCache) ClearType(object any) {
	t := reflect.TypeOf(object)
	delete(mc, t)
}

// TableCache stores mappings between models and the
// GoTypes they model with separate [ModelCache]s for
// each table.
//
// Because changes to table structure can be made at
// anytime, independent of this library, it's not possible
// to be certain about the integrity of any [Model]
// and the state of the database table at anytime. As a
// minimum, we have to assume a table's structure will not
// change (or table deleted) for the duration of an
// operation on a [Model]. However, the user programmer
// (you) will usually have a clear idea about what can
// change and when, thus, the user programmer is charged
// with managing the cache.
//
// Because the cache is a map, Go's len, delete, and
// clear functions work directly on instances of it.
type TableCache map[string]ModelCache

// Map checks the cache for an existing model of the
// object's type and returns it if found, else a call to
// [Map] is made and the result cached only if the database
// table exists.
func (cm TableCache) Map(
	db *sql.DB,
	object any,
) (Model, bool, error) {
	objectType := derefObjectType(object)

	return cm.mapTypeAsModel(
		db,
		objectType.Name(),
		objectType,
	)
}

// MapAs is the same as [TableCache.Map] but works the
// same as [MapAs] function instead.
func (cm TableCache) MapAs(
	db *sql.DB,
	table string,
	object any,
) (Model, bool, error) {
	objectType := derefObjectType(object)

	return cm.mapTypeAsModel(
		db,
		table,
		objectType,
	)
}

func (cm TableCache) mapTypeAsModel(
	db *sql.DB,
	table string,
	objectType reflect.Type,
) (Model, bool, error) {
	mc, ok := cm[table]

	if !ok {
		mc = ModelCache{}
		cm[table] = mc
	}

	return mc.mapTypeAsModel(
		db,
		table,
		objectType,
	)
}

// Clear removes all entries from the cache. This may also
// be done using Go's clear function.
func (cm TableCache) Clear() {
	clear(cm)
}

// ClearType removes the model represented by the object,
// if it exists.
func (cm TableCache) ClearType(object any) {
	t := reflect.TypeOf(object)
	for _, mc := range cm {
		delete(mc, t)
	}
}

// ClearTable removes all models associated with the
// specified table name. This may also be done using Go's
// delete function.
func (cm TableCache) ClearTable(table string) {
	delete(cm, table)
}
