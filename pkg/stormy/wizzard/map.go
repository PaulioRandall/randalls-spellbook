package wizzard

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
)

// Map parses the object and modifies the resultant model
// to align with its associated table within the database.
// If the table doesn't currently exist then the model will
// be the same as that returned by [Parse].
func Map(db *sql.DB, object any) (Model, bool, error) {
	objectType := derefObjectType(object)
	return mapTypeAsModel(db, objectType.Name(), objectType)
}

// MapAs is the same as [Map] except the table name is
// provided explicitly.
func MapAs(db *sql.DB, table string, object any) (Model, bool, error) {
	objectType := derefObjectType(object)
	return mapTypeAsModel(db, table, objectType)
}

func mapTypeAsModel(
	db *sql.DB,
	table string,
	objectType reflect.Type,
) (Model, bool, error) {
	var zero Model

	model, e := parseTypeAsModel(table, objectType)
	if e != nil {
		// Issue parsing object. Error is already dressed up.
		return zero, false, e
	}

	cols, found, e := queryTable(db, model.SqlName)
	if e != nil {
		// Issue querying table data.
		goto Err
	}

	if !found {
		// Table doesn't exist yet so pass back without
		// filtering.
		return model, false, nil
	}

	model.Props, e = filterAndCheckColumns(
		model.Props,
		cols,
	)

	if e != nil {
		// Issue filtering table data.
		goto Err
	}

	return model, true, nil

Err:
	return zero, false, ErrForType.
		Fmt(model.GoName).
		Wrap(e).
		WrapIn(ErrForTable).
		Fmt(table)
}

func queryTable(
	db *sql.DB,
	tableName string,
) ([]scumble.TableInfo, bool, error) {

	// Check if table exists at all.
	_, e := scumble.QuerySqliteSchema(db, tableName)
	if errors.Is(e, scumble.ErrEntityNotFound) {
		return nil, false, nil
	}

	if e != nil {
		return nil, false, e
	}

	cols, e := scumble.QueryTableInfo(db, tableName)
	if e != nil {
		return nil, false, e
	}

	return cols, true, nil
}

func filterAndCheckColumns(
	props []Property,
	cols []scumble.TableInfo,
) ([]Property, error) {
	var filteredProps []Property

	for _, prop := range props {
		var col scumble.TableInfo

		// Find the table column that the property maps to.
		for _, c := range cols {
			if prop.SqlName == c.Name {
				col = c
			}
		}

		if col == (scumble.TableInfo{}) {
			// Ignore the property if not in database table.
			continue
		}

		// Update key status as the Parse function assumes the
		// first exported field is the key, which might not be
		// true if the model is not the same as the one used to
		// create the database table.
		prop.IsKey = col.PrimaryKey > 0

		// Check types are compatible.
		if prop.SqlType != col.Type {
			return nil, ErrForField.
				Fmt(prop.SqlName).
				WrapIn(ErrTypeMismatch).
				Fmt(prop.SqlType, col.Type)
		}

		filteredProps = append(filteredProps, prop)
	}

	return filteredProps, nil
}
