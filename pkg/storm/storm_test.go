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
