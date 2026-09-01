package sqlick

import (
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

func TestMain(m *testing.M) {
	deleteTestDbIfExists()
	code := m.Run()
	os.Exit(code)
}

func deleteTestDbIfExists() {
	e := os.Remove(testDb)
	if e != nil && !os.IsNotExist(e) {
		panic(e)
	}
}

func Test_Sqlick_AddStructTable_1(t *testing.T) {
	db := NewSqliteDB(testDb)

	e := db.AddStructTable(TestCheese{})

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
	db := NewSqliteDB(testDb)

	e := db.Open()
	require.Equal(t, nil, e)
	defer db.Close()

	e = db.AddStructTable(TestCheeseMaker{})
	require.Equal(t, nil, e)

	e = db.AddStructTable(TestCheese{})
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
	db := NewSqliteDB(testDb)

	bobs := TestCheeseMaker{
		Id:      1,
		Name:    "Bob's Cheese",
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

	e = db.AddStructTable(TestCheeseMaker{})
	require.Equal(t, nil, e)

	e = db.CreateTables()
	require.Equal(t, nil, e)

	e = db.Insert(bobs)
	require.Equal(t, nil, e)

	e = db.Insert(francs)
	require.Equal(t, nil, e)

	rows, e := db.db.Query(`
		SELECT
			Id,
			Name,
			Country
		FROM
			TestCheeseMaker
		WHERE
			Id IN (1, 2)
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
