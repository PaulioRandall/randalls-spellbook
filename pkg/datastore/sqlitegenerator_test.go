package datastore

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Test_sqlitegenerator_1(t *testing.T) {
	given := DbTable{
		Name: "TestEntity",
		Columns: []DbColumn{
			DbColumn{
				Name:     "EntityId",
				DataType: "string",
			},
			DbColumn{
				Name:     "Name",
				DataType: "string",
			},
			DbColumn{
				Name:     "Age",
				DataType: "int",
			},
			DbColumn{
				Name:     "Height",
				DataType: "float64",
			},
		},
	}

	mapToSqliteType := func(dataType string) (string, bool) {
		switch dataType {
		case "string":
			return "TEXT", true
		case "int":
			return "INTEGER", true
		case "float64":
			return "REAL", true
		default:
			return "", false
		}
	}

	exp := joinLines(
		"CREATE TABLE IF NOT EXISTS TestEntity (",
		"  EntityId TEXT NOT NULL,",
		"  Name TEXT NOT NULL,",
		"  Age INTEGER NOT NULL,",
		"  Height REAL NOT NULL",
		")",
	)

	act, e := generateCreateTableSql(given, mapToSqliteType)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}
