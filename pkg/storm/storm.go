package storm

import (
	"database/sql"
	"os"
	"path/filepath"
	ref "reflect"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

// TODO: Make Storm thread safe. Lock on function entry
//       and defer the unlock.
// TODO: Check the ordering of fields when scanning
//       database rows. Ensure the order the rows appear in
//       the select statement is the order they are
//       scanned to avoid the wrong info appearing in the
//       wrong struct fields.
// TODO: API call cache, create a ModelTable cache that
//       lives and dies in a single API call so bulk
//       inserts of the same kind don't fetch metadata we
//       already know will be the same.

var (
	// ErrNotOpen occurs when trying to perform an operation
	// before opening the database.
	ErrNotOpen = sin.Template(
		"Database not open",
	)

	// ErrTableRequest occurs within an error chain when any
	// error occurs involving a specific table/struct, except
	// for 'not found' errors.
	ErrTableRequest = sin.Template(
		"Request failed for table '%s'",
	)

	// ErrScanningRow is returned when an error occurs
	// scanning database results.
	ErrScanningRow = sin.Template(
		"Scanning row '%d'",
	)

	// ErrObjectNotFound is returned when a search for a
	// specific object/row failed.
	ErrObjectNotFound = sin.Template(
		"Object '%s' not found with ID '%v'",
	)

	// ErrDatabaseFile is returned when an error occurs with
	// or while opening or closing the database.
	ErrDatabaseFile = sin.Template(
		"Database IO error '%s'",
	)

	// ErrNoSuchTable is returned when an object is passed
	// to a function which does not have a registered table
	// for its type.
	ErrNoSuchTable = sin.Template(
		"No matching table for object type '%s'",
	)

	// ErrBadIdType is returned when an ID passed to a
	// function, e.g. SelectById, is not of the same type as
	// the ID field of the associated model type. This may
	// be returned even for compatible types like int when
	// int64 is expected.
	ErrBadIdType = sin.Template(
		"ID type mismatch for '%s', got %s, want %s",
	)
)

// Storm is the core type and the interface to the SQLite
// database.
type Storm struct {
	path string
	db   *sql.DB
}

// New returns a new [Storm] for the database represented
// by path.
func New(path string) *Storm {
	return &Storm{
		path: path,
	}
}

// Open opens the database. If not an 'in-memory' path then
// the missing directories in the directory path are
// created.
func (st *Storm) Open() error {
	if st.IsOpen() {
		return nil
	}

	e := makeParentDirs(st.path)
	if e != nil {
		return e
	}

	db, e := sql.Open("sqlite", st.path)
	if e != nil {
		return sin.Err("Unable to open SQLite database").
			Wrap(e).
			WrapIn(ErrDatabaseFile).
			Fmt(st.path)
	}

	st.db = db
	return nil
}

// IsOpen returns true if the database is open.
func (st *Storm) IsOpen() bool {
	return st.db != nil
}

// Close closes the database. Use with defer as usual.
func (st *Storm) Close() error {
	if !st.IsOpen() {
		return nil
	}

	defer func() {
		st.db = nil
	}()

	e := st.db.Close()
	if e == nil {
		return nil
	}

	return sin.Err("Unable to close SQLite database").
		Wrap(e).
		WrapIn(ErrDatabaseFile).
		Fmt(st.path)
}

// Table returns the full table details the passed model
// represents. All columns in the table are included, not
// just those that map to the passed model type.
func (st *Storm) Table(model any) (scumble.SqlTable, error) {
	return scumble.QueryTable(
		st.db,
		typeName(model),
	)
}

func makeParentDirs(path string) error {
	if path == ":memory" {
		// SQlite in-memory database. There is no path!
		return nil
	}

	parent := filepath.Dir(path)
	e := os.MkdirAll(parent, os.ModePerm)
	if e == nil {
		return nil
	}

	return sin.Err(
		"Unable to verify or create path to SQLite database",
	).Wrap(e)
}

func (st *Storm) getTable(model any) (ModelTable, bool, error) {
	return MapModel(st.db, model)
}

func (st *Storm) getOrCreateTable(model any) (ModelTable, error) {
	var zero ModelTable

	table, exists, e := MapModel(st.db, model)
	if e != nil {
		return zero, e
	}

	if exists {
		return table, nil
	}

	e = st.createTable(table)
	if e != nil {
		return zero, e
	}

	return table, nil
}

func typeName(model any) string {
	return typeOf(model).Name()
}

func typeOf(model any) ref.Type {
	return ref.TypeOf(model)
}
