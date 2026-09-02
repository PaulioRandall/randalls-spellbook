package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testColUserId = Column{
	GoName: "UserId",
	GoType: int64Type,
}

var testColName = Column{
	GoName: "Name",
	GoType: strType,
}

var testColHeight = Column{
	GoName: "Height",
	GoType: float64Type,
}

var testTableEntity = Table{
	GoName: "Entity",
	Columns: []Column{
		testColUserId,
		testColName,
		testColHeight,
	},
}

func Test_sqliteGenerator_TableCreate_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	exp := joinLines(
		"CREATE TABLE IF NOT EXISTS Entity (",
		"  UserId INTEGER NOT NULL,",
		"  Name TEXT NOT NULL,",
		"  Height REAL NOT NULL,",
		"  PRIMARY KEY (UserId)",
		")",
	)

	act, e := sqlGen.TableCreate(testTableEntity)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableCreate_2(t *testing.T) {
	// Error when there's no SQL type mapping for a column's
	// Go type.
	sqlGen := NewSqliteGenerator()

	given := Table{
		GoName: "Entity",
		Columns: []Column{
			Column{
				GoName: "UserId",
				GoType: int64ArrType,
			},
		},
	}

	_, e := sqlGen.TableCreate(given)
	require.ErrorIs(t, e, ErrNoSqlType)
}

func Test_sqliteGenerator_TableInsert_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	exp := joinLines(
		"INSERT INTO Entity (",
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

	act, e := sqlGen.TableInsert(testTableEntity, 3)
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

	exp := joinLines(
		"SELECT",
		"  UserId,",
		"  Name,",
		"  Height",
		"FROM",
		"  Entity",
	)

	act, e := sqlGen.TableSelectAll(testTableEntity)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableSelectById_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	exp := joinLines(
		"SELECT",
		"  UserId,",
		"  Name,",
		"  Height",
		"FROM",
		"  Entity",
		"WHERE",
		"  UserId = ?",
	)

	act, e := sqlGen.TableSelectById(testTableEntity)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableUpdateById_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	exp := joinLines(
		"UPDATE",
		"  Entity",
		"SET",
		"  Name = ?,",
		"  Height = ?",
		"WHERE",
		"  UserId = ?",
	)

	act, e := sqlGen.TableUpdateById(testTableEntity)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableDeleteAll_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	exp := joinLines(
		"DELETE FROM",
		"  Entity",
	)

	act, e := sqlGen.TableDeleteAll(testTableEntity)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}

func Test_sqliteGenerator_TableDeleteById_1(t *testing.T) {
	// Happy path.
	sqlGen := NewSqliteGenerator()

	exp := joinLines(
		"DELETE FROM",
		"  Entity",
		"WHERE",
		"  UserId = ?",
	)

	act, e := sqlGen.TableDeleteById(testTableEntity)
	require.Equal(t, nil, e)
	require.Equal(t, exp, act)
}
