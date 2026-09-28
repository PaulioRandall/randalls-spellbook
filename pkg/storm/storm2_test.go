package storm

import (
	//"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
)

func openStormDatabase(t *testing.T) *Storm {
	st := New(":memory:")

	e := st.Open()
	require.NoError(t, e)

	return st
}

func queryDatabase(
	t *testing.T,
	st *Storm,
	colCount int,
	sql string,
	args ...any,
) [][]string {
	rows, e := st.db.Query(sql, args...)
	require.NoError(t, e)

	var result [][]any

	for rows.Next() {
		ptrs := make([]any, colCount)
		for i := range ptrs {
			s := ""
			ptrs[i] = &s
		}

		e := rows.Scan(ptrs...)
		require.NoError(t, e)
		result = append(result, ptrs)
	}

	strResult := make([][]string, len(result))

	for i, col := range result {
		strCol := make([]string, len(col))

		for j := range col {
			strCol[j] = *(col[j].(*string))
		}

		strResult[i] = strCol
	}

	return strResult
}

func Test_Storm_Open_Close_1(t *testing.T) {
	st := openStormDatabase(t)
	defer st.Close()

	require.Equal(t, true, st.IsOpen())

	e := st.Close()
	require.NoError(t, e)
	require.Equal(t, false, st.IsOpen())
}

func Test_Storm_Create_1(t *testing.T) {
	// GIVEN a valid table.
	// THEN  it should be created.

	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int
	}

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(Player{})
	require.NoError(t, e)

	tableSchema, e := schema.QuerySqliteSchema(st.db, "Player")
	require.NoError(t, e)
	require.Equal(t, "Player", tableSchema.Name)
	require.Equal(t, "Player", tableSchema.TableName)
}

func Test_Storm_Insert_1(t *testing.T) {
	// GIVEN table hasn't been created yet.
	// THEN  table should be created and data inserted.

	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int
	}

	st := openStormDatabase(t)
	defer st.Close()

	data := Player{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.Insert(data)
	require.NoError(t, e)

	// Check table was created.
	tableSchema, e := schema.QuerySqliteSchema(st.db, "Player")
	require.NoError(t, e)
	require.Equal(t, "Player", tableSchema.Name)
	require.Equal(t, "Player", tableSchema.TableName)

	// Check data was entered.
	rows, e := st.db.Query(`
		SELECT
			Id,
			Name,
			Rating
		FROM
			Player
	`)
	require.NoError(t, e)

	var act []Player
	for rows.Next() {
		var p Player
		e := rows.Scan(&p.Id, &p.Name, &p.Rating)
		require.NoError(t, e)
		act = append(act, p)
	}

	require.Equal(t, []Player{data}, act)
}

func Test_Storm_Insert_2(t *testing.T) {
	// GIVEN table was already created.
	// THEN  data should be inserted.

	type Player struct {
		Id      int
		Name    string
		Rating  float64
		ignored *int
	}

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(Player{})
	require.NoError(t, e)

	data := Player{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e = st.Insert(data)
	require.NoError(t, e)

	// Check data was entered.
	rows, e := st.db.Query(`
		SELECT
			Id,
			Name,
			Rating
		FROM
			Player
	`)
	require.NoError(t, e)

	var act []Player
	for rows.Next() {
		var p Player
		e := rows.Scan(&p.Id, &p.Name, &p.Rating)
		require.NoError(t, e)
		act = append(act, p)
	}

	require.Equal(t, []Player{data}, act)
}
