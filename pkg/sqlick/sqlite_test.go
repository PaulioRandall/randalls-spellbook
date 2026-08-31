package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SqliteGenerator_TableCreate_1(t *testing.T) {
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

	act, e := sqlGen.TableCreate(given)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_TableCreate_2(t *testing.T) {
	// Error when there's no SQL type mapping for a column's
	// Go type.
	sqlGen := NewSqliteGenerator()

	given := Table{
		GoName: "TestEntity",
		Columns: []Column{
			Column{
				GoName: "UserId",
				GoType: "[]int",
			},
		},
	}

	_, e := sqlGen.TableCreate(given)
	require.ErrorIs(t, e, ErrNoSqlType)
}

func Test_SqliteGenerator_TableInsert_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	table := Table{
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
		"INSERT INTO TestEntity (",
		"  UserId,",
		"  Name,",
		"  Height",
		") VALUES (",
		"  ?,",
		"  ?,",
		"  ?",
		")",
	)

	act, e := sqlGen.TableInsert(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_TableSelect_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	table := Table{
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
		"SELECT",
		"  UserId,",
		"  Name,",
		"  Height",
		"FROM",
		"  TestEntity",
	)

	act, e := sqlGen.TableSelect(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_TableSelectById_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	table := Table{
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
		"SELECT",
		"  UserId,",
		"  Name,",
		"  Height",
		"FROM",
		"  TestEntity",
		"WHERE",
		"  UserId = ?",
	)

	act, e := sqlGen.TableSelectById(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_TableUpdateById_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	table := Table{
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
		"UPDATE",
		"  TestEntity",
		"SET",
		"  Name = ?,",
		"  Height = ?",
		"WHERE",
		"  UserId = ?",
	)

	act, e := sqlGen.TableUpdateById(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_TableDeleteById_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	table := Table{
		GoName: "TestEntity",
		Columns: []Column{
			Column{
				GoName: "UserId",
				GoType: "int",
			},
		},
	}

	exp := joinLines(
		"DELETE FROM",
		"  TestEntity",
		"WHERE",
		"  UserId = ?",
	)

	act, e := sqlGen.TableDeleteById(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}
