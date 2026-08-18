package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/entity"
)

type applyQueries func(tx *sql.Tx) error

// Datastore is an interface for accessing data within the
// project. It may be backed by any form of data storage,
// local or remote.
type Datastore interface {
	// Path returns the URL to the datastore.
	Path() string

	// Open opens the datastore, creating it if it doesn't
	// already exist.
	Open() error

	// IsOpen returns true if the datastore is open.
	IsOpen() error

	// Close closes the datastore and cleans up all
	// resources.
	Close() error

	// GetMedia gets the media with the given EntityID.
	GetMedia(entityId entity.EntityId) (entity.Media, error)
}

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
func NewSQLiteDatabase(path string) *SQLiteDatabase {
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
	db, e := sql.Open("sqlite", path)

	if e != nil {
		return fmt.Errorf("Unable to open SQLite database: %w", e)
	}

	ds.db = db
	return execTransaction(db, applyCreateTableQueries)
}

// IsOpen satisfies the Datastore interface.
func (ds *SQLiteDatabase) IsOpen() error {
	return ds.db != nil
}

// Close satisfies the Datastore interface.
func (ds *SQLiteDatabase) Close() error {
	defer func() {
		ds.db = nil
	}()
	ds.db.Close()
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
	for q := range queries {
		if _, e := tx.Exec(q); e != nil {
			return e
		}
	}

	return nil
}

/*
// TODO: generated code for reference.

// insertRecords inserts 5 sample rows into the "users" table.
func insertRecords(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO users (id, name) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	names := []string{"Alice", "Bob", "Carol", "Dave", "Eve"}
	for i, name := range names {
		id := i + 1
		if _, err := stmt.Exec(id, name); err != nil {
			return fmt.Errorf("insert record %d: %w", id, err)
		}
	}
	return nil
}

// queryRecordsByIDs retrieves rows whose id is in the given list of IDs.
func queryRecordsByIDs(db *sql.DB, ids []int) ([]string, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Build a "?,?,?" placeholder string matching the number of IDs.
	placeholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	query := fmt.Sprintf("SELECT id, name FROM users WHERE id IN (%s)", placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query records: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		results = append(results, fmt.Sprintf("%d: %s", id, name))
	}
	return results, rows.Err()
}

func main() {
	db, err := createDatabase("example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatal(err)
	}

	if err := insertRecords(db); err != nil {
		log.Fatal(err)
	}

	results, err := queryRecordsByIDs(db, []int{2, 4})
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		fmt.Println(r)
	}
}

*/
