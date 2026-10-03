package wizzard

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Map_1(t *testing.T) {
	// When mapping a struct to a database table
	// if the table doesn't exist within the database
	// then the full unmodified model is returned

	db := createTestDb(t)
	defer db.Close()

	act, exists, e := Map(db, Dummy{})
	exp := modelOfDummy()

	require.NoError(t, e)
	require.Equal(t, false, exists)
	require.Equal(t, exp, act)
}

func Test_Map_2(t *testing.T) {
	// When mapping a struct to a database table
	// if the table exists in the database
	// and the model matches the table field-to-column
	// then the full model is returned

	db := createDummyTestDb(t)
	defer db.Close()

	act, exists, e := Map(db, Dummy{})
	exp := modelOfDummy()

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

	type Dummy struct {
		Id         int
		Name       string
		NotInTable float64
	}

	db := createDummyTestDb(t)
	defer db.Close()

	act, _, e := Map(db, Dummy{})
	exp := modelOfDummy().Props
	exp = exp[:2] // Remove last property/field

	require.NoError(t, e)
	require.Equal(t, exp, act.Props)
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

	type Dummy struct {
		Name   string // Default primary key
		Rating float64
		Id     int // Real primary key
	}

	db := createDummyTestDb(t)
	defer db.Close()

	act, _, e := Map(db, Dummy{})
	exp := modelOfDummy().Props

	// Reorder props of default model.
	exp[0], exp[1], exp[2] = exp[1], exp[2], exp[0]
	for i := range exp {
		exp[i].GoIndex = i
	}

	require.NoError(t, e)
	require.Equal(t, exp, act.Props)
}

func Test_Map_5(t *testing.T) {
	// When mapping a struct to a database table
	// if the table exists in the database
	// but a field that maps to a column doesn't have a type
	// that is compatible with the SQL column's type
	// then a named error is returned

	type Dummy struct {
		Name int
	}

	db := createDummyTestDb(t)
	defer db.Close()

	_, _, e := Map(db, Dummy{})
	require.ErrorIs(t, e, ErrForTable)
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrTypeMismatch)
}

func Test_MapAs_1(t *testing.T) {
	// When mapping a struct to a database table
	// if the table doesn't exist within the database
	// then the full unmodified model is returned

	db := createTestDb(t)
	defer db.Close()

	act, exists, e := MapAs(db, "AliasName", Dummy{})
	exp := modelOfDummy()
	exp.SqlName = "AliasName"

	require.NoError(t, e)
	require.Equal(t, false, exists)
	require.Equal(t, exp, act)
}
