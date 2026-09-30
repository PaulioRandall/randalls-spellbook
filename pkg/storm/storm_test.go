package storm

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

func openStormDatabase(t *testing.T) *Storm {
	st := New(":memory:")

	e := st.Open()
	require.NoError(t, e)

	return st
}

func requireDummyTableExists(t *testing.T, st *Storm) {
	tableSchema, e := scumble.QuerySqliteSchema(st.db, "TestTable")
	require.NoError(t, e)
	require.Equal(t, "TestTable", tableSchema.Name)
	require.Equal(t, "TestTable", tableSchema.TableName)
}

func requireDummyTableRows(t *testing.T, st *Storm, data ...TestTable) {
	act := queryDummyTable(t, st)
	require.Equal(t, data, act)
}

func queryDummyTable(t *testing.T, st *Storm) []TestTable {
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

func Test_Storm_Open_Close_1(t *testing.T) {
	st := openStormDatabase(t)
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
