package storm

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/scumble"
)

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
