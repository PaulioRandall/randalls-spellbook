package mapper

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
)

// NEXT
// 1. Write tests for this package.
// 2. Reorg, refactor, & clean up package.

func MapModel(db *sql.DB, model any) (Table, error) {
	var zero Table

	pTable := ParseModel(model)
	qCols, found, e := queryModel(db, pTable.SqlName)

	if e != nil {
		return zero, e
	}

	if !found {
		// Table doesn't exist yet so pass the Table derived
		// from the model without filtered columns.
		return pTable, nil
	}

	pTable.Columns = filterAndCheckColumns(pTable.Columns, qCols)
	return pTable, nil
}

func filterAndCheckColumns(
	pCols []Column,
	qCols []sqlQueryCol,
) []Column {
	var filteredCols []Column

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
			// TODO: Refactor to return an error instead.
			msg := fmt.Sprintf(
				"For column '%s', model's field type '%s' is not compatible with existing column type '%s'",
				pCol.SqlName,
				pCol.SqlType,
				qCol.SqlType,
			)
			panic(msg)
		}

		filteredCols = append(filteredCols, pCol)
	}

	return filteredCols
}

func findQueryColumn(
	qCols []sqlQueryCol,
	pCol Column,
) (sqlQueryCol, bool) {
	for _, qCol := range qCols {
		if pCol.SqlName == qCol.SqlName {
			return qCol, true
		}
	}
	return sqlQueryCol{}, false
}

type sqlQueryCol struct {
	SqlType    string
	SqlName    string
	PrimaryKey bool
}

func queryModel(db *sql.DB, tableName string) ([]sqlQueryCol, bool, error) {
	_, e := schema.QuerySqliteSchema(db, tableName)

	// TODO: Redundant?
	if errors.Is(e, schema.ErrEntityNotFound) {
		return nil, false, nil
	}

	if e != nil {
		return nil, false, e
	}

	colInfo, e := schema.QueryTableInfo(db, tableName)
	if e != nil {
		return nil, false, e
	}

	cols := make([]sqlQueryCol, len(colInfo), len(colInfo))

	for i, col := range colInfo {
		cols[i] = sqlQueryCol{
			SqlType:    col.Type,
			SqlName:    col.Name,
			PrimaryKey: col.PrimaryKey > 0,
		}
	}

	return cols, true, nil
}
