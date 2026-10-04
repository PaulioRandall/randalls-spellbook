package stormy

import (
	"database/sql"

	"github.com/PaulioRandall/randalls-spellbook/pkg/stormy/wizzard"
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
func (st *Stormy) Custom[R any](op Operation[R]) (R, error) {
	var empty R

	if !st.IsOpen() {
		return empty, st.errNotOpen()
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.prepareMapper()

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
