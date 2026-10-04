package storm

import (
	"database/sql"
	"net/url"
	"os"
	"path/filepath"
	ref "reflect"
	"strings"
	"sync"

	_ "github.com/glebarez/go-sqlite"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/wizzard"
)

// CacheMode represents an approach to caching parsed
// models.
type CacheMode int

const (
	// CacheModeNone means caching is disabled.
	CacheModeNone CacheMode = iota

	// CacheModeRequest means the cache is reset for each
	// request. This is useful if the structure or existence
	// of tables is changed through custom operations but
	// you are sure it will not change during calls to
	// Storm methods.
	CacheModeRequest

	// CacheModeSession means entries persist across the
	// session, but dropping a table will remove all entries
	// for that table. Cache will be cleared on database
	// close. Use this mode if the database is only written
	// to using a single Storm instance (which is the
	// most common scenario, thus this is the default mode).
	CacheModeSession
)

// Storm is the core type for interfacing with the
// database. Operations share a single mutex so only
// a single operation is permitted at once.
type Storm struct {
	path      string
	db        *sql.DB
	mutex     sync.Mutex
	cacheMode CacheMode
	mapper    wizzard.TableCache
}

// New returns a new [Storm] for the database represented
// by path.
func New(path string) *Storm {
	return &Storm{
		path:      path,
		cacheMode: CacheModeSession,
		mapper:    wizzard.TableCache{},
	}
}

// Open create a new Storm for the database represented
// by path, and opens it before returning.
func Open(path string) (*Storm, error) {
	st := &Storm{
		path:      path,
		cacheMode: CacheModeSession,
		mapper:    wizzard.TableCache{},
	}
	return st, st.Open()
}

// CacheMode returns the current caching mode.
func (st *Storm) CacheMode() CacheMode {
	return st.cacheMode
}

// SetCacheMode sets the caching mode. If the mode is
// already set then nothing happens, else the cache
// content is cleared before returning.
func (st *Storm) SetCacheMode(mode CacheMode) {
	if st.cacheMode == mode {
		return
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.cacheMode = mode
	clear(st.mapper)
}

// CacheClear clears the cache regardless of mode.
func (st *Storm) CacheClear() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	clear(st.mapper)
}

// CacheClearModel removes all cache entries associated
// with a specific model regardless of caching mode.
func (st *Storm) CacheClearModel(model any) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.mapper.ClearType(model)
}

// CacheClearModel removes all cache entries associated
// with a specific table regardless of caching mode.
func (st *Storm) CacheClearTable(name string) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.mapper.ClearTable(name)
}

// Open opens the database. If not an 'in-memory' path then
// the missing directories in the directory path are
// created.
func (st *Storm) Open() error {
	if st.IsOpen() {
		return nil
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

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

// IsOpen returns true if the database is open.
func (st *Storm) IsOpen() bool {
	return st.db != nil
}

// Close closes the database. Use Go's defer as usual.
// The cache content is always cleared on close regardless
// of caching mode.
func (st *Storm) Close() error {
	if !st.IsOpen() {
		return nil
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

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

// Table returns the full table details the passed object
// maps to. All columns in the table are included, not
// just those that are mapped.
func (st *Storm) Table(model any) (scumble.SqlTable, error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	return scumble.QueryTable(
		st.db,
		typeName(model),
	)
}

// TableAs is the same as [Storm.Table] except the table
// name is provided explicitly.
func (st *Storm) TableAs(table string) (scumble.SqlTable, error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	return scumble.QueryTable(
		st.db,
		table,
	)
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

func (st *Storm) prepareMapper() {
	if st.cacheMode == CacheModeRequest {
		clear(st.mapper)
	}
}

func (st *Storm) mapModel(
	model any,
) (wizzard.Model, bool, error) {
	if st.cacheMode == CacheModeNone {
		return wizzard.Map(st.db, model)
	}
	return st.mapper.Map(st.db, model)
}

func (st *Storm) mapModelAs(
	table string,
	model any,
) (wizzard.Model, bool, error) {
	if st.cacheMode == CacheModeNone {
		return wizzard.MapAs(st.db, table, model)
	}
	return st.mapper.MapAs(st.db, table, model)
}

func (st *Storm) errNotOpen() error {
	return ErrNotOpen.
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func (st *Storm) errForTable(
	table string,
	cause error,
) error {
	return ErrForTable.
		Fmt(table).
		Wrap(cause).
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func (st *Storm) errForModel(
	model any,
	cause error,
) error {
	if _, ok := model.(string); !ok {
		model = typeName(model)
	}

	return ErrForModel.
		Fmt(model).
		Wrap(cause).
		WrapIn(ErrForDatabase).
		Fmt(st.path)
}

func typeName(model any) string {
	return typeOf(model).Name()
}

func typeOf(model any) ref.Type {
	return ref.TypeOf(model)
}
