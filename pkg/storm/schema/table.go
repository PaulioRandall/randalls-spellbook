package schema

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

type SqlTable struct {
	Name    string
	Columns []SqlColumn
}

type SqlColumn struct {
	Name       string
	Type       string
	Default    any
	PrimaryKey bool
}

var (
	ErrQueryTable = sin.Template(
		"For table '%s'",
	)

	ErrParseDefaultValue = sin.Template(
		"Parse issue with default value for column '%s'",
	)
)

// PrimaryKeyColumn returns the column representing the
// table's primary key. Zero value is returned if not
// found.
func (t SqlTable) PrimaryKeyColumn() SqlColumn {
	for _, col := range t.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return SqlColumn{}
}

// QueryTable queries for and returns a database
// representation of the SQL table called tableName. If
// it doesn't exist a zero value is returned.
func QueryTable(db *sql.DB, tableName string) (SqlTable, error) {
	_, e := QuerySqliteSchema(db, tableName)
	if errors.Is(e, ErrEntityNotFound) {
		return SqlTable{}, nil
	}

	if e != nil {
		return SqlTable{}, ErrQueryTable.Fmt(tableName).Wrap(e)
	}

	tableInfo, e := QueryTableInfo(db, tableName)
	if e != nil {
		return SqlTable{}, ErrQueryTable.Fmt(tableName).Wrap(e)
	}

	columns, e := parseColumns(tableInfo)
	if e != nil {
		return SqlTable{}, ErrQueryTable.Fmt(tableName).Wrap(e)
	}

	table := SqlTable{
		Name:    tableName,
		Columns: columns,
	}

	return table, nil
}

func parseColumns(tableInfo []TableInfo) ([]SqlColumn, error) {
	results := make([]SqlColumn, len(tableInfo))

	for i, col := range tableInfo {
		dv, e := parseDefault(col)

		if e != nil {
			return nil, ErrParseDefaultValue.Fmt(col.Name).Wrap(e)
		}

		results[i] = SqlColumn{
			Name:       col.Name,
			Type:       col.Type,
			Default:    dv,
			PrimaryKey: col.PrimaryKey > 0,
		}
	}

	return results, nil
}

func parseDefault(col TableInfo) (any, error) {
	if !col.HasDefaultValue {
		return nil, nil
	}

	switch col.Type {
	case "INTEGER":
		return strconv.Atoi(col.DefaultValue)
	case "REAL":
		return strconv.ParseFloat(col.DefaultValue, 64)
	case "TEXT":
		return col.DefaultValue, nil
	default:
		return nil, sin.Fmt(
			"Unsupported SQLite column type '%s'",
			col.Type,
		)
	}
}
