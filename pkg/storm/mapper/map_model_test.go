package mapper

import (
	"database/sql"
	ref "reflect"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"
)

func createAndPopulateTestDb(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	require.NoError(t, e)

	_, e = db.Exec(`
		CREATE TABLE Player (
			Id INTEGER,
			Name TEXT,
			Rating REAL,
			PRIMARY KEY (Id)
		)
	`)
	require.NoError(t, e)

	return db
}

var (
	TestIdCol = Column{
		GoType:     ref.TypeOf(int(0)),
		GoIndex:    0,
		GoName:     "Id",
		SqlType:    "INTEGER",
		SqlName:    "Id",
		SqlDefault: int(0),
		PrimaryKey: true,
	}

	TestNameCol = Column{
		GoType:     ref.TypeOf(""),
		GoIndex:    1,
		GoName:     "Name",
		SqlType:    "TEXT",
		SqlName:    "Name",
		SqlDefault: "''",
		PrimaryKey: false,
	}

	TestRatingCol = Column{
		GoType:     ref.TypeOf(float64(0)),
		GoIndex:    2,
		GoName:     "Rating",
		SqlType:    "REAL",
		SqlName:    "Rating",
		SqlDefault: float64(0),
		PrimaryKey: false,
	}
)

func Test_MapModel_1(t *testing.T) {
	// GIVEN model not in database
	// THEN  mapped Table is representing the full model

	type NotPlayer struct {
		Id     int
		Name   string
		Rating float64
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, exists, e := MapModel(db, NotPlayer{})
	exp := Table{
		GoType:  ref.TypeOf(NotPlayer{}),
		GoName:  "NotPlayer",
		SqlName: "NotPlayer",
		Columns: []Column{
			TestIdCol,
			TestNameCol,
			TestRatingCol,
		},
	}

	require.NoError(t, e)
	require.Equal(t, false, exists)
	require.Equal(t, exp, act)
}

func Test_MapModel_2(t *testing.T) {
	// GIVEN model in database with same fields/columns
	// THEN  mapped Table is representing the full model

	type Player struct {
		Id     int
		Name   string
		Rating float64
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, exists, e := MapModel(db, Player{})
	exp := Table{
		GoType:  ref.TypeOf(Player{}),
		GoName:  "Player",
		SqlName: "Player",
		Columns: []Column{
			TestIdCol,
			TestNameCol,
			TestRatingCol,
		},
	}

	require.NoError(t, e)
	require.Equal(t, true, exists)
	require.Equal(t, exp, act)
}

func Test_MapModel_3(t *testing.T) {
	// GIVEN model with exported field not in database
	// THEN  field is not omitted from Table

	type Player struct {
		Id         int
		Name       string
		NotInTable float64
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, _, e := MapModel(db, Player{})
	exp := Table{
		GoType:  ref.TypeOf(Player{}),
		GoName:  "Player",
		SqlName: "Player",
		Columns: []Column{
			TestIdCol,
			TestNameCol,
		},
	}

	require.NoError(t, e)
	require.Equal(t, exp, act)
}

func Test_MapModel_4(t *testing.T) {
	// GIVEN ID field is not first field in model
	// THEN  Table is returned with corrected ID field
	// AND   GoIndex matches model struct

	type Player struct {
		// Defaults to being ID column unless updated to match
		// database.
		Name   string
		Rating float64
		Id     int
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, _, e := MapModel(db, Player{})

	exp := Table{
		GoType:  ref.TypeOf(Player{}),
		GoName:  "Player",
		SqlName: "Player",
		Columns: []Column{
			TestNameCol,
			TestRatingCol,
			TestIdCol,
		},
	}

	// Because we reordered the fields within Players{}
	exp.Columns[0].GoIndex = 0
	exp.Columns[1].GoIndex = 1
	exp.Columns[2].GoIndex = 2

	require.NoError(t, e)
	require.Equal(t, exp, act)
}

func Test_MapModel_5(t *testing.T) {
	// GIVEN field with type that is not compatible with
	//       type in database table
	// THEN  panic ensues

	type Player struct {
		Name int
	}

	db := createAndPopulateTestDb(t)
	defer db.Close()

	_, _, e := MapModel(db, Player{})
	require.ErrorIs(t, e, ErrMapModel)
	require.ErrorIs(t, e, ErrFieldTypeMismatch)
}
