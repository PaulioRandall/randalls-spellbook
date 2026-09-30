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
	Id      int64
	MakerId int64
	Name    string
	Rating  float64
	wasInit bool
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

var cheddar = testCheese{
	Id:      3,
	MakerId: bobs.Id,
	Name:    "Cheddar",
	Rating:  7.7,
	wasInit: false,
}

func openCreateInsert(
	t *testing.T,
	models []any,
	objects ...any,
) *Storm {
	db := New(":memory:")

	e := db.Open()
	require.NoError(t, e)

	if len(models) > 0 {
		e = db.Create(models...)
		require.NoError(t, e)
	}

	if len(objects) > 0 {
		e = db.Insert(objects...)
		require.NoError(t, e)
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
	require.NoError(t, e)

	var result []testCheeseMaker

	for rows.Next() {
		tcm := testCheeseMaker{}
		e = rows.Scan(&tcm.Id, &tcm.Name, &tcm.Country)
		require.NoError(t, e)
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
	require.NoError(t, e)

	var result []string

	for rows.Next() {
		var name string
		e = rows.Scan(&name)
		require.NoError(t, e)
		result = append(result, name)
	}

	return result
}

func Test_Storm_DeleteById_1(t *testing.T) {
	db := openCreateInsert(
		t,
		[]any{testCheeseMaker{}, testCheese{}},
		bobs, francs,
	)
	defer db.Close()

	e := db.DeleteById(testCheeseMaker{}, int64(1))
	require.NoError(t, e)

	records := selectAllTestCheeseMakers(t, db)
	require.Equal(t, francs, records[0])
	require.Equal(t, 1, len(records))
}
