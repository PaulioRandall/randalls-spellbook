package scumble

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_QuerySqliteSchema_1(t *testing.T) {
	// GIVEN the player table exists
	// WHEN  queried
	// THEN  a SqliteSchema object is returned with known
	//       fields populated as expected.

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, e := QuerySqliteSchema(db, "players")

	require.NoError(t, e)
	require.Equal(t, "table", act.Type)
	require.Equal(t, "players", act.Name)
	require.Equal(t, "players", act.TableName)

	// We don't know exactly what act.Rootpage and act.Sql
	// will contain so we can't fully test them.
}

func Test_QuerySqliteSchema_2(t *testing.T) {
	// GIVEN a table dosen't exist
	// WHEN  queried
	// THEN  returns an ErrEntityNotFound error.

	db := createAndPopulateTestDb(t)
	defer db.Close()

	_, e := QuerySqliteSchema(db, "meh")
	require.ErrorIs(t, e, ErrEntityNotFound)
}
