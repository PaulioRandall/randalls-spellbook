package schema

import (
	"database/sql"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

// TableInfo is a Go representation of table_info within
// SQLite3. Note that internally generated columns are
// not returned via table_info PRAGMA. Column ID (cid) is
// omitted because it's not a reliable reference to
// anything outside the query that generated it.
//
// See https://sqlite.org/pragma.html#pragma_table_info
type TableInfo struct {
	// Name or 'name' is the name of the column as the user
	// defined it.
	Name string

	// Type or 'type' is the column type.
	Type string

	// NotNull or 'notnull' is true if the column was given
	// a not null constraint.
	NotNull bool

	// HasDefaultValue is true if the 'dflt_value' column is
	// not null. This field is needed to differentiate
	// between null and an empty string without resorting to
	// making DefaultValue a pointer to a string. This
	// approach felt nicer when I wrote the test and example
	// code.
	HasDefaultValue bool

	// DefaultValue or 'dflt_value' is the default value used
	// when the a row is created but a value for the column
	// was not provided. Sqlite stores and returns the
	// value as text (string) so you must parse it yourself.
	DefaultValue string

	// PrimaryKey or 'pk' is the 1-based index of the column
	// within the primary key. It will be 0 if the column is
	// not part of the primary key.
	PrimaryKey int
}

var (
	// ErrQueryTableInfo occurs when querying table_info
	// PRAGMA.
	ErrQueryTableInfo = sin.Template(
		"Could not query table_info for table '%s'",
	)

	// ErrScanTableInfo occurs when scanning rows returned
	// from table_info PRAGMA.
	ErrScanTableInfo = sin.Template(
		"Could not scan table_info row %d for table '%s'",
	)
)

// QueryTableInfo queries for and returns the info on all
// user columns (non-internally generated columns) within
// table_info for the passed tableName. If the table does
// not exist or does not contain any user columns then nil
// is returned.
func QueryTableInfo(
	db *sql.DB,
	tableName string,
) ([]TableInfo, error) {

	query := `
		SELECT
			name,
			type,
			"notnull",
			dflt_value,
			pk
		FROM
			pragma_table_info(?)
		`

	rows, e := db.Query(query, tableName)
	if e != nil {
		return nil, ErrQueryTableInfo.Fmt(tableName).Wrap(e)
	}

	defer rows.Close()
	var cols []TableInfo

	for i := 0; rows.Next(); i++ {
		var ti TableInfo
		var defaultValue sql.NullString

		e := rows.Scan(
			&ti.Name,
			&ti.Type,
			&ti.NotNull,
			&defaultValue,
			&ti.PrimaryKey,
		)

		if e != nil {
			return nil, ErrScanTableInfo.Fmt(i, tableName).Wrap(e)
		}

		ti.HasDefaultValue = defaultValue.Valid
		ti.DefaultValue = defaultValue.String

		cols = append(cols, ti)
	}

	return cols, rows.Err()
}
