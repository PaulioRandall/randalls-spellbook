package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sqliteGenerator_TableCreate_1(t *testing.T) {
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

func Test_sqliteGenerator_TableCreate_2(t *testing.T) {
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

func Test_sqliteGenerator_TableInsert_1(t *testing.T) {
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
		") VALUES (",
		"  ?,",
		"  ?,",
		"  ?",
		") VALUES (",
		"  ?,",
		"  ?,",
		"  ?",
		")",
	)

	act, e := sqlGen.TableInsert(table, 3)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableInsert_2(t *testing.T) {
	// Error when too few rows specified.
	sqlGen := NewSqliteGenerator()
	_, e := sqlGen.TableInsert(Table{}, 0)
	require.ErrorIs(t, e, ErrTooFewRows)
}

func Test_sqliteGenerator_TableSelectAll_1(t *testing.T) {
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

	act, e := sqlGen.TableSelectAll(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableSelectById_1(t *testing.T) {
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

func Test_sqliteGenerator_TableUpdateById_1(t *testing.T) {
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

func Test_sqliteGenerator_TableDeleteAll_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	table := Table{
		GoName: "TestEntity",
	}

	exp := joinLines(
		"DELETE FROM",
		"  TestEntity",
	)

	act, e := sqlGen.TableDeleteAll(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableDeleteById_1(t *testing.T) {
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
