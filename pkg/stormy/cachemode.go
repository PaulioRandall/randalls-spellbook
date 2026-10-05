package stormy

// CacheMode determines how and when parsed models are
// cached.
type CacheMode int

const (
	// CacheModeNone means caching is disabled.
	CacheModeNone CacheMode = iota

	// CacheModeRequest means the cache is reset for each
	// request. This is useful if the structure or existence
	// of tables is changed external to stormy but not during
	// operations like [Stormy.Put] which can accepts
	// heterogeneous data types but input is usually
	// homogeneous.
	CacheModeRequest

	// CacheModeSession means entries persist across the
	// session, but dropping a table will remove all entries
	// for that table. Cache will be cleared on database
	// close. Use this mode if the database is only written
	// to via a single Stormy object (which is the
	// most common scenario, thus this is the default mode).
	CacheModeSession
)

// CacheMode returns the current caching mode.
func (st *Stormy) CacheMode() CacheMode {
	return st.cacheMode
}

// SetCacheMode sets the caching mode. If the mode is
// already set then nothing happens, else the cache
// content is cleared before returning.
func (st *Stormy) SetCacheMode(mode CacheMode) {
	if st.cacheMode == mode {
		return
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.cacheMode = mode
	clear(st.mapper)
}

// CacheClear clears the cache regardless of caching mode.
func (st *Stormy) CacheClear() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	clear(st.mapper)
}

// CacheClearModel removes all cache entries associated
// with a specific model regardless of caching mode.
func (st *Stormy) CacheClearModel(model any) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.mapper.ClearType(model)
}

// CacheClearModel removes all cache entries associated
// with a specific table regardless of caching mode.
func (st *Stormy) CacheClearTable(name string) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.mapper.ClearTable(name)
}
