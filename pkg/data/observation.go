package data

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
)

// Observation represents a single observation within a
// single media.
type Observation struct {
	// EntityId unique to the entity.
	EntityId string `json:"entityId"`

	// MediaId is the EntityId of the Media.
	MediaId string `json:"mediaId"`

	// StartTime is the number of seconds from the start of
	// the media that the observation starts. The end time
	// is calculated by adding the Duration property.
	StartTime float64 `json:"startTime"`

	// Duration is the number of seconds the observation
	// occurs. Duration may be 0, representing an instant in
	// time rather than a range.
	Duration float64 `json:"duration"`

	// Description is the user's description and notes about
	// the observation.
	Description string `json:"description"`
}

// CleanObservation returns a new Observation entity from
// an exsiting one ensuring all properties are valid. The
// description is trimmed, and if the EntityId is empty
// then a new one is allocated.
//
// If checks pass then the observation is returned else an
// empty Observation and an error. An error will occur if
// the StartTime or Duration are not finite non-negative
// numbers, or MediaId is empty.
func CleanObservation(ob Observation) (Observation, error) {
	empty := Observation{}

	ob.Description = strings.TrimSpace(ob.Description)
	ob.StartTime = roundFloat64(ob.StartTime, 3)
	ob.Duration = roundFloat64(ob.Duration, 3)

	if ob.EntityId == "" {
		ob.EntityId = randomEntityId()
	}

	if ob.MediaId == "" {
		return empty, errors.New("MediaId must be set")
	}

	if math.IsNaN(ob.StartTime) || ob.StartTime < 0 {
		return empty, fmt.Errorf(
			"StartTime '%.2f' must be finite and non-negative",
			ob.StartTime,
		)
	}

	if math.IsNaN(ob.Duration) || ob.Duration < 0 {
		return empty, fmt.Errorf(
			"Duration '%.2f' must be finite and non-negative",
			ob.Duration,
		)
	}

	return ob, nil
}

func createObservationTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS observation (
			id INTEGER PRIMARY KEY,
			entity_id TEXT UNIQUE NOT NULL ON CONFLICT ROLLBACK,
			media_id TEXT NOT NULL,
			start_time REAL NOT NULL,
			duration REAL NOT NULL,
			description TEXT NOT NULL,
			FOREIGN KEY (media_id)
				REFERENCES media(entity_id)
				ON DELETE CASCADE
		)
	`

	_, e := db.Exec(query)

	if e != nil {
		return fmt.Errorf("Failed to create observation table: %w", e)
	}

	return nil
}

func insertObservation(db *sql.DB, ob Observation) error {
	ob, e := CleanObservation(ob)
	if e != nil {
		return e
	}

	query := `
		INSERT INTO observation (
			entity_id,
			media_id,
			start_time,
			duration,
			description
		) VALUES (?, ?, ?, ?, ?)
	`

	_, e = db.Exec(
		query,
		ob.EntityId,
		ob.MediaId,
		ob.StartTime,
		ob.Duration,
		ob.Description,
	)

	if e != nil {
		return fmt.Errorf("Failed to insert observation: %w", e)
	}

	return nil
}

func listObservationsByMediaId(
	db *sql.DB,
	mediaId string,
) ([]Observation, error) {
	query := `
		SELECT
			entity_id,
			media_id,
			start_time,
			duration,
			description
		FROM
			observation
		WHERE
			media_id = ?
	`

	rows, e := db.Query(query, mediaId)
	if e != nil {
		return nil, fmt.Errorf("Failed to query observation by ID: %w", e)
	}

	defer rows.Close()

	obs, e := parseObservationTableRows(rows)
	if e != nil {
		return nil, e
	}

	return obs, nil
}

func parseObservationTableRows(
	rows *sql.Rows,
) ([]Observation, error) {
	var result []Observation

	for rows.Next() {
		var ob Observation

		e := rows.Scan(
			&ob.EntityId,
			&ob.MediaId,
			&ob.StartTime,
			&ob.Duration,
			&ob.Description,
		)

		if e != nil {
			return nil, fmt.Errorf("Failed to scan observation row: %w", e)
		}

		result = append(result, ob)
	}

	return result, rows.Err()
}
