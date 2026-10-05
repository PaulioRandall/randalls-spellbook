package stormy

// OpenMode determines when automatic database opening and
// closing is allowed to take place. Remember that
// closing an in-memory database (or the last shared
// connection to one) will result in all data being lost.
type OpenMode int

const (
	// OpenModeManual means no automatic opening or closing
	// will occur.
	OpenModeManual OpenMode = iota

	// OpenModeRequest means that if the database is closed
	// and an operation is called, e.g. Create, then the
	// database will open, perform the operation, then close.
	// However, if the database was open prior to the request
	// then it will remain open after the operation.
	OpenModeRequest

	// OpenModePersist means that if the database is closed
	// and an operation is called, e.g. Create, then the
	// database will open, perform the operation, then remain
	// open. Closing must be performed manually. This is the
	// default mode.
	OpenModePersist
)

// OpenMode returns the current automated database opening
// and closing strategy.
func (st *Stormy) OpenMode() OpenMode {
	return st.openMode
}

// SetOpenMode sets the automated database opening and
// closing strategy. The act of setting the mode will not
// result in opening or closing the database.
func (st *Stormy) SetOpenMode(mode OpenMode) {
	if st.openMode == mode {
		return
	}

	st.mutex.Lock()
	defer st.mutex.Unlock()

	st.openMode = mode
}
