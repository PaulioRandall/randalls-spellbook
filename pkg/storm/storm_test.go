package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCheeseMaker struct {
	Id      int64
	Name    string
	Country string
}

type testCheese struct {
	Id       int64
	MakerId  int64
	Name     string
	Strength int64
	Rating   float64
	Notes    string
}

var bobs = testCheeseMaker{
	Id:      1,
	Name:    "Bob's Cheeses",
	Country: "England",
}

var francs = testCheeseMaker{
	Id:      2,
	Name:    "Franc's Fromage",
	Country: "France",
}

func openCreateInsert(
	t *testing.T,
	models []any,
	makers []testCheeseMaker,
) *Storm {
	db := New(":memory:")

	e := db.Open()
	require.Equal(t, nil, e)

	defer func() {
		if r := recover(); r != nil {
			db.Close()
		}
	}()

	for _, m := range models {
		e = db.Create(m)
		require.Equal(t, nil, e)
	}

	for _, m := range makers {
		e = db.Insert(m)
		require.Equal(t, nil, e)
	}

	return db
}

func selectAllTestCheeseMakers(
	t *testing.T,
	db *Storm,
) []testCheeseMaker {
	rows, e := db.db.Query(`
		SELECT
			Id,
			Name,
			Country
		FROM
			testCheeseMaker
	`)
	require.Equal(t, nil, e)

	var result []testCheeseMaker

	for rows.Next() {
		tcm := testCheeseMaker{}
		e = rows.Scan(&tcm.Id, &tcm.Name, &tcm.Country)
		require.Equal(t, nil, e)
		result = append(result, tcm)
	}

	return result
}

func selectTestTableNamesFromSqliteSchema(
	t *testing.T,
	db *Storm,
) []string {
	rows, e := db.db.Query(`
		SELECT
			name
		FROM
			sqlite_schema
		WHERE
			name IN (
				'testCheeseMaker',
				'testCheese'
			)
	`)
	require.Equal(t, nil, e)

	var result []string

	for rows.Next() {
		var name string
		e = rows.Scan(&name)
		require.Equal(t, nil, e)
		result = append(result, name)
	}

	return result
}

func Test_Storm_Open_Close_1(t *testing.T) {
	db := openCreateInsert(
		t,
		nil,
		nil,
	)
	defer db.Close()

	require.Equal(t, true, db.IsOpen())

	e := db.Close()
	require.Equal(t, nil, e)
	require.Equal(t, false, db.IsOpen())
}

func Test_Storm_Create_1(t *testing.T) {
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		nil,
	)
	defer db.Close()

	tableNames := selectTestTableNamesFromSqliteSchema(t, db)

	require.Equal(t, "testCheeseMaker", tableNames[0])
	require.Equal(t, "testCheese", tableNames[1])
	require.Equal(t, 2, len(tableNames))
}

func Test_Storm_Insert_1(t *testing.T) {
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		[]testCheeseMaker{bobs, francs},
	)
	defer db.Close()

	records := selectAllTestCheeseMakers(t, db)
	require.Equal(t, bobs, records[0])
	require.Equal(t, francs, records[1])
	require.Equal(t, 2, len(records))
}

func Test_Storm_Update_1(t *testing.T) {
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		[]testCheeseMaker{bobs, francs},
	)
	defer db.Close()

	bobsUpdated := testCheeseMaker{
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
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		[]testCheeseMaker{bobs, francs},
	)
	defer db.Close()

	records, e := db.SelectAll(testCheeseMaker{})
	require.Equal(t, nil, e)

	act, ok := records.([]testCheeseMaker)
	require.Equal(t, true, ok)
	require.Equal(t, 2, len(act))
	require.Equal(t, bobs, act[0])
	require.Equal(t, francs, act[1])
}

func Test_Storm_SelectById_1(t *testing.T) {
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		[]testCheeseMaker{bobs, francs},
	)
	defer db.Close()

	records, e := db.SelectById(testCheeseMaker{}, francs.Id)
	require.Equal(t, nil, e)

	act, ok := records.(testCheeseMaker)
	require.Equal(t, true, ok)
	require.Equal(t, francs, act)
}

func Test_Storm_DeleteById_1(t *testing.T) {
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		[]testCheeseMaker{bobs, francs},
	)
	defer db.Close()

	e := db.DeleteById(testCheeseMaker{}, int64(1))
	require.Equal(t, nil, e)

	records := selectAllTestCheeseMakers(t, db)
	require.Equal(t, francs, records[0])
	require.Equal(t, 1, len(records))
}

func Test_Storm_Drop_1(t *testing.T) {
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		nil,
	)
	defer db.Close()

	e := db.Drop(testCheeseMaker{})
	require.Equal(t, nil, e)

	tableNames := selectTestTableNamesFromSqliteSchema(t, db)

	require.Equal(t, "testCheese", tableNames[0])
	require.Equal(t, 1, len(tableNames))
}
