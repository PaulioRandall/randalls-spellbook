package storm

import (
	"database/sql"
	"errors"
	ref "reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

var (
	ErrMapModel = sin.Template(
		"Regarding model '%s' (table '%s')",
	)

	ErrFieldTypeMismatch = sin.Template(
		"Model's field type '%s' is not compatible with existing column type '%s'",
	)
)

type modelMapperCacheEntry struct {
	exists bool
	table  ModelTable
}

type modelMapper struct {
	st    *Storm
	cache map[ref.Type]modelMapperCacheEntry
}

func newModelMapper(st *Storm) *modelMapper {
	return &modelMapper{
		st:    st,
		cache: map[ref.Type]modelMapperCacheEntry{},
	}
}

func (tm *modelMapper) get(m any) (ModelTable, bool, error) {
	t := typeOf(m)

	entry, ok := tm.cache[t]
	if ok {
		return entry.table, entry.exists, nil
	}

	table, exists, e := MapModel(tm.st.db, m)
	if e != nil {
		return ModelTable{}, false, e
	}

	tm.cache[t] = modelMapperCacheEntry{
		exists: exists,
		table:  table,
	}

	return table, exists, nil
}

func (tm *modelMapper) getOrCreate(m any) (ModelTable, error) {
	table, exists, e := tm.get(m)
	if e != nil || exists {
		return table, e
	}

	e = tm.st.createTable(table)
	if e != nil {
		return ModelTable{}, e
	}

	return table, nil
}

type sqlQueryCol struct {
	SqlType    string
	SqlName    string
	PrimaryKey bool
}

func MapModel(db *sql.DB, model any) (ModelTable, bool, error) {
	var zero ModelTable

	pTable := ParseModel(model)
	qCols, found, e := queryModel(db, pTable.SqlName)

	if e != nil {
		return zero, false, ErrMapModel.
			Fmt(pTable.GoName, pTable.SqlName).
			Wrap(e)
	}

	if !found {
		// ModelTable doesn't exist yet so pass back without
		// filtering.
		return pTable, false, nil
	}

	pTable.Columns, e = filterAndCheckColumns(pTable.Columns, qCols)
	if e != nil {
		return zero, false, ErrMapModel.
			Fmt(pTable.GoName, pTable.SqlName).
			Wrap(e)
	}

	return pTable, true, nil
}

func filterAndCheckColumns(
	pCols []ModelColumn,
	qCols []sqlQueryCol,
) ([]ModelColumn, error) {
	var filteredCols []ModelColumn

	for _, pCol := range pCols {
		qCol, ok := findQueryColumn(qCols, pCol)

		if !ok {
			// Ignore column as not in database table.
			continue
		}

		// Update primary key status as ParseModel assumes the
		// first exported field is the primary key, which might
		// not be true if the passed model is not the same type
		// as the the one that was used to create the database
		// table.
		pCol.PrimaryKey = qCol.PrimaryKey

		if pCol.SqlType != qCol.SqlType {
			return nil, sin.
				Fmt("For column '%s'", pCol.SqlName).
				WrapIn(ErrFieldTypeMismatch).
				Fmt(pCol.SqlType, qCol.SqlType)
		}

		filteredCols = append(filteredCols, pCol)
	}

	return filteredCols, nil
}

func findQueryColumn(
	qCols []sqlQueryCol,
	pCol ModelColumn,
) (sqlQueryCol, bool) {
	for _, qCol := range qCols {
		if pCol.SqlName == qCol.SqlName {
			return qCol, true
		}
	}
	return sqlQueryCol{}, false
}

func queryModel(
	db *sql.DB, tableName string,
) ([]sqlQueryCol, bool, error) {

	// Check if table exists at all.
	_, e := scumble.QuerySqliteSchema(db, tableName)
	if errors.Is(e, scumble.ErrEntityNotFound) {
		return nil, false, nil
	}

	if e != nil {
		return nil, false, e
	}

	colInfo, e := scumble.QueryTableInfo(db, tableName)
	if e != nil {
		return nil, false, e
	}

	cols := make([]sqlQueryCol, len(colInfo))

	for i, col := range colInfo {
		cols[i] = sqlQueryCol{
			SqlType:    col.Type,
			SqlName:    col.Name,
			PrimaryKey: col.PrimaryKey > 0,
		}
	}

	return cols, true, nil
}
