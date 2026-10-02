package modtab

import (
	ref "reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Map_1(t *testing.T) {
	// When mapping a struct to a database table
	// if the table doesn't exist within the database
	// then the full unmodified model is returned

	type UnusedModel struct {
		Id     int
		Name   string
		Rating float64
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, exists, e := Map(db, UnusedModel{})
	exp := ModelTable{
		GoType:  ref.TypeOf(UnusedModel{}),
		GoName:  "UnusedModel",
		SqlName: "UnusedModel",
		Columns: []ModelColumn{
			TestIdCol,
			TestNameCol,
			TestRatingCol,
		},
	}

	require.NoError(t, e)
	require.Equal(t, false, exists)
	require.Equal(t, exp, act)
}

func Test_Map_2(t *testing.T) {
	// When mapping a struct to a database table
	// if the table exists in the database
	// and the model matches the table field-to-column
	// then the full model is returned

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, exists, e := Map(db, Player{})
	exp := ModelTable{
		GoType:  ref.TypeOf(Player{}),
		GoName:  "Player",
		SqlName: "Player",
		Columns: []ModelColumn{
			TestIdCol,
			TestNameCol,
			TestRatingCol,
		},
	}

	require.NoError(t, e)
	require.Equal(t, true, exists)
	require.Equal(t, exp, act)
}

func Test_Map_3(t *testing.T) {
	// When mapping a struct to a database table
	// if the table exists in the database
	// but some fields in the struct don't exist within
	// the table structure
	// then those fields are omitted from the model

	type Player struct {
		Id         int
		Name       string
		NotInTable float64
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, _, e := Map(db, Player{})
	exp := ModelTable{
		GoType:  ref.TypeOf(Player{}),
		GoName:  "Player",
		SqlName: "Player",
		Columns: []ModelColumn{
			TestIdCol,
			TestNameCol,
		},
	}

	require.NoError(t, e)
	require.Equal(t, exp, act)
}

func Test_Map_4(t *testing.T) {
	// When mapping a struct to a database table
	// if the table exists in the database
	// but the primary key column maps to a field that isn't
	// the first field within the struct (since the first
	// defaults to being the primary key) then the
	// parsed model is modified so that primary key field
	// that maps to the primary key column is set as the
	// primary key and all other fields are set as
	// not being primary keys

	type Player struct {
		Name   string // Default primary key
		Rating float64
		Id     int // Real primary key
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, _, e := Map(db, Player{})

	exp := ModelTable{
		GoType:  ref.TypeOf(Player{}),
		GoName:  "Player",
		SqlName: "Player",
		Columns: []ModelColumn{
			TestNameCol,
			TestRatingCol,
			TestIdCol,
		},
	}

	// The pre-created package level test columns won't
	// have the correct GoIndex because we messed up the
	// field ordering to test that ordering is accounted for
	// in the output, thus we need to correct the GoIndex so
	// the test assertion passes
	for i := range exp.Columns {
		exp.Columns[i].GoIndex = i
	}

	require.NoError(t, e)
	require.Equal(t, exp, act)
}

func Test_Map_5(t *testing.T) {
	// When mapping a struct to a database table
	// if the table exists in the database
	// but a field that maps to a column doesn't have a type
	// that is compatible with the SQL column's type
	// then a named error is returned

	type Player struct {
		Name int
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	_, _, e := Map(db, Player{})
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrTypeMismatch)
}
