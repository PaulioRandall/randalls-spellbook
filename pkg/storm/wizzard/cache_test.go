package wizzard

import (
	ref "reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CacheMapper_Map_1(t *testing.T) {
	// When mapping a struct to a database table
	// if the cache doesn't contain the model
	// and the database table doesn't currently exist
	// then the model will not be added to the cache

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	db := createTestDb(t)
	defer db.Close()

	mapper := CachedMapper{}
	_, _, e := mapper.Map(db, Player{})

	exp := CachedMapper{
		"Player": ModelCache{},
	}

	require.NoError(t, e)
	require.Equal(t, exp, mapper)
}

func Test_CacheMapper_Map_2(t *testing.T) {
	// When mapping a struct to a database table
	// if the cache already contains the model
	// and the database table already exists
	// then the model will be added to the cache

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	mapper := CachedMapper{}
	model, _, e := mapper.Map(db, Player{})

	exp := CachedMapper{
		"Player": ModelCache{
			ref.TypeOf(Player{}): model,
		},
	}

	require.NoError(t, e)
	require.Equal(t, exp, mapper)
}

func Test_CacheMapper_Map_3(t *testing.T) {
	// When mapping a struct to a database table
	// if the cache already contains the model
	// then the model will be sourced from the cache

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	db := createTestDb(t)
	defer db.Close()

	// Empty Model because the CachedMapper will never
	// store the an empty version.
	mapper := CachedMapper{
		"Player": ModelCache{
			ref.TypeOf(Player{}): Model{},
		},
	}

	act, _, _ := mapper.Map(db, Player{})
	require.Equal(t, Model{}, act)
}

func Test_CacheMapper_ClearType_1(t *testing.T) {
	// When clearing a type from the cache
	// if the cache contains a model associated with the type
	// then the related cache entries will be removed

	type Alpha struct{}
	type Beta struct{}
	type Charlie struct{}

	a1 := Model{SqlName: "Player"}
	b1 := Model{SqlName: "Player"}
	c1 := Model{SqlName: "Player"}

	a2 := Model{SqlName: "NotPlayer"}
	b2 := Model{SqlName: "NotPlayer"}
	c2 := Model{SqlName: "NotPlayer"}

	mapper := CachedMapper{
		"Player": ModelCache{
			ref.TypeOf(Alpha{}):   a1,
			ref.TypeOf(Beta{}):    b1,
			ref.TypeOf(Charlie{}): c1,
		},
		"NotPlayer": ModelCache{
			ref.TypeOf(Alpha{}):   a2,
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	exp := CachedMapper{
		"Player": ModelCache{
			ref.TypeOf(Beta{}):    b1,
			ref.TypeOf(Charlie{}): c1,
		},
		"NotPlayer": ModelCache{
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	mapper.ClearType(Alpha{})
	require.Equal(t, exp, mapper)
}

func Test_CacheMapper_ClearTable_1(t *testing.T) {
	// When clearing a table from the cache
	// if the cache contains a model associated with the
	// table that's associated with the type
	// then the related cache entries will be removed

	type Alpha struct{}
	type Beta struct{}
	type Charlie struct{}

	a1 := Model{SqlName: "Player"}
	b1 := Model{SqlName: "Player"}
	c1 := Model{SqlName: "Player"}

	a2 := Model{SqlName: "NotPlayer"}
	b2 := Model{SqlName: "NotPlayer"}
	c2 := Model{SqlName: "NotPlayer"}

	mapper := CachedMapper{
		"Player": ModelCache{
			ref.TypeOf(Alpha{}):   a1,
			ref.TypeOf(Beta{}):    b1,
			ref.TypeOf(Charlie{}): c1,
		},
		"NotPlayer": ModelCache{
			ref.TypeOf(Alpha{}):   a2,
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	exp := CachedMapper{
		"NotPlayer": ModelCache{
			ref.TypeOf(Alpha{}):   a2,
			ref.TypeOf(Beta{}):    b2,
			ref.TypeOf(Charlie{}): c2,
		},
	}

	mapper.ClearTable("Player")
	require.Equal(t, exp, mapper)
}
