package data

import (
	"database/sql"
	"fmt"
)

// sqlite_schema is the Go representation of sqlite_schema
// table within SQLite3. It is often referred to by its
// alias, sqlite_master.
//
// See https://sqlite.org/schematab.html
type sqlite_schema struct {
	// objectType as defined as 'type' in
	// https://sqlite.org/schematab.html. Because 'type' is a
	// reserved word in Go.
	//
	// "The sqlite_schema.type column will be one of the
	// following text strings: 'table', 'index', 'view', or
	// 'trigger' according to the type of object defined. The
	// 'table' string is used for both ordinary and virtual
	// tables."
	objectType string

	// name as defined in https://sqlite.org/schematab.html.
	//
	// "The sqlite_schema.name column will hold the name of
	// the object. UNIQUE and PRIMARY KEY constraints on
	// tables cause SQLite to create internal indexes with
	// names of the form "sqlite_autoindex_TABLE_N" where
	// TABLE is replaced by the name of the table that
	// contains the constraint and N is an integer beginning
	// with 1 and increasing by one with each constraint seen
	// in the table definition. In a WITHOUT ROWID table,
	// there is no sqlite_schema entry for the PRIMARY KEY,
	// but the "sqlite_autoindex_TABLE_N" name is set aside
	// for the PRIMARY KEY as if the sqlite_schema entry did
	// exist. This will affect the numbering of subsequent
	// UNIQUE constraints. The "sqlite_autoindex_TABLE_N"
	// name is never allocated for an INTEGER PRIMARY KEY,
	// either in rowid tables or WITHOUT ROWID tables."
	name string

	// tbl_name as defined in
	// https://sqlite.org/schematab.html.
	//
	// "The sqlite_schema.tbl_name column holds the name of a
	// table or view that the object is associated with. For
	// a table or view, the tbl_name column is a copy of the
	// name column. For an index, the tbl_name is the name of
	// the table that is indexed. For a trigger, the tbl_name
	// column stores the name of the table or view that
	// causes the trigger to fire."
	tbl_name string

	// rootpage as defined in
	// https://sqlite.org/schematab.html.
	//
	// "The sqlite_schema.rootpage column stores the page
	// number of the root b-tree page for tables and indexes.
	// For rows that define views, triggers, and virtual
	// tables, the rootpage column is 0 or NULL."
	rootpage int

	// sql as defined in https://sqlite.org/schematab.html.
	//
	// "The sqlite_schema.sql column stores SQL text that
	// describes the object. This SQL text is a CREATE TABLE,
	// CREATE VIRTUAL TABLE, CREATE INDEX, CREATE VIEW, or
	// CREATE TRIGGER statement that if evaluated against the
	// database file when it is the main database of a
	// database connection would recreate the object. The
	// text is usually a copy of the original statement used
	// to create the object but with normalizations applied
	// so that the text conforms to the following rules:
	// - The CREATE, TABLE, VIEW, TRIGGER, and INDEX keywords
	//   at the beginning of the statement are converted to
	//   all upper case letters.
	// - The TEMP or TEMPORARY keyword is removed if it
	//   occurs after the initial CREATE keyword. Any
	//   database name qualifier that occurs prior to the
	//   name of the object being created is removed.
	// - Leading spaces are removed.
	// - All spaces following the first two keywords are
	//   converted into a single space.
	// The text in the sqlite_schema.sql column is a copy of
	// the original CREATE statement text that created the
	// object, except normalized as described above and as
	// modified by subsequent ALTER TABLE statements. The
	// sqlite_schema.sql is NULL for the internal indexes
	// that are automatically created by UNIQUE or
	// PRIMARY KEY constraints."
	sql string
}

func listSqliteSchema(
	db *sql.DB,
) ([]sqlite_schema, error) {
	query := `
		SELECT
			type,
			name,
			tbl_name,
			rootpage,
			sql
		FROM
			sqlite_schema
	`

	rows, e := db.Query(query)
	if e != nil {
		return nil, fmt.Errorf(
			"Failed to query project table: %w",
			e,
		)
	}

	defer rows.Close()
	return parseSqliteSchemaTableRows(rows)
}

func parseSqliteSchemaTableRows(
	rows *sql.Rows,
) (
	[]sqlite_schema,
	error,
) {
	var result []sqlite_schema

	for rows.Next() {
		var ss sqlite_schema

		e := rows.Scan(
			&ss.objectType,
			&ss.name,
			&ss.tbl_name,
			&ss.rootpage,
			&ss.sql,
		)

		if e != nil {
			return nil, fmt.Errorf(
				"Failed to scan project table row: %w",
				e,
			)
		}

		result = append(result, ss)
	}

	return result, rows.Err()
}
