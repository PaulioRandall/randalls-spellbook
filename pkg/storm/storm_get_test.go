package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
