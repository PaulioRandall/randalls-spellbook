package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
