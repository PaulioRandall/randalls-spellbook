package wizzard

import (
	ref "reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TableCache_Map_1(t *testing.T) {
	// When mapping a struct to a database table
	// if the cache doesn't contain the model
	// and the database table doesn't currently exist
	// then the model will not be added to the cache

	db := createTestDb(t)
	defer db.Close()

	cache := TableCache{}
	_, _, e := cache.Map(db, Dummy{})

	exp := TableCache{
		"Dummy": ModelCache{},
	}

	require.NoError(t, e)
	require.Equal(t, exp, cache)
}

func Test_TableCache_Map_2(t *testing.T) {
	// When mapping a struct to a database table
	// if the cache already contains the model
	// and the database table already exists
	// then the model will be added to the cache

	db := createDummyTestDb(t)
	defer db.Close()

	cache := TableCache{}
	model, _, e := cache.Map(db, Dummy{})

	exp := TableCache{
		"Dummy": ModelCache{
			ref.TypeOf(Dummy{}): model,
		},
	}

	require.NoError(t, e)
	require.Equal(t, exp, cache)
}

func Test_TableCache_Map_3(t *testing.T) {
	// When mapping a struct to a database table
	// if the cache already contains the model
	// then the model will be sourced from the cache

	db := createTestDb(t)
	defer db.Close()

	// Empty Model because the TableCache will never
	// store the an empty version.
	cache := TableCache{
		"Dummy": ModelCache{
			ref.TypeOf(Dummy{}): Model{},
		},
	}

	act, _, _ := cache.Map(db, Dummy{})
	require.Equal(t, Model{}, act)
}

func Test_TableCache_ClearType_1(t *testing.T) {
	// When clearing a type from the cache
	// if the cache contains a model associated with the type
	// then the related cache entries will be removed

	type Alpha struct{}
	type Beta struct{}
	type Charlie struct{}

	a1 := Model{SqlName: "Dummy"}
	b1 := Model{SqlName: "Dummy"}
	c1 := Model{SqlName: "Dummy"}

	a2 := Model{SqlName: "NotDummy"}
	b2 := Model{SqlName: "NotDummy"}
	c2 := Model{SqlName: "NotDummy"}

	cache := TableCache{
		"Dummy": ModelCache{
			ref.TypeOf(Alpha{}):   a1,
			ref.TypeOf(Beta{}):    b1,
			ref.TypeOf(Charlie{}): c1,
		},
		"NotDummy": ModelCache{
			ref.TypeOf(Alpha{}):   a2,
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	exp := TableCache{
		"Dummy": ModelCache{
			ref.TypeOf(Beta{}):    b1,
			ref.TypeOf(Charlie{}): c1,
		},
		"NotDummy": ModelCache{
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	cache.ClearType(Alpha{})
	require.Equal(t, exp, cache)
}

func Test_TableCache_ClearTable_1(t *testing.T) {
	// When clearing a table from the cache
	// if the cache contains a model associated with the
	// table that's associated with the type
	// then the related cache entries will be removed

	type Alpha struct{}
	type Beta struct{}
	type Charlie struct{}

	a1 := Model{SqlName: "Dummy"}
	b1 := Model{SqlName: "Dummy"}
	c1 := Model{SqlName: "Dummy"}

	a2 := Model{SqlName: "NotDummy"}
	b2 := Model{SqlName: "NotDummy"}
	c2 := Model{SqlName: "NotDummy"}

	cache := TableCache{
		"Dummy": ModelCache{
			ref.TypeOf(Alpha{}):   a1,
			ref.TypeOf(Beta{}):    b1,
			ref.TypeOf(Charlie{}): c1,
		},
		"NotDummy": ModelCache{
			ref.TypeOf(Alpha{}):   a2,
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	exp := TableCache{
		"NotDummy": ModelCache{
			ref.TypeOf(Alpha{}):   a2,
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	cache.ClearTable("Dummy")
	require.Equal(t, exp, cache)
}
