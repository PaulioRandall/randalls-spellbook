package storm

import (
// "github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

func (st *Storm) prepareTable(model any) error {
	tableName := typeName(model)

	metadata, e := querySqliteSchema(st.db, tableName)
	if e != nil {
		return e
	}

	if len(metadata) == 0 {
		e = st.createTable(model)
		if e != nil {
			return e
		}
	}

	tableInfo, e := queryPragmaTableInfo(st.db, tableName)
	if e != nil {
		return e
	}

	_ = tableInfo
	// NEXT: Filter tableInfo using the struct's fields so
	//       only columns that map to a field remain.
	// THEN: Design new struct types that hold info about
	//       the table and columns to be operated on.

	return nil
}
