package scumble

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_QueryTable_1(t *testing.T) {
	db := createAndPopulateTestDb(t)
	defer db.Close()

	act, e := QueryTable(db, "Player")
	require.NoError(t, e)

	exp := SqlTable{
		Name: "Player",
		Columns: []SqlColumn{
			SqlColumn{
				Name:       "Id",
				Type:       "INTEGER",
				Default:    int(0),
				PrimaryKey: true,
			},
			SqlColumn{
				Name:       "Name",
				Type:       "TEXT",
				Default:    "''",
				PrimaryKey: false,
			},
			SqlColumn{
				Name:       "Rating",
				Type:       "REAL",
				Default:    float64(0),
				PrimaryKey: false,
			},
		},
	}

	require.Equal(t, exp, act)
}
