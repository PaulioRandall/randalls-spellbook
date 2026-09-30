package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
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
