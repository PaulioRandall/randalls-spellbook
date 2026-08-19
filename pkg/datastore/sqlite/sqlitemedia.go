package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/PaulioRandall/randalls-spellbook/pkg/entity"
)

var queryCreateMediaTable string = `
	CREATE TABLE media (
		id INTEGER PRIMARY KEY,
		entity_id TEXT UNIQUE NOT NULL ON CONFLICT ROLLBACK,
		media_type TEXT NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		local_path TEXT NOT NULL
	)
`

var queryInsertIntoMediaTable string = `
	INSERT INTO media (
		entity_id,
		media_type,
		name,
		description,
		local_path
	) VALUES (?, ?, ?, ?, ?)
`

var queryAllMedia string = `
	SELECT
		entity_id,
		media_type,
		name,
		description,
		local_path
	FROM
		media
`

var queryMediaById string = `
	SELECT
		entity_id,
		media_type,
		name,
		description,
		local_path
	FROM
		media
	WHERE
		entity_id = ?
`

// InsertMedia satisfies the Datastore interface.
func (ds *SQLiteDatabase) InsertMedia(
	media entity.Media,
) error {
	_, e := ds.db.Exec(
		queryInsertIntoMediaTable,
		media.EntityId,
		media.MediaType,
		media.Name,
		media.Description,
		media.LocalPath,
	)

	if e != nil {
		return fmt.Errorf("Failed to insert media: %w", e)
	}

	return nil
}

// GetAllMedia satisfies the Datastore interface.
func (ds *SQLiteDatabase) GetAllMedia() ([]entity.Media, error) {
	rows, e := ds.db.Query(queryAllMedia)
	if e != nil {
		return nil, fmt.Errorf("Failed to query media table: %w", e)
	}

	defer rows.Close()
	return parseMediaRows(rows)
}

// GetMediaById satisfies the Datastore interface.
func (ds *SQLiteDatabase) GetMediaById(
	entityId entity.EntityId,
) (entity.Media, error) {
	empty := entity.Media{}

	rows, e := ds.db.Query(queryMediaById, entityId)
	if e != nil {
		return empty, fmt.Errorf("Failed to query media table: %w", e)
	}

	defer rows.Close()

	media, e := parseMediaRows(rows)
	if e != nil {
		return empty, e
	}

	if len(media) == 0 {
		return empty, nil
	}

	return media[0], nil
}

func parseMediaRows(rows *sql.Rows) ([]entity.Media, error) {
	var result []entity.Media

	for rows.Next() {
		var m entity.Media

		e := rows.Scan(
			&m.EntityId,
			&m.MediaType,
			&m.Name,
			&m.Description,
			&m.LocalPath,
		)

		if e != nil {
			return nil, fmt.Errorf("Failed to scan media row: %w", e)
		}

		result = append(result, m)
	}

	return result, rows.Err()
}
