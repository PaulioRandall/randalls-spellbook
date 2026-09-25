package storm

import (
	"database/sql"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

type TableInfo struct {
	ColumnId     int
	Name         string
	Type         string
	Notnull      bool
	DefaultValue any
	PrimaryKey   any
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

func queryPragmaTableInfo(
	db *sql.DB,
	tableName string,
) ([]TableInfo, error) {

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

	var result []TableInfo
	var ti TableInfo

	for i := 0; rows.Next(); i++ {
		e := rows.Scan(
			&ti.ColumnId,
			&ti.Name,
			&ti.Type,
			&ti.Notnull,
			&ti.DefaultValue,
			&ti.PrimaryKey,
		)

		if e != nil {
			return nil, ErrScanTableInfo.Fmt(i, tableName).Wraps(e)
		}

		result = append(result, ti)
	}

	return result, rows.Err()
}
