package storm

import (
	ref "reflect"
)

type modTab = ModelTable
type modelCache map[modelCacheKey]modTab
type modelCacheKey struct {
	goType    ref.Type
	tableName string
}

func (mc modelCache) get(model any) (modTab, bool) {
	typ := typeOf(model)
	return mc.getForName(model, typ.Name())
}

func (mc modelCache) getForName(
	model any,
	tableName string,
) (modTab, bool) {
	key := modelCacheKey{
		goType:    typeOf(model),
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
		goType:    typeOf(model),
		tableName: table.SqlName,
	}

	mc[key] = table
}

func (mc modelCache) clear() {
	clear(mc)
}

func (mc modelCache) clearTable(tableName string) {
	for key := range mc {
		if key.tableName == tableName {
			delete(mc, key)
		}
	}
}

func (mc modelCache) clearModel(model any) {
	typ := typeOf(model)

	for key := range mc {
		if key.goType == typ {
			delete(mc, key)
		}
	}
}

func (mc modelCache) setTestEntry(model any, tableName string) {
	mc.set(model, modTab{
		SqlName: tableName,
	})
}
