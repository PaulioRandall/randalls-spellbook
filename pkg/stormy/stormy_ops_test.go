package stormy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sqlick"
)

func Test_Stormy_Create_1(t *testing.T) {
	// When creating a table
	// if the table doesn't exist
	// then the table is created.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	requireTableExists(t, st, "TestTable")
}

func Test_Stormy_Create_2(t *testing.T) {
	// When creating a table
	// if the table already exists
	// then nothing happens.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	e = st.Create(TestTable{})
	require.NoError(t, e)
}

func Test_Stormy_CreateAs_1(t *testing.T) {
	// When creating a table
	// if the table doesn't exist
	// then the table is created.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.CreateAs("TestDummy", TestTable{})
	require.NoError(t, e)

	requireTableExists(t, st, "TestDummy")
}

func Test_Stormy_CreateAs_2(t *testing.T) {
	// When creating a table
	// if the table already exists
	// then nothing happens.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.CreateAs("TestDummy", TestTable{})
	require.NoError(t, e)

	e = st.CreateAs("TestDummy", TestTable{})
	require.NoError(t, e)
}

func Test_Stormy_Put_1(t *testing.T) {
	// When inserting data
	// if the table doesn't exist yet
	// then the table is created
	// and then the data is inserted.

	st := openStormyDatabase(t)
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

func Test_Stormy_Put_2(t *testing.T) {
	// When inserting data
	// if the table already exists
	// then the data is inserted.

	st := openStormyDatabase(t)
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

func Test_Stormy_Put_3(t *testing.T) {
	// When updating data
	// if the table and row exists
	// then the data is updated without error.

	st := openStormyDatabase(t)
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

func Test_Stormy_PutAs_1(t *testing.T) {
	// When inserting data
	// if the table doesn't exist yet
	// then the table is created
	// and then the data is inserted.

	st := openStormyDatabase(t)
	defer st.Close()

	data := TestTable{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.PutAs("TestDummy", data)
	require.NoError(t, e)

	requireTableExists(t, st, "TestDummy")
	requireTableContains(t, st, "TestDummy", data)
}

func Test_Stormy_PutAs_2(t *testing.T) {
	// When inserting data
	// if the table already exists
	// then the data is inserted.

	st := openStormyDatabase(t)
	defer st.Close()

	data := TestTable{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.CreateAs("TestDummy", TestTable{})
	require.NoError(t, e)

	e = st.PutAs("TestDummy", data)
	require.NoError(t, e)

	requireTableExists(t, st, "TestDummy")
	requireTableContains(t, st, "TestDummy", data)
}

func Test_Stormy_PutAs_3(t *testing.T) {
	// When updating data
	// if the table and row exists
	// then the data is updated without error.

	st := openStormyDatabase(t)
	defer st.Close()

	original := TestTable{
		Id:     123,
		Name:   "Abc",
		Rating: 123.456,
	}

	e := st.PutAs("TestDummy", original)
	require.NoError(t, e)

	updated := TestTable{
		Id:     123,
		Name:   "xyZ",
		Rating: 987.654,
	}

	e = st.PutAs("TestDummy", updated)
	require.NoError(t, e)

	requireTableExists(t, st, "TestDummy")
	requireTableContains(t, st, "TestDummy", updated)
}

func Test_Stormy_List_1(t *testing.T) {
	// When selecting all data in a table
	// if the table exists
	// then all data is returned.

	st := openStormyDatabase(t)
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

	act, e := st.List(TestTable{}, "")
	require.NoError(t, e)

	require.Equal(t, a, act[0])
	require.Equal(t, b, act[1])
	require.Equal(t, 2, len(act))
}

func Test_Stormy_List_2(t *testing.T) {
	// When selecting all data in a table
	// if the table doesn't exists
	// then an empty result set (nil) is returned
	// and no error occurs.

	st := openStormyDatabase(t)
	defer st.Close()

	act, e := st.List(TestTable{}, "")
	require.NoError(t, e)

	require.Equal(t, 0, len(act))
}

func Test_Stormy_ListAs_1(t *testing.T) {
	// When selecting all data in a table
	// if the table exists
	// then all data is returned.

	st := openStormyDatabase(t)
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

	e := st.PutAs("TestDummy", a)
	require.NoError(t, e)

	e = st.PutAs("TestDummy", b)
	require.NoError(t, e)

	act, e := st.ListAs("TestDummy", TestTable{}, "")
	require.NoError(t, e)

	require.Equal(t, a, act[0])
	require.Equal(t, b, act[1])
	require.Equal(t, 2, len(act))
}

func Test_Stormy_ListAs_2(t *testing.T) {
	// When selecting all data in a table
	// if the table doesn't exists
	// then an empty result set (nil) is returned
	// and no error occurs.

	st := openStormyDatabase(t)
	defer st.Close()

	act, e := st.ListAs("TestDummy", TestTable{}, "")
	require.NoError(t, e)

	require.Equal(t, 0, len(act))
}

func Test_Stormy_Get_1(t *testing.T) {
	// When selecting a specific object/row
	// if the table and row exist
	// then the data is returned.

	st := openStormyDatabase(t)
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

	act, found, e := st.Get(TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
	require.Equal(t, true, found)
	require.Equal(t, a, act)
}

func Test_Stormy_Get_2(t *testing.T) {
	// When selecting a specific object/row
	// if the table exists
	// but row does not exist
	// then not found error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	_, found, e := st.Get(TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
	require.Equal(t, false, found)
}

func Test_Stormy_Get_3(t *testing.T) {
	// When selecting a specific object/row
	// if the table does not exists
	// then not found error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	_, found, e := st.Get(TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
	require.Equal(t, false, found)
}

func Test_Stormy_GetAs_1(t *testing.T) {
	// When selecting a specific object/row
	// if the table and row exist
	// then the data is returned.

	st := openStormyDatabase(t)
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

	e := st.PutAs("TestDummy", a)
	require.NoError(t, e)

	e = st.PutAs("TestDummy", b)
	require.NoError(t, e)

	act, found, e := st.GetAs("TestDummy", TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
	require.Equal(t, true, found)
	require.Equal(t, a, act)
}

func Test_Stormy_GetAs_2(t *testing.T) {
	// When selecting a specific object/row
	// if the table exists
	// but row does not exist
	// then not found error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.CreateAs("TestDummy", TestTable{})
	require.NoError(t, e)

	_, found, e := st.GetAs("TestDummy", TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
	require.Equal(t, false, found)
}

func Test_Stormy_GetAs_3(t *testing.T) {
	// When selecting a specific object/row
	// if the table does not exists
	// then not found error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	_, found, e := st.GetAs("TestDummy", TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
	require.Equal(t, false, found)
}

func Test_Stormy_Delete_1(t *testing.T) {
	// When deleting a specific object/row
	// if the table does not exists
	// then nothing happens
	// and no error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.Delete(TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
}

func Test_Stormy_Delete_2(t *testing.T) {
	// When deleting a specific object/row
	// if the table exists
	// but row does not exist
	// then nothing happens
	// and no error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	e = st.Delete(TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
}

func Test_Stormy_Delete_3(t *testing.T) {
	// When deleting a specific object/row
	// if the table and row exist
	// then the row is deleted.

	st := openStormyDatabase(t)
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

	e = st.Delete(TestTable{}, "Id = ?", 1)
	require.NoError(t, e)

	requireDummyTableRows(t, st, b)
}

func Test_Stormy_DeleteAs_1(t *testing.T) {
	// When deleting a specific object/row
	// if the table does not exists
	// then nothing happens
	// and no error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.DeleteAs("TestDummy", TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
}

func Test_Stormy_DeleteAs_2(t *testing.T) {
	// When deleting a specific object/row
	// if the table exists
	// but row does not exist
	// then nothing happens
	// and no error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.CreateAs("TestDummy", TestTable{})
	require.NoError(t, e)

	e = st.DeleteAs("TestDummy", TestTable{}, "Id = ?", 1)
	require.NoError(t, e)
}

func Test_Stormy_DeleteAs_3(t *testing.T) {
	// When deleting a specific object/row
	// if the table and row exist
	// then the row is deleted.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.CreateAs("TestDummy", TestTable{})
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

	e = st.PutAs("TestDummy", a)
	require.NoError(t, e)

	e = st.PutAs("TestDummy", b)
	require.NoError(t, e)

	requireTableContains(t, st, "TestDummy", a, b)

	e = st.DeleteAs("TestDummy", TestTable{}, "Id = ?", 1)
	require.NoError(t, e)

	requireTableContains(t, st, "TestDummy", b)
}

func Test_Stormy_Drop_1(t *testing.T) {
	// When removing a table
	// if the table does not exist
	// then nothing happens
	// and no error is returned.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.Drop(TestTable{})
	require.NoError(t, e)
}

func Test_Stormy_Drop_2(t *testing.T) {
	// When removing a table
	// if the table exists
	// then the table is removed.

	st := openStormyDatabase(t)
	defer st.Close()

	e := st.Create(TestTable{})
	require.NoError(t, e)

	e = st.Drop(TestTable{})
	require.NoError(t, e)

	_, e = sqlick.QuerySqliteSchema(st.db, "TestTable")
	require.ErrorIs(t, e, sqlick.ErrEntityNotFound)
}
