package storm

import (
	ref "reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/mapper"
)

type modTab = mapper.ModelTable
type modelCache map[modelCacheKey]modTab

type modelCacheKey struct {
	goType    ref.Type
	tableName string
}

func (mc modelCache) get(model any) (modTab, bool) {
	typ := ref.TypeOf(model)

	key := modelCacheKey{
		goType:    typ,
		tableName: typ.Name(),
	}

	table, found := mc[key]
	return table, found
}

func (mc modelCache) getForName(
	model any,
	tableName string,
) (modTab, bool) {
	key := modelCacheKey{
		goType:    ref.TypeOf(model),
		tableName: tableName,
	}

	table, found := mc[key]
	return table, found
}

func (mc modelCache) set(
	model any,
	table modTab,
) {
	key := modelCacheKey{
		goType:    ref.TypeOf(model),
		tableName: table.SqlName,
	}

	mc[key] = table
}
