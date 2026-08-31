package sqlick

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

type Sqlick interface {
	Open() error
	IsOpen() bool
	Close() error
	AddStructTable(any) error
	CreateTables() error
}

type sqliteSqlick struct {
	path   string
	gen    SqlGenerator
	tables []Table
	db     *sql.DB
}

func NewSqliteDB(path string) *sqliteSqlick {
	return &sqliteSqlick{
		path: path,
		gen:  sqliteGenerator{},
	}
}

// Open satisfies the Sqlick interface.
func (ss *sqliteSqlick) Open() error {
	db, e := sql.Open("sqlite", ss.path)

	if e != nil {
		return fmt.Errorf(
			"Unable to open SQLite database: %w",
			e,
		)
	}

	ss.db = db
	return nil
}

// IsOpen satisfies the Sqlick interface.
func (ss *sqliteSqlick) IsOpen() bool {
	return ss.db != nil
}

// Close satisfies the Sqlick interface.
func (ss *sqliteSqlick) Close() error {
	if !ss.IsOpen() {
		return nil
	}

	defer func() {
		ss.db = nil
	}()

	return ss.db.Close()
}

// AddStructTable satisfies the Sqlick interface.
func (ss *sqliteSqlick) AddStructTable(object any) error {
	table, e := Parse(object)

	if e != nil {
		return fmt.Errorf("Failed to add struct table: %w", e)
	}

	ss.tables = append(ss.tables, table)
	return nil
}

// CreateTables satisfies the Sqlick interface.
func (ss *sqliteSqlick) CreateTables() error {
	for _, table := range ss.tables {
		e := ss.createTable(table)
		if e == nil {
			continue
		}

		return fmt.Errorf(
			"Failed to create table '%s': %w",
			table.GoName,
			e,
		)
	}

	return nil
}

func (ss *sqliteSqlick) createTable(table Table) error {
	query, e := ss.gen.TableCreate(table)
	if e != nil {
		return e
	}

	_, e = ss.db.Exec(query)
	return e
}
