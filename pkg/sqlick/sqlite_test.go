package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SqliteGenerator_CreateTable_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	given := Table{
		GoName: "TestEntity",
		Columns: []Column{
			Column{
				GoName: "UserId",
				GoType: "int",
			},
			Column{
				GoName: "Name",
				GoType: "string",
			},
			Column{
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

	act, e := sqlGen.CreateTable(given)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}
