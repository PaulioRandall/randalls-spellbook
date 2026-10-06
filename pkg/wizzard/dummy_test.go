package wizzard

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sqlick"
)

type Dummy struct {
	Id     int
	Name   string
	Rating float64
}

func createTestDb(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	require.NoError(t, e)
	return db
}

func createDummyTestDb(t *testing.T) *sql.DB {
	db := createTestDb(t)

	_, e := db.Exec(`
		CREATE TABLE Dummy (
			Id INTEGER NOT NULL DEFAULT 0,
			Name TEXT NOT NULL DEFAULT '',
			Rating REAL NOT NULL DEFAULT 0,
			PRIMARY KEY (Id)
		)
	`)
	require.NoError(t, e)

	return db
}

func modelOfDummy() Model {
	return Model{
		GoType:  reflect.TypeOf(Dummy{}),
		GoName:  "Dummy",
		SqlName: "Dummy",
		Props: []Property{
			Property{
				GoType:     reflect.TypeOf(int(0)),
				GoIndex:    0,
				GoName:     "Id",
				SqlType:    "INTEGER",
				SqlName:    "Id",
				SqlDefault: int(0),
				IsKey:      true,
			},
			Property{
				GoType:     reflect.TypeOf(""),
				GoIndex:    1,
				GoName:     "Name",
				SqlType:    "TEXT",
				SqlName:    "Name",
				SqlDefault: "''",
				IsKey:      false,
			},
			Property{
				GoType:     reflect.TypeOf(float64(0)),
				GoIndex:    2,
				GoName:     "Rating",
				SqlType:    "REAL",
				SqlName:    "Rating",
				SqlDefault: float64(0),
				IsKey:      false,
			},
		},
	}
}

func requireDummyTableDoesNotExist(t *testing.T, db *sql.DB) {
	_, e := sqlick.QuerySqliteSchema(db, "Dummy")
	require.ErrorIs(t, e, sqlick.ErrEntityNotFound)
}

func requireDummyTableExists(t *testing.T, db *sql.DB) {
	tableSchema, e := sqlick.QuerySqliteSchema(db, "Dummy")
	require.NoError(t, e)
	require.Equal(t, "Dummy", tableSchema.Name)
	require.Equal(t, "Dummy", tableSchema.TableName)

	tableInfo, e := sqlick.QueryTableInfo(db, "Dummy")
	require.NoError(t, e)

	idColInfo := sqlick.TableInfo{
		Name:            "Id",
		Type:            "INTEGER",
		NotNull:         true,
		HasDefaultValue: true,
		DefaultValue:    "0",
		PrimaryKey:      1,
	}
	require.Equal(t, idColInfo, tableInfo[0])

	nameColInfo := sqlick.TableInfo{
		Name:            "Name",
		Type:            "TEXT",
		NotNull:         true,
		HasDefaultValue: true,
		DefaultValue:    "''",
		PrimaryKey:      0,
	}
	require.Equal(t, nameColInfo, tableInfo[1])

	ratingColInfo := sqlick.TableInfo{
		Name:            "Rating",
		Type:            "REAL",
		NotNull:         true,
		HasDefaultValue: true,
		DefaultValue:    "0",
		PrimaryKey:      0,
	}
	require.Equal(t, ratingColInfo, tableInfo[2])
}

func requireDummyTableContains(
	t *testing.T,
	db *sql.DB,
	data ...Dummy,
) []Dummy {
	rows, e := db.Query(`
		SELECT
			Id,
			Name,
			Rating
		FROM
			Dummy
	`)
	require.NoError(t, e)
	defer rows.Close()

	var dbData []Dummy

	for rows.Next() {
		var d Dummy
		e := rows.Scan(&d.Id, &d.Name, &d.Rating)
		require.NoError(t, e)
		dbData = append(dbData, d)
	}

	for i, exp := range data {
		require.Equal(t, exp, dbData[i])
	}

	return dbData
}
