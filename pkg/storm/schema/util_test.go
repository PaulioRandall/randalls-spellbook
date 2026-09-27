package schema

import (
	"database/sql"
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/stretchr/testify/require"
)

func createAndPopulateTestDb(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	require.NoError(t, e)

	_, e = db.Exec(`
		CREATE TABLE players (
			id INTEGER,
			name TEXT NOT NULL,
			rating REAL NOT NULL DEFAULT 2.50,
			PRIMARY KEY (id)
		)
	`)
	require.NoError(t, e)

	return db
}
