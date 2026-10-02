package modtab

import (
	"database/sql"
	"errors"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
)

// Map parses the model and modifies it to align with its
// associated table within the database. If the table
// doesn't currently exist then the returned result will be
// the same as that returned by [Parse].
func Map(db *sql.DB, model any) (ModelTable, bool, error) {
	var zero ModelTable

	modelTable, e := Parse(model)
	if e != nil {
		// Issue parsing model.
		return zero, false, e
	}

	tableCols, found, e := queryTable(db, modelTable.SqlName)
	if e != nil {
		// Issue querying table data.
		return zero, false, ErrForModel.
			Fmt(modelTable.GoName).
			Wrap(e)
	}

	if !found {
		// Table doesn't exist yet so pass back without
		// filtering.
		return modelTable, false, nil
	}

	modelTable.Columns, e = filterAndCheckColumns(
		modelTable.Columns,
		tableCols,
	)

	if e != nil {
		// Issue filtering table data.
		return zero, false, ErrForModel.
			Fmt(modelTable.GoName).
			Wrap(e)
	}

	return modelTable, true, nil
}

func queryTable(
	db *sql.DB, tableName string,
) ([]scumble.TableInfo, bool, error) {

	// Check if table exists at all.
	_, e := scumble.QuerySqliteSchema(db, tableName)
	if errors.Is(e, scumble.ErrEntityNotFound) {
		return nil, false, nil
	}

	if e != nil {
		return nil, false, e
	}

	tableCols, e := scumble.QueryTableInfo(db, tableName)
	if e != nil {
		return nil, false, e
	}

	return tableCols, true, nil
}

func filterAndCheckColumns(
	modelCols []ModelColumn,
	tableCols []scumble.TableInfo,
) ([]ModelColumn, error) {
	var filteredCols []ModelColumn

	for _, modelCol := range modelCols {
		var tableCol scumble.TableInfo

		// Find the table column that the model column maps to.
		for _, c := range tableCols {
			if modelCol.SqlName == c.Name {
				tableCol = c
			}
		}
		if tableCol == (scumble.TableInfo{}) {
			// Ignore column if not in database table.
			continue
		}

		// Update primary key status as Parse assumes the
		// first exported field is the primary key, which might
		// not be true if the passed model is not the same type
		// as the the one that was used to create the database
		// table.
		modelCol.PrimaryKey = tableCol.PrimaryKey > 0

		// Check types are compatible.
		if modelCol.SqlType != tableCol.Type {
			return nil, ErrForField.
				Fmt(modelCol.SqlName).
				WrapIn(ErrTypeMismatch).
				Fmt(modelCol.SqlType, tableCol.Type)
		}

		filteredCols = append(filteredCols, modelCol)
	}

	return filteredCols, nil
}
