package database

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
		local_path TEXT NOT NULL,
	);
`

var queryMediaById string = `
	SELECT
		entity_id,
		media_type,
		name,
		description,
		local_path,
	FROM media
	WHERE entity_id = '?'
`

// GetMedia satisfies the Datastore interface.
func (ds *SQLiteDatabase) GetMedia(
	entityId entity.EntityId,
) (entity.Media, error) {
	empty := Media{}

	rows, e := db.Query(queryMediaById, entityId)
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
