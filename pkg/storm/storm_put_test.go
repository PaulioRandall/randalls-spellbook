package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
