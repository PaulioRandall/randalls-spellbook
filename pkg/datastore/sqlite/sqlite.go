package sqlite

import (
	"database/sql"
	"fmt"
	"slices"

	_ "github.com/glebarez/go-sqlite"
)

var queryTablesExist string = `
	SELECT
		name
	FROM
		sqlite_master
	WHERE
		type='table'
`

var tableNames []string = []string{
	"media",
}

// applyQueries is a function that applies queries to a
// SQLite transaction.
type applyQueries func(tx *sql.Tx) error

// SQLiteDatabase satisfies the Datastore interface using
// SQLite3 as the underlying database.
type SQLiteDatabase struct {
	// path is the path to the SLQite file.
	path string

	// db is the database connection. It must be set to nil
	// when the database is closed. It is used to determine
	// if the database is open.
	db *sql.DB
}

// New creates a SQLite database that satisfies the
// Datastore interface. This function does not create or
// open the database.
func New(path string) *SQLiteDatabase {
	return &SQLiteDatabase{
		path: path,
	}
}

// Path returns the path to the SLQite file.
func (ds *SQLiteDatabase) Path() string {
	return ds.path
}

// Open satisfies the Datastore interface. It opens a local
// SQLite3 database, creating it if it doesn't exist yet.
// Tables are created if
func (ds *SQLiteDatabase) Open() error {
	db, e := sql.Open("sqlite", ds.path)

	if e != nil {
		return fmt.Errorf("Unable to open SQLite database: %w", e)
	}

	ds.db = db

	missingTables, e := listMissingTables(ds.db)
	if e != nil {
		return e
	}

	return createTables(ds.db, missingTables)
}

// IsOpen satisfies the Datastore interface.
func (ds *SQLiteDatabase) IsOpen() bool {
	return ds.db != nil
}

// Close satisfies the Datastore interface.
func (ds *SQLiteDatabase) Close() error {
	defer func() {
		ds.db = nil
	}()
	return ds.db.Close()
}

// listMissingTables returns all tables that don't exist
// within the database.
func listMissingTables(db *sql.DB) ([]string, error) {
	rows, e := db.Query(queryTablesExist)
	if e != nil {
		return nil, fmt.Errorf("Failed to query sqlite_master: %w", e)
	}

	fetchedNames, e := parseTableNameRows(rows)
	if e != nil {
		return nil, e
	}

	return listUniqueItems(tableNames, fetchedNames), nil
}

// listUniqueItems returns a slice of items present in
// sliceA but not in sliceB.
func listUniqueItems[T comparable](sliceA, sliceB []T) []T {
	result := []T{}

	for _, item := range sliceA {
		if !slices.Contains(sliceB, item) {
			result = append(result, item)
		}
	}

	return result
}

// contiansSubSlice returns true if all items in sliceB
// exist within sliceA, ignoring order. Ensure you order
// the arguements correctly; the sub slice must be the
// second argument.
func contiansSubSlice[T comparable](sliceA, sliceB []T) bool {
	for _, item := range sliceB {
		if !slices.Contains(sliceA, item) {
			return false
		}
	}
	return true
}

// createTables creates all tables within tableNames. If
// a table name is provided with no known query for
// creating the table then panic ensues.
func createTables(db *sql.DB, tableNames []string) error {
	applyQueriesToTx := func(tx *sql.Tx) error {
		var e error

		for _, name := range tableNames {
			switch name {
			case "media":
				e = applyCreateTableQueries(tx)
			default:
				panic("Missing handler for creating table")
			}

			if e != nil {
				return e
			}
		}

		return nil
	}

	return execTransaction(db, applyQueriesToTx)
}

// parseTableNameRows parses tables names from SQL query
// result rows.
func parseTableNameRows(rows *sql.Rows) ([]string, error) {
	var result []string

	for rows.Next() {
		var name string
		e := rows.Scan(&name)

		if e != nil {
			return nil, fmt.Errorf("Failed to scan table name row: %w", e)
		}

		result = append(result, name)
	}

	return result, rows.Err()
}

func execTransaction(db *sql.DB, doQueries applyQueries) error {
	tx, e := db.Begin()
	if e != nil {
		return fmt.Errorf("Failed to setup transaction: %w", e)
	}

	defer tx.Rollback()

	e = doQueries(tx)
	if e != nil {
		return e
	}

	e = tx.Commit()
	if e != nil {
		return fmt.Errorf("Failed to commit transaction: %w", e)
	}

	return nil
}

func applyCreateTableQueries(tx *sql.Tx) error {
	queries := []string{
		queryCreateMediaTable,
	}

	e := addQueriesToTransaction(tx, queries)
	if e != nil {
		return fmt.Errorf("Failed to create tables: %w", e)
	}

	return nil
}

func addQueriesToTransaction(tx *sql.Tx, queries []string) error {
	for _, q := range queries {
		if _, e := tx.Exec(q); e != nil {
			return e
		}
	}

	return nil
}
