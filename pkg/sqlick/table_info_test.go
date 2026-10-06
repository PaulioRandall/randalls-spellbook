package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_QueryTableInfo_1(t *testing.T) {
	// GIVEN the player table exists
	// WHEN  queried
	// THEN  all expected TableInfo objects are returned and
	//       populated as expected.

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, e := QueryTableInfo(db, "players")
	exp := []TableInfo{
		TableInfo{
			Name:            "id",
			Type:            "INTEGER",
			NotNull:         false,
			HasDefaultValue: false,
			DefaultValue:    "",
			PrimaryKey:      1,
		},
		TableInfo{
			Name:            "name",
			Type:            "TEXT",
			NotNull:         true,
			HasDefaultValue: false,
			DefaultValue:    "",
			PrimaryKey:      0,
		},
		TableInfo{
			Name:            "rating",
			Type:            "REAL",
			NotNull:         true,
			HasDefaultValue: true,
			DefaultValue:    "2.50",
			PrimaryKey:      0,
		},
	}

	require.NoError(t, e)
	require.Equal(t, exp, act)
}

func Test_QueryTableInfo_2(t *testing.T) {
	// GIVEN a table dosen't exist
	// THEN  returns an empty slice.

	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, e := QueryTableInfo(db, "meh")

	require.NoError(t, e)
	require.Equal(t, []TableInfo(nil), act)
}
