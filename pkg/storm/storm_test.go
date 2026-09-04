package storm

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCheeseMaker struct {
	Id      int64
	Name    string
	Country string
}

type TestCheese struct {
	Id       int64
	MakerId  int64 // Map to CheeseMaker.Id
	Name     string
	Strength int64   // 1: Mild, 4: Extra Mature
	Rating   float64 // Out of 5
	Notes    string  // E.g. "Nutty after taste"
}

var bobs = TestCheeseMaker{
	Id:      1,
	Name:    "Bob's Cheeses",
	Country: "England",
}

var francs = TestCheeseMaker{
	Id:      2,
	Name:    "Franc's Fromage",
	Country: "France",
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

func openCreateAndInsertTestCheeseMakers(
	t *testing.T,
	db *Storm,
	objects ...TestCheeseMaker,
) {
	e := db.Open()
	require.Equal(t, nil, e)

	e = db.Create(TestCheeseMaker{})
	require.Equal(t, nil, e)

	for _, obj := range objects {
		e = db.Insert(obj)
		require.Equal(t, nil, e)
	}
}

func selectAllTestCheeseMakers(
	t *testing.T,
	db *Storm,
) []TestCheeseMaker {
	rows, e := db.db.Query(`
		SELECT
			Id,
			Name,
			Country
		FROM
			TestCheeseMaker
	`)
	require.Equal(t, nil, e)

	var result []TestCheeseMaker

	for rows.Next() {
		tcm := TestCheeseMaker{}
		e = rows.Scan(&tcm.Id, &tcm.Name, &tcm.Country)
		require.Equal(t, nil, e)
		result = append(result, tcm)
	}

	return result
}

func Test_Storm_Open_Close_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)
	require.Equal(t, false, db.IsOpen())

	e := db.Open()
	require.Equal(t, nil, e)

	defer func() {
		if r := recover(); r != nil {
			db.Close()
			panic(r)
		}
	}()

	require.Equal(t, true, db.IsOpen())

	e = db.Close()
	require.Equal(t, nil, e)
	require.Equal(t, false, db.IsOpen())
}

func Test_Storm_Create_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)

	e := db.Open()
	require.Equal(t, nil, e)
	defer db.Close()

	e = db.Create(TestCheeseMaker{})
	require.Equal(t, nil, e)

	e = db.Create(TestCheese{})
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

func Test_Storm_Insert_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)

	openCreateAndInsertTestCheeseMakers(
		t,
		db,
		bobs,
		francs,
	)

	defer db.Close()

	records := selectAllTestCheeseMakers(t, db)
	require.Equal(t, bobs, records[0])
	require.Equal(t, francs, records[1])
	require.Equal(t, 2, len(records))
}

func Test_Storm_Update_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)

	openCreateAndInsertTestCheeseMakers(
		t,
		db,
		bobs,
		francs,
	)

	defer db.Close()

	bobsUpdated := TestCheeseMaker{
		Id:      bobs.Id,
		Name:    bobs.Name,
		Country: "United Kingdom",
	}

	e := db.Update(bobsUpdated)
	require.Equal(t, nil, e)

	records := selectAllTestCheeseMakers(t, db)
	require.Equal(t, bobsUpdated, records[0])
	require.Equal(t, francs, records[1])
	require.Equal(t, 2, len(records))
}

func Test_Storm_SelectAll_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)

	openCreateAndInsertTestCheeseMakers(
		t,
		db,
		bobs,
		francs,
	)

	defer db.Close()

	result, e := db.SelectAll(TestCheeseMaker{})
	require.Equal(t, nil, e)

	act, ok := result.([]TestCheeseMaker)
	require.Equal(t, true, ok)
	require.Equal(t, 2, len(act))
	require.Equal(t, bobs, act[0])
	require.Equal(t, francs, act[1])
}

func Test_Storm_SelectById_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)

	openCreateAndInsertTestCheeseMakers(
		t,
		db,
		bobs,
		francs,
	)

	defer db.Close()

	result, e := db.SelectById(TestCheeseMaker{}, francs.Id)
	require.Equal(t, nil, e)

	act, ok := result.(TestCheeseMaker)
	require.Equal(t, true, ok)
	require.Equal(t, francs, act)
}

func Test_Storm_DeleteById_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)

	openCreateAndInsertTestCheeseMakers(
		t,
		db,
		bobs,
		francs,
	)

	defer db.Close()

	e := db.DeleteById(TestCheeseMaker{}, int64(1))
	require.Equal(t, nil, e)

	records := selectAllTestCheeseMakers(t, db)
	require.Equal(t, francs, records[0])
	require.Equal(t, 1, len(records))
}

func Test_Storm_Drop_1(t *testing.T) {
	dropTablesIfExists()
	db := New(testDb)

	e := db.Open()
	require.Equal(t, nil, e)
	defer db.Close()

	e = db.Create(TestCheeseMaker{})
	require.Equal(t, nil, e)

	e = db.Create(TestCheese{})
	require.Equal(t, nil, e)

	e = db.Drop(TestCheeseMaker{})
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
	require.Equal(t, "TestCheese", act)

	require.Equal(t, false, rows.Next())
}
