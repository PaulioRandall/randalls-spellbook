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

func Test_SqliteGenerator_CreateTable_2(t *testing.T) {
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

	_, e := sqlGen.CreateTable(given)
	require.ErrorIs(t, e, ErrNoSqlType)
}

func Test_SqliteGenerator_Insert_1(t *testing.T) {
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

	act, e := sqlGen.Insert(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_Select_1(t *testing.T) {
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

	act, e := sqlGen.Select(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_SelectById_1(t *testing.T) {
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

	act, e := sqlGen.SelectById(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_Update_1(t *testing.T) {
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

	act, e := sqlGen.Update(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_SqliteGenerator_Delete_1(t *testing.T) {
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

	act, e := sqlGen.Delete(table)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}
