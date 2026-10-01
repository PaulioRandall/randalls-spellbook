package storm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
)

func Test_Storm_Create_1(t *testing.T) {
	// Creates table in database GIVEN valid model AND no
	// table in database yet.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	requireDummyTableExists(t, st)
}

func Test_Storm_Upsert_1(t *testing.T) {
	// Creates table in database and inserts data GIVEN
	// table not yet in database.

	st := openStormDatabase(t)
	defer st.Close()

	data := TestTable{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.Put(data)
	require.NoError(t, e)

	requireDummyTableExists(t, st)
	requireDummyTableRows(t, st, data)
}

func Test_Storm_Upsert_2(t *testing.T) {
	// Inserts data GIVEN table is already created.

	st := openStormDatabase(t)
	defer st.Close()

	data := TestTable{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.Create(TestTable{})
	require.NoError(t, e)

	e = st.Put(data)
	require.NoError(t, e)

	requireDummyTableExists(t, st)
	requireDummyTableRows(t, st, data)
}

func Test_Storm_Upsert_3(t *testing.T) {
	// Updates data GIVEN object is already within table.

	st := openStormDatabase(t)
	defer st.Close()

	original := TestTable{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.Put(original)
	require.NoError(t, e)

	updated := TestTable{
		Id:     123,
		Name:   "xyZ",
		Rating: 987.654,
	}

	e = st.Put(updated)
	require.NoError(t, e)

	requireDummyTableExists(t, st)
	requireDummyTableRows(t, st, updated)
}

func Test_Storm_List_1(t *testing.T) {
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

	e := st.Put(a)
	require.NoError(t, e)

	e = st.Put(b)
	require.NoError(t, e)

	act, e := st.List(TestTable{})
	require.NoError(t, e)

	require.Equal(t, a, act[0])
	require.Equal(t, b, act[1])
	require.Equal(t, 2, len(act))
}

func Test_Storm_List_2(t *testing.T) {
	// Selects nothing from database GIVEN table doesn't
	// exist.

	st := openStormDatabase(t)
	defer st.Close()

	act, e := st.List(TestTable{})
	require.NoError(t, e)

	require.Equal(t, 0, len(act))
}

func Test_Storm_Get_1(t *testing.T) {
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

	e := st.Put(a)
	require.NoError(t, e)

	e = st.Put(b)
	require.NoError(t, e)

	act, e := st.Get(TestTable{}, 1)
	require.NoError(t, e)

	require.Equal(t, a, act)
}

func Test_Storm_Get_2(t *testing.T) {
	// Returns not found error GIVEN data not in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	_, e = st.Get(TestTable{}, 1)
	require.ErrorIs(t, e, ErrObjectNotFound)
}

func Test_Storm_Get_3(t *testing.T) {
	// Returns not found error GIVEN table doesn't exist
	// in database.

	st := openStormDatabase(t)
	defer st.Close()

	_, e := st.Get(TestTable{}, 1)
	require.ErrorIs(t, e, ErrObjectNotFound)
}

func Test_Storm_Delete_1(t *testing.T) {
	// Does nothing and no error GIVEN table doesn't exist.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Delete(TestTable{}, 1)
	require.NoError(t, e)
}

func Test_Storm_Delete_2(t *testing.T) {
	// Does nothing and no error GIVEN row not in database.

	st := openStormDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	e = st.Delete(TestTable{}, 1)
	require.NoError(t, e)
}

func Test_Storm_Delete_3(t *testing.T) {
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

	e = st.Put(a)
	require.NoError(t, e)

	e = st.Put(b)
	require.NoError(t, e)

	requireDummyTableRows(t, st, a, b)

	e = st.Delete(TestTable{}, 1)
	require.NoError(t, e)

	requireDummyTableRows(t, st, b)
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

	_, e = scumble.QuerySqliteSchema(st.db, "TestTable")
	require.ErrorIs(t, e, scumble.ErrEntityNotFound)
}
