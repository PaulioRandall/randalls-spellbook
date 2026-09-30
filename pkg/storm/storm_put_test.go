package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Storm_Put_1(t *testing.T) {
	// Creates table in database and inserts data GIVEN
	// valid model AND table not yet in database.

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

func Test_Storm_Put_2(t *testing.T) {
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

	e = st.Put(data)
	require.NoError(t, e)

	requireDummyTableRows(t, st, data)
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

	e := st.Put(original)
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
