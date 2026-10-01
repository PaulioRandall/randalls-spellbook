package storm

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/modtab"
)

// Create creates tables, represented by the passed models,
// within the database. If a table already exists then
// the model is ignored. It's safe to call Create at
// anytime, but doing it all upfront is can help avoid
// debugging issues.
//
// You are responsible for choosing appropriate data types
// for your structures. If you choose to use an unsigned
// integer type (unit etc) or shorter integer type
// (int8 etc) then you are responsible for ensuring the
// number polarity or maximum value (respectively) within
// the models and tables. I recommend sticking to int or
// int64 as these are the least issue prone.
//
// # Table creation rules
//
//   - Model (struct) name becomes the table name.
//   - The exported model fields become columns.
//   - Field name is the column name.
//   - Field type is mapped to a SQLite column type.
//   - Only primitive types may be used as field types.
//   - All columns have NOT NULL constraint.
//   - All columns have DEFAULT set to the zero value of
//     the field type.
//   - By default, the first field in the model is
//     designated the PRIMARY KEY.
//
// # Type mappings
//
//	INTEGER:
//		int, int8, int16, int32, int64,
//		uint, uint8, uint16, uint32, uint64
//	REAL:
//		float32, float64
//	TEXT:
//		string
func (st *Storm) Create(models ...any) error {
	if !st.IsOpen() {
		return st.errNotOpen()
	}

	mapper := modtab.CachedMapper{}

	for _, m := range models {
		table, exists, e := mapper.Map(st.db, m)
		if e != nil {
			return st.errForModel(m, e)
		}

		if exists {
			continue
		}

		e = table.Create(st.db)
		if e != nil {
			return st.errForModel(m, e)
		}
	}

	return nil
}
