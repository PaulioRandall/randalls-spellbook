package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GenerateSql_1(t *testing.T) {
	// Happy path.

	given := SqlickTable{
		GoName: "TestEntity",
		Columns: []SqlickColumn{
			SqlickColumn{
				GoName: "UserId",
				GoType: "int",
			},
			SqlickColumn{
				GoName: "Name",
				GoType: "string",
			},
			SqlickColumn{
				GoName: "Height",
				GoType: "float64",
			},
		},
	}

	exp := joinLines(
		"CREATE TABLE IF NOT EXISTS TestEntity (",
		"  UserId INTEGER NOT NULL,",
		"  Name TEXT NOT NULL,",
		"  Height REAL NOT NULL,",
		"  PRIMARY KEY (UserId)",
		")",
	)

	act, e := GenerateSql(given)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}
