package data

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

// SQLiteDatabase satisfies the Store interface using
// SQLite3 as the underlying database.
type sqliteDatabase struct {
	// path is the path to the SLQite file.
	path string

	// db is the database connection. It must be set to nil
	// when the database is closed. It is used to determine
	// if the database is open.
	db *sql.DB
}

// NewStore creates a SQLite database. This function does
// not create or open the database.
func NewStore(path string) Store {
	return &sqliteDatabase{
		path: path,
	}
}

// Path satisfies the Store interface.
func (sqlite *sqliteDatabase) Path() string {
	return sqlite.path
}

// Open satisfies the Store interface. It opens a local
// SQLite3 database file, creating it if it doesn't exist.
func (sqlite *sqliteDatabase) Open() error {
	db, e := sql.Open("sqlite", sqlite.path)

	if e != nil {
		return fmt.Errorf("Unable to open SQLite database: %w", e)
	}

	sqlite.db = db
	return createTables(sqlite.db)
}

// IsOpen satisfies the Store interface.
func (sqlite *sqliteDatabase) IsOpen() bool {
	return sqlite.db != nil
}

// Close satisfies the Store interface.
func (sqlite *sqliteDatabase) Close() error {
	defer func() {
		sqlite.db = nil
	}()
	return sqlite.db.Close()
}

// GetProject satisfies the Store interface.
func (sqlite *sqliteDatabase) GetProject() (Project, error) {
	return getProject(sqlite.db)
}

// InsertMedia satisfies the Store interface.
func (sqlite *sqliteDatabase) InsertMedia(
	media Media,
) (Media, error) {
	return insertMedia(sqlite.db, media)
}

// ListMedia satisfies the Store interface.
func (sqlite *sqliteDatabase) ListMedia() ([]Media, error) {
	return listMedia(sqlite.db)
}

// GetMediaById satisfies the Store interface.
func (sqlite *sqliteDatabase) GetMediaById(
	entityId string,
) (Media, error) {
	return getMediaByEntityId(sqlite.db, entityId)
}

// DeleteMediaById satisfies the Store interface.
func (sqlite *sqliteDatabase) DeleteMediaById(
	entityId string,
) error {
	return deleteMediaByEntityId(sqlite.db, entityId)
}

// InsertObservation satisfies the Store interface.
func (sqlite *sqliteDatabase) InsertObservation(
	obs Observation,
) error {
	return insertObservation(sqlite.db, obs)
}

// ListObservationsByMediaId satisfies the Store interface.
func (sqlite *sqliteDatabase) ListObservationsByMediaId(
	mediaId string,
) ([]Observation, error) {
	return listObservationsByMediaId(sqlite.db, mediaId)
}

// createTables creates all database tables.
func createTables(db *sql.DB) error {
	type tableCreator func(db *sql.DB) error

	tableCreators := []tableCreator{
		createProjectTable,
		createMediaTable,
		createObservationTable,
	}

	for _, createTable := range tableCreators {
		e := createTable(db)
		if e != nil {
			return e
		}
	}

	return nil
}
