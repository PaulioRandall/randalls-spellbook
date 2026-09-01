package sqlick

import (
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCheeseMaker struct {
	Id      int
	Name    string
	Country string
}

type TestCheese struct {
	Id       int
	MakerId  int // Map to CheeseMaker.Id
	Name     string
	Strength int     // 1: Mild, 4: Extra Mature
	Rating   float64 // Out of 5
	Notes    string  // E.g. "Nutty after taste"
}

const testDb string = "../../testdata/test.sqlite"

func dropTablesIfExists() {
	_, e := os.Stat(testDb)
	if os.IsNotExist(e) {
		return
	}

	db, e := sql.Open("sqlite", testDb)
	panicOnError(e)
	defer db.Close()

	_, e = db.Exec("DROP TABLE IF EXISTS TestCheeseMaker")
	panicOnError(e)

	_, e = db.Exec("DROP TABLE IF EXISTS TestCheese")
	panicOnError(e)
}

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func Test_Sqlick_Register_1(t *testing.T) {
	dropTablesIfExists()
	db := NewSqliteDB(testDb)

	e := db.Register(TestCheese{})

	exp := Table{
		GoType: reflect.TypeOf(TestCheese{}),
		GoName: "TestCheese",
		Columns: []Column{
			Column{
				GoName: "Id",
				GoType: "int",
			},
			Column{
				GoName: "MakerId",
				GoType: "int",
			},
			Column{
				GoName: "Name",
				GoType: "string",
			},
			Column{
				GoName: "Strength",
				GoType: "int",
			},
			Column{
				GoName: "Rating",
				GoType: "float64",
			},
			Column{
				GoName: "Notes",
				GoType: "string",
			},
		},
	}

	require.Equal(t, nil, e)
	require.Equal(t, 1, len(db.tables))
	require.Equal(t, exp, db.tables[0])
}

func Test_Sqlick_Open_Close_1(t *testing.T) {
	dropTablesIfExists()
	db := NewSqliteDB(testDb)

	require.Equal(t, false, db.IsOpen())

	e := db.Open()
	require.Equal(t, nil, e)
	require.Equal(t, true, db.IsOpen())

	e = db.Close()
	require.Equal(t, nil, e)
	require.Equal(t, false, db.IsOpen())
}

func Test_Sqlick_CreateTables_1(t *testing.T) {
	dropTablesIfExists()
	db := NewSqliteDB(testDb)

	e := db.Open()
	require.Equal(t, nil, e)
	defer db.Close()

	e = db.Register(TestCheeseMaker{})
	require.Equal(t, nil, e)

	e = db.Register(TestCheese{})
	require.Equal(t, nil, e)

	e = db.CreateTables()
	require.Equal(t, nil, e)

	rows, e := db.db.Query(`
		SELECT
			name
		FROM
			sqlite_schema
		WHERE
			name IN (
				'TestCheeseMaker',
				'TestCheese'
			)
	`)
	require.Equal(t, nil, e)

	var act string

	require.Equal(t, true, rows.Next())
	e = rows.Scan(&act)
	require.Equal(t, nil, e)
	require.Equal(t, "TestCheeseMaker", act)

	require.Equal(t, true, rows.Next())
	e = rows.Scan(&act)
	require.Equal(t, nil, e)
	require.Equal(t, "TestCheese", act)

	require.Equal(t, false, rows.Next())
}

func Test_Sqlick_Insert_1(t *testing.T) {
	dropTablesIfExists()
	db := NewSqliteDB(testDb)

	bobs := TestCheeseMaker{
		Id:      1,
		Name:    "Bob's Cheeses",
		Country: "England",
	}

	francs := TestCheeseMaker{
		Id:      2,
		Name:    "Franc's Fromage",
		Country: "France",
	}

	e := db.Open()
	require.Equal(t, nil, e)
	defer db.Close()

	e = db.Register(TestCheeseMaker{})
	require.Equal(t, nil, e)

	e = db.CreateTables()
	require.Equal(t, nil, e)

	// Function under test.
	e = db.Insert(bobs)
	require.Equal(t, nil, e)

	// Function under test.
	e = db.Insert(francs)
	require.Equal(t, nil, e)

	rows, e := db.db.Query(`
		SELECT
			Id,
			Name,
			Country
		FROM
			TestCheeseMaker
	`)
	require.Equal(t, nil, e)

	var act TestCheeseMaker

	require.Equal(t, true, rows.Next())
	e = rows.Scan(&act.Id, &act.Name, &act.Country)
	require.Equal(t, nil, e)
	require.Equal(t, bobs, act)

	require.Equal(t, true, rows.Next())
	e = rows.Scan(&act.Id, &act.Name, &act.Country)
	require.Equal(t, nil, e)
	require.Equal(t, francs, act)

	require.Equal(t, false, rows.Next())
}

func Test_Sqlick_Update_1(t *testing.T) {
	dropTablesIfExists()
	db := NewSqliteDB(testDb)

	bobs := TestCheeseMaker{
		Id:      1,
		Name:    "Bob's Cheeses",
		Country: "England",
	}

	e := db.Open()
	require.Equal(t, nil, e)
	defer db.Close()

	e = db.Register(TestCheeseMaker{})
	require.Equal(t, nil, e)

	e = db.CreateTables()
	require.Equal(t, nil, e)

	e = db.Insert(bobs)
	require.Equal(t, nil, e)

	bobsUpdated := TestCheeseMaker{
		Id:      bobs.Id,
		Name:    bobs.Name,
		Country: "United Kingdom",
	}

	// Function under test.
	e = db.Update(bobsUpdated)
	require.Equal(t, nil, e)

	rows, e := db.db.Query(`
		SELECT
			Id,
			Name,
			Country
		FROM
			TestCheeseMaker
	`)
	require.Equal(t, nil, e)

	var act TestCheeseMaker

	require.Equal(t, true, rows.Next())
	e = rows.Scan(&act.Id, &act.Name, &act.Country)
	require.Equal(t, nil, e)
	require.Equal(t, bobsUpdated, act)

	require.Equal(t, false, rows.Next())
}
