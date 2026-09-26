package storm

import (
	"errors"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/modelparser"
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
)

func (st *Storm) prepareTable(model any) error {
	structModel := modelparser.Parse(model)
	tableSchema, e := schema.Query(st.db, structModel.SqlName)

	if errors.Is(e, schema.ErrTableNotFound) {
		return st.createTable(model)
	}

	if e != nil {
		return e
	}

	_ = structModel
	_ = tableSchema
	// NEXT: Filter tableInfo using the struct's fields so
	//       only columns that map to a field remain.
	// THEN: Design new struct types that hold info about
	//       the table and columns to be operated on.

	return nil
}

// diffTable returns the differences between a table model
// and the table and row schemas.
func diffTable(model modelparser.Table2) {

}
