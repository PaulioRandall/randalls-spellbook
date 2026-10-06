package stormy

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sqlick"

	"github.com/PaulioRandall/randalls-spellbook/pkg/wizzard"
)

// Table returns the full table details of the passed
// object's type name. The returned model is a container
// for information only and considered invalid for
// operations.
func (st *Stormy) Table(object any) (wizzard.Model, error) {
	return st.TableAs(typeName(object))
}

// TableAs is the same as [Stormy.Table] except the table
// name is provided explicitly.
func (st *Stormy) TableAs(table string) (wizzard.Model, error) {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	t, e := sqlick.QueryTable(st.db, table)

	if e != nil {
		return wizzard.Model{}, st.errForTable(table, e)
	}

	model := wizzard.Model{
		SqlName: t.Name,
		Props:   make([]wizzard.Property, len(t.Columns)),
	}

	for i, c := range t.Columns {
		model.Props[i] = wizzard.Property{
			SqlName:    c.Name,
			SqlType:    c.Type,
			SqlDefault: defaultGoTypeForSqlType(c.Type),
			IsKey:      c.PrimaryKey,
		}
	}

	return model, nil
}
