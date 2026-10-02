package modtab

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

	exp := CachedMapper{}

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
		ref.TypeOf(Player{}): model,
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
		ref.TypeOf(Player{}): Model{},
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

	a := Model{SqlName: "Player"}
	b := Model{SqlName: "Player"}
	c := Model{SqlName: "NotPlayer"}

	mapper := CachedMapper{
		ref.TypeOf(Alpha{}):   a,
		ref.TypeOf(Beta{}):    b,
		ref.TypeOf(Charlie{}): c,
	}

	exp := CachedMapper{
		ref.TypeOf(Beta{}):    b,
		ref.TypeOf(Charlie{}): c,
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

	a := Model{SqlName: "Player"}
	b := Model{SqlName: "Player"}
	c := Model{SqlName: "NotPlayer"}

	mapper := CachedMapper{
		ref.TypeOf(Alpha{}):   a,
		ref.TypeOf(Beta{}):    b,
		ref.TypeOf(Charlie{}): c,
	}

	exp := CachedMapper{
		ref.TypeOf(Charlie{}): c,
	}

	mapper.ClearTable("Player")
	require.Equal(t, exp, mapper)
}
