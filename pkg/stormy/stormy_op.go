package stormy

import (
	"database/sql"

	"github.com/PaulioRandall/randalls-spellbook/pkg/wizzard"
)

// Operation is the function type used for [Stormy.Custom]
// operations.
type Operation[R any] func(OperationContext) (R, error)

// OperationContext provides access to object mapping, the
// database, and functions for managing model cache to
// execute queries. The context only lives while the
// operation is being executed; panic will ensue if any of
// its functions are called outside this time window.
type OperationContext interface {
	// Database returns the *sql.DB underpinning the [Stormy]
	// object.
	Database() *sql.DB

	// Map generates a model to support database operations.
	// The model will be added to the model cache if the
	// associated table exists within the database and
	// caching is enabled.
	Map(object any) (wizzard.Model, bool, error)

	// MapAs generates a model to support database operations
	// for an explicit table. The model will be added to the
	// model cache if the associated table exists within the
	// database and caching is enabled.
	MapAs(table string, object any) (wizzard.Model, bool, error)

	// CacheClear removes all models in the model cache.
	CacheClear()

	// CacheClearModel removes all models of the passed
	// object from the model cache.
	CacheClearModel(object any)

	// CacheClearTable removes all models associated with a
	// specific table from the model cache.
	CacheClearTable(name string)
}

type opCtx struct {
	st   *Stormy
	dead *bool
}

func (oc opCtx) panicIfDead() {
	if *oc.dead {
		panic("Operation context is dead! Context is only usable during an operation.")
	}
}

func (oc opCtx) Database() *sql.DB {
	oc.panicIfDead()
	return oc.st.db
}

func (oc opCtx) Map(object any) (wizzard.Model, bool, error) {
	oc.panicIfDead()
	return oc.st.mapModel(object)
}

func (oc opCtx) MapAs(table string, object any) (wizzard.Model, bool, error) {
	oc.panicIfDead()
	return oc.st.mapModelAs(table, object)
}

func (oc opCtx) CacheClear() {
	oc.panicIfDead()
	clear(oc.st.mapper)
}

func (oc opCtx) CacheClearModel(model any) {
	oc.panicIfDead()
	oc.st.mapper.ClearType(model)
}

func (oc opCtx) CacheClearTable(name string) {
	oc.panicIfDead()
	oc.st.mapper.ClearTable(name)
}

// Custom executes the passed operation. The Stormy object
// is locked during execution so no other named or custom
// operations may occur in parallel. If you value your
// sanity, do not close the database connection from within
// the operation.
func (st *Stormy) Custom[R any](op Operation[R]) (_ R, e error) {
	var empty R

	if e := st.setupOperation(); e != nil {
		return empty, e
	}
	defer st.tearDownOperation(&e)

	var dead bool
	defer func() {
		dead = true
	}()

	r, e := op(opCtx{
		st:   st,
		dead: &dead,
	})

	if e != nil {
		return empty, ErrForDatabase.Fmt(st.path).Wrap(e)
	}

	return r, nil
}

// setupOperation must always be called at the start of an
// operation.
func (st *Stormy) setupOperation() error {
	st.mutex.Lock()

	if st.IsOpen() {
		return nil
	}

	if st.openMode == OpenModeManual {
		st.mutex.Unlock()
		return st.errNotOpen()
	}

	if st.cacheMode == CacheModeRequest {
		clear(st.mapper)
	}

	switch st.openMode {
	case OpenModeRequest, OpenModePersist:
		st.autoOpened = true
		return st.open()
	default:
		panic("Sanity check! Unsupported OpenMode")
	}
}

// tearDownOperation must always be defer called straight
// after setupOperation.
func (st *Stormy) tearDownOperation(errPtr *error) {
	defer st.mutex.Unlock()

	if st.autoOpened && st.openMode == OpenModeRequest {
		e := st.close()

		// Don't hide the original error.
		if *errPtr == nil {
			*errPtr = e
		}
	}
}
