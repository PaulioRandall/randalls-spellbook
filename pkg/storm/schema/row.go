package schema

import (
	"database/sql"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

type ColumnSchema struct {
	ColumnId     int
	Name         string
	Type         string
	Notnull      bool
	DefaultValue any
	PrimaryKey   int
}

var (
	// ErrQueryTableInfo occurs when querying table_info
	// PRAGMA.
	ErrQueryTableInfo = curse.Proto(
		"Could not query table_info for table '%s'",
	)

	// ErrScanTableInfo occurs when scanning rows returned
	// from table_info PRAGMA.
	ErrScanTableInfo = curse.Proto(
		"Could not scan table_info row %d for table '%s'",
	)
)

func QueryColumns(
	db *sql.DB,
	tableName string,
) ([]ColumnSchema, error) {

	query := `
		SELECT
			cid,
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
		return nil, ErrQueryTableInfo.Fmt(tableName).Wraps(e)
	}
	defer rows.Close()

	var result []ColumnSchema
	var cs ColumnSchema

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(
			&cs.ColumnId,
			&cs.Name,
			&cs.Type,
			&cs.Notnull,
			&cs.DefaultValue,
			&cs.PrimaryKey,
		)

		if e != nil {
			return nil, ErrScanTableInfo.Fmt(i, tableName).Wraps(e)
		}

		result = append(result, cs)
	}

	return result, rows.Err()
}
