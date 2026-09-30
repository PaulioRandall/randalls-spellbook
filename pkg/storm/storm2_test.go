package storm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
)

type Dummy struct {
	Id      int
	Name    string
	Rating  float64
	ignored *int
}

func openStormDatabase(t *testing.T) *Storm {
	st := New(":memory:")

	e := st.Open()
	require.NoError(t, e)

	return st
}

func requireDummyTableExists(t *testing.T, st *Storm) {
	tableSchema, e := schema.QuerySqliteSchema(st.db, "Dummy")
	require.NoError(t, e)
	require.Equal(t, "Dummy", tableSchema.Name)
	require.Equal(t, "Dummy", tableSchema.TableName)
}

func requireDummyTableRows(t *testing.T, st *Storm, data ...Dummy) {
	act := queryDummyTable(t, st)
	require.Equal(t, data, act)
}

func queryDummyTable(t *testing.T, st *Storm) []Dummy {
	rows, e := st.db.Query(`
		SELECT
			Id,
			Name,
			Rating
		FROM
			Dummy
	`)
	require.NoError(t, e)

	var result []Dummy

	for rows.Next() {
		var p Dummy
		e := rows.Scan(&p.Id, &p.Name, &p.Rating)
		require.NoError(t, e)
		result = append(result, p)
	}

	return result
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
	// Creates table in database GIVEN valid model AND no
	// table in database yet.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(Dummy{})
	require.NoError(t, e)

	requireDummyTableExists(t, st)
}

func Test_Storm_Insert_1(t *testing.T) {
	// Creates table in database and inserts data GIVEN
	// valid model AND table not yet in database.

	st := openStormDatabase(t)
	defer st.Close()

	data := Dummy{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.Insert(data)
	require.NoError(t, e)

	requireDummyTableExists(t, st)
	requireDummyTableRows(t, st, data)
}

func Test_Storm_Insert_2(t *testing.T) {
	// Inserts data into database GIVEN table already in
	// database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(Dummy{})
	require.NoError(t, e)

	data := Dummy{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e = st.Insert(data)
	require.NoError(t, e)

	requireDummyTableRows(t, st, data)
}

func Test_Storm_Drop_1(t *testing.T) {
	// Does nothing GIVEN model table not in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Drop(Dummy{})
	require.NoError(t, e)
}

func Test_Storm_Drop_2(t *testing.T) {
	// Drops table GIVEN model table in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(Dummy{})
	require.NoError(t, e)

	e = st.Drop(Dummy{})
	require.NoError(t, e)

	_, e = schema.QuerySqliteSchema(st.db, "Dummy")
	require.ErrorIs(t, e, schema.ErrEntityNotFound)
}

func Test_Storm_Update_1(t *testing.T) {
	// Updates data in database GIVEN table exists AND
	// item exists in table.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(Dummy{})
	require.NoError(t, e)

	original := Dummy{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e = st.Insert(original)
	require.NoError(t, e)

	updated := Dummy{
		Id:     123,
		Name:   "Xyz",
		Rating: 987.654,
	}

	e = st.Update(updated)
	require.NoError(t, e)

	requireDummyTableRows(t, st, updated)
}
