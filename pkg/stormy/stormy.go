package stormy

import (
	"database/sql"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"

	"github.com/PaulioRandall/randalls-spellbook/pkg/stormy/wizzard"
)

// Stormy is the core type for interfacing with the
// database. Operations share a single mutex so only
// a single operation is permitted at once.
type Stormy struct {
	path       string
	db         *sql.DB
	mutex      sync.Mutex
	openMode   OpenMode
	autoOpened bool
	cacheMode  CacheMode
	mapper     wizzard.TableCache
}

// New returns a new [Stormy] object for the database
// represented by path.
func New(path string) *Stormy {
	return &Stormy{
		path:      path,
		openMode:  OpenModePersist,
		cacheMode: CacheModeSession,
		mapper:    wizzard.TableCache{},
	}
}

// Open creates a new [Stormy] object for the database
// represented by path, and opens it before returning.
func Open(path string) (*Stormy, error) {
	st := New(path)
	return st, st.Open()
}

// Open opens the database. If not an 'in-memory' path then
// the missing directories in the directory path are
// created.
func (st *Stormy) Open() error {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.autoOpened = false
	return st.open()
}

// IsOpen returns true if the database is open.
func (st *Stormy) IsOpen() bool {
	return st.db != nil
}

// Database returns *sql.DB underpinning this Stormy
// instance. Operations performed directly on the sql.DB
// will not benefit from internal synchronisation and other
// safe guards. Use with care.
func (st *Stormy) Database() *sql.DB {
	return st.db
}

// Close closes the database. Use Go's defer as usual.
// The cache content is always cleared on close regardless
// of caching mode.
func (st *Stormy) Close() error {
	st.mutex.Lock()
	defer st.mutex.Unlock()
	return st.close()
}

func (st *Stormy) open() error {
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
			WrapIn(ErrForDatabase).
			Fmt(st.path)
	}

	st.db = db
	return nil
}

func (st *Stormy) close() error {
	if !st.IsOpen() {
		return nil
	}

	clear(st.mapper)

	defer func() {
		st.db = nil
	}()

	e := st.db.Close()
	if e == nil {
		return nil
	}

	return sin.Err("Unable to close SQLite database").
		Wrap(e).
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func defaultGoTypeForSqlType(sqlType string) any {
	switch sqlType {
	case "INTEGER":
		return int(0)
	case "REAL":
		return float64(0)
	case "TEXT":
		return string("")
	default:
		panic("Unsupport SQL type: " + sqlType)
	}
}

func (st *Stormy) mapModel(
	model any,
) (wizzard.Model, bool, error) {
	if st.cacheMode == CacheModeNone {
		return wizzard.Map(st.db, model)
	}
	return st.mapper.Map(st.db, model)
}

func (st *Stormy) mapModelAs(
	table string,
	model any,
) (wizzard.Model, bool, error) {
	if st.cacheMode == CacheModeNone {
		return wizzard.MapAs(st.db, table, model)
	}
	return st.mapper.MapAs(st.db, table, model)
}

func (st *Stormy) errNotOpen() error {
	return ErrNotOpen.
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func (st *Stormy) errForTable(
	table string,
	cause error,
) error {
	return ErrForTable.
		Fmt(table).
		Wrap(cause).
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func (st *Stormy) errForType(
	object any,
	cause error,
) error {
	if _, ok := object.(string); !ok {
		object = typeName(object)
	}

	return ErrForType.
		Fmt(object).
		Wrap(cause).
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func (st *Stormy) errForModel(
	model wizzard.Model,
	cause error,
) error {
	return ErrForType.
		Fmt(model.GoName).
		Wrap(cause).
		WrapIn(ErrForTable).
		Fmt(model.SqlName).
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func makeParentDirs(path string) error {
	if isInMemoryDatabase(path) {
		// There is no path!
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

// isInMemoryDatabase determines if the path will open an
// in-memory database. There are probably a few very
// uncommon and highly niche edge cases that are not
// covered. Tough! I CBA to deal with them and don't
// currently trust Claudes output on the matter of
// detecting in-memory database paths.
//
// The path is being used like: sql.Open("sqlite", path).
func isInMemoryDatabase(path string) bool {
	name, _, _ := strings.Cut(path, "?")
	if name == ":memory:" || name == "file::memory:" {
		return true
	}

	if !strings.HasPrefix(path, "file:") {
		// Can't be in-memory if it doesn't use the file scheme
		// or is not the special filename ":memory:".
		return false
	}

	u, e := url.Parse(path)
	if e != nil {
		// Fails to parse then assume it's not in-memory.
		return false
	}

	if u.Query().Get("mode") == "memory" {
		// Named in-memory database.
		return true
	}

	if u.Query().Get("vfs") == "memdb" {
		// Alternative way to specify in-memory.
		return true
	}

	// Finally, handle encoded path.
	name = u.Opaque
	if name == "" {
		name = u.Path
	}

	name, _ = url.PathUnescape(name)
	return name == ":memory:"
}

func typeName(model any) string {
	return reflect.TypeOf(model).Name()
}
