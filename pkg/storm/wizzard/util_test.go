package wizzard

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
	TestIdCol = Property{
		GoType:     reflect.TypeOf(int(0)),
		GoIndex:    0,
		GoName:     "Id",
		SqlType:    "INTEGER",
		SqlName:    "Id",
		SqlDefault: int(0),
		IsKey:      true,
	}

	TestNameCol = Property{
		GoType:     reflect.TypeOf(""),
		GoIndex:    1,
		GoName:     "Name",
		SqlType:    "TEXT",
		SqlName:    "Name",
		SqlDefault: "''",
		IsKey:      false,
	}

	TestRatingCol = Property{
		GoType:     reflect.TypeOf(float64(0)),
		GoIndex:    2,
		GoName:     "Rating",
		SqlType:    "REAL",
		SqlName:    "Rating",
		SqlDefault: float64(0),
		IsKey:      false,
	}
)
