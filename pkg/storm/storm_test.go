package storm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/storm/schema"
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
	tableSchema, e := schema.QuerySqliteSchema(st.db, "TestTable")
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

func Test_Storm_Create_1(t *testing.T) {
	// Creates table in database GIVEN valid model AND no
	// table in database yet.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	requireDummyTableExists(t, st)
}

func Test_Storm_Insert_1(t *testing.T) {
	// Creates table in database and inserts data GIVEN
	// valid model AND table not yet in database.

	st := openStormDatabase(t)
	defer st.Close()

	data := TestTable{
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

	e := st.Create(TestTable{})
	require.NoError(t, e)

	data := TestTable{
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

	e := st.Drop(TestTable{})
	require.NoError(t, e)
}

func Test_Storm_Drop_2(t *testing.T) {
	// Drops table GIVEN model table in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	e = st.Drop(TestTable{})
	require.NoError(t, e)

	_, e = schema.QuerySqliteSchema(st.db, "TestTable")
	require.ErrorIs(t, e, schema.ErrEntityNotFound)
}

func Test_Storm_Update_1(t *testing.T) {
	// Updates data in database GIVEN table exists AND
	// item exists in table.

	st := openStormDatabase(t)
	defer st.Close()

	original := TestTable{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.Insert(original)
	require.NoError(t, e)

	updated := TestTable{
		Id:     123,
		Name:   "Xyz",
		Rating: 987.654,
	}

	e = st.Update(updated)
	require.NoError(t, e)

	requireDummyTableRows(t, st, updated)
}

func Test_Storm_Select_1(t *testing.T) {
	// Selects data from database GIVEN data exists.

	st := openStormDatabase(t)
	defer st.Close()

	a := TestTable{
		Id:     1,
		Name:   "A",
		Rating: 1.1,
	}

	b := TestTable{
		Id:     2,
		Name:   "B",
		Rating: 2.2,
	}

	e := st.Insert(a)
	require.NoError(t, e)

	e = st.Insert(b)
	require.NoError(t, e)

	act, e := st.Select(TestTable{})
	require.NoError(t, e)

	require.Equal(t, a, act[0])
	require.Equal(t, b, act[1])
	require.Equal(t, 2, len(act))
}

func Test_Storm_Select_2(t *testing.T) {
	// Selects nothing from database GIVEN table doesn't
	// exist.

	st := openStormDatabase(t)
	defer st.Close()

	act, e := st.Select(TestTable{})
	require.NoError(t, e)

	require.Equal(t, 0, len(act))
}

func Test_Storm_SelectById_1(t *testing.T) {
	// Returns requested data from database GIVEN data
	// exists.

	st := openStormDatabase(t)
	defer st.Close()

	a := TestTable{
		Id:     1,
		Name:   "A",
		Rating: 1.1,
	}

	b := TestTable{
		Id:     2,
		Name:   "B",
		Rating: 2.2,
	}

	e := st.Insert(a)
	require.NoError(t, e)

	e = st.Insert(b)
	require.NoError(t, e)

	act, e := st.SelectById(TestTable{}, 1)
	require.NoError(t, e)

	require.Equal(t, a, act)
}

func Test_Storm_SelectById_2(t *testing.T) {
	// Returns not found error GIVEN data not in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	_, e = st.SelectById(TestTable{}, 1)
	require.ErrorIs(t, e, ErrObjectNotFound)
}

func Test_Storm_SelectById_3(t *testing.T) {
	// Returns not found error GIVEN table doesn't exist
	// in database.

	st := openStormDatabase(t)
	defer st.Close()

	_, e := st.SelectById(TestTable{}, 1)
	require.ErrorIs(t, e, ErrObjectNotFound)
}

func Test_Storm_DeleteById_1(t *testing.T) {
	// Does nothing and no error GIVEN table doesn't exist.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.DeleteById(TestTable{}, 1)
	require.NoError(t, e)
}

func Test_Storm_DeleteById_2(t *testing.T) {
	// Does nothing and no error GIVEN row not in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	e = st.DeleteById(TestTable{}, 1)
	require.NoError(t, e)
}

func Test_Storm_DeleteById_3(t *testing.T) {
	// Deletes the  GIVEN row not in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	a := TestTable{
		Id:     1,
		Name:   "A",
		Rating: 1.1,
	}

	b := TestTable{
		Id:     2,
		Name:   "B",
		Rating: 2.2,
	}

	e = st.Insert(a)
	require.NoError(t, e)

	e = st.Insert(b)
	require.NoError(t, e)

	requireDummyTableRows(t, st, a, b)

	e = st.DeleteById(TestTable{}, 1)
	require.NoError(t, e)

	requireDummyTableRows(t, st, b)
}
