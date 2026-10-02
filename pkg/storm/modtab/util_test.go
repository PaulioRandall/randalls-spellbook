package modtab

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"
)

func createTestDb(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	require.NoError(t, e)
	return db
}

func createAndPopulateTestDb(t *testing.T) *sql.DB {
	db := createTestDb(t)

	_, e := db.Exec(`
		CREATE TABLE Player (
			Id INTEGER,
			Name TEXT,
			Rating REAL,
			PRIMARY KEY (Id)
		)
	`)
	require.NoError(t, e)

	return db
}

var (
	TestIdCol = ModelColumn{
		GoType:     reflect.TypeOf(int(0)),
		GoIndex:    0,
		GoName:     "Id",
		SqlType:    "INTEGER",
		SqlName:    "Id",
		SqlDefault: int(0),
		PrimaryKey: true,
	}

	TestNameCol = ModelColumn{
		GoType:     reflect.TypeOf(""),
		GoIndex:    1,
		GoName:     "Name",
		SqlType:    "TEXT",
		SqlName:    "Name",
		SqlDefault: "''",
		PrimaryKey: false,
	}

	TestRatingCol = ModelColumn{
		GoType:     reflect.TypeOf(float64(0)),
		GoIndex:    2,
		GoName:     "Rating",
		SqlType:    "REAL",
		SqlName:    "Rating",
		SqlDefault: float64(0),
		PrimaryKey: false,
	}
)
