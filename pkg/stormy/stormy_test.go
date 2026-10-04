package stormy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
)

type TestTable struct {
	Id      int
	Name    string
	ignored *int
	Rating  float64
}

func openStormyDatabase(t *testing.T) *Stormy {
	st := New(":memory:")
	e := st.Open()
	require.NoError(t, e)
	return st
}

func requireTableExists(t *testing.T, st *Stormy, name string) {
	tableSchema, e := scumble.QuerySqliteSchema(st.db, name)
	require.NoError(t, e)
	require.Equal(t, name, tableSchema.Name)
	require.Equal(t, name, tableSchema.TableName)
}

func requireDummyTableExists(t *testing.T, st *Stormy) {
	tableSchema, e := scumble.QuerySqliteSchema(st.db, "TestTable")
	require.NoError(t, e)
	require.Equal(t, "TestTable", tableSchema.Name)
	require.Equal(t, "TestTable", tableSchema.TableName)
}

func requireTableContains(
	t *testing.T,
	st *Stormy,
	table string,
	data ...TestTable,
) []TestTable {
	rows, e := st.db.Query(`
		SELECT
			Id,
			Name,
			Rating
		FROM
	` + table)
	require.NoError(t, e)
	defer rows.Close()

	var dbData []TestTable

	for rows.Next() {
		var tt TestTable
		e := rows.Scan(&tt.Id, &tt.Name, &tt.Rating)
		require.NoError(t, e)
		dbData = append(dbData, tt)
	}

	for i, exp := range data {
		require.Equal(t, exp, dbData[i])
	}

	return dbData
}

func requireDummyTableRows(t *testing.T, st *Stormy, data ...TestTable) {
	act := queryDummyTable(t, st)
	require.Equal(t, data, act)
}

func queryDummyTable(t *testing.T, st *Stormy) []TestTable {
	rows, e := st.db.Query(`
		SELECT
			Id,
			Name,
			Rating
		FROM
			TestTable
	`)
	require.NoError(t, e)

	var result []TestTable

	for rows.Next() {
		var p TestTable
		e := rows.Scan(&p.Id, &p.Name, &p.Rating)
		require.NoError(t, e)
		result = append(result, p)
	}

	return result
}

func Test_Stormy_Open_Close_1(t *testing.T) {
	st := openStormyDatabase(t)
	defer st.Close()

	require.Equal(t, true, st.IsOpen())

	e := st.Close()
	require.NoError(t, e)
	require.Equal(t, false, st.IsOpen())
}

func Test_isInMemoryDatabase(t *testing.T) {
	cases := map[string]bool{
		// The simple case.
		":memory:": true,

		// The simple case with file scheme.
		"file::memory:": true,

		// Special name is case-sensitive (AFAIK).
		":MEMORY:": false,

		// Scheme is case-sensitive (AFAIK).
		"FILE::memory:": false,

		// File scheme should not automatically be considered
		// in-memory.
		"file:meh": false,

		// The simple case with query params
		":memory:?abc=123": true,

		// The simple case with file scheme and query params
		"file::memory:?abc=123": true,

		// Other than ':memory:', if it doesn't use the file
		// scheme then it's not an in-memory database (AFAIK).
		"./:memory:": false,

		// Treat URLs without scheme as not in-memory; probably
		// fails parsing.
		"example.com": false,

		// Named in-memory.
		"file:my-db?mode=memory": true,

		// In-memory specified using VFS parameter.
		"file:my-db?vfs=memdb": true,

		// Encoded path.
		"file:%3Amemory%3A": true,
	}

	for path, is := range cases {
		require.Equal(
			t,
			is,
			isInMemoryDatabase(path),
			`Path: %s`,
			path,
		)
	}
}
