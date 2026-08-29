package data

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// MediaTypeVideo is only used by video implementations
// of Media.
const (
	MediaTypeVideo string = "video"
)

var supportedMediaTypes []string = []string{
	MediaTypeVideo,
}

// SupportsMediaType returns true if the given media type
// is supported.
func SupportsMediaType(mt string) bool {
	for _, smt := range supportedMediaTypes {
		if mt == smt {
			return true
		}
	}
	return false
}

// Media holds information about a user defined media
// entity. It maps directly to a single table in the
// database.
type Media struct {
	// EntityId of the video.
	//
	// It must be unique within the project, never empty, and
	// never change.
	EntityId string `json:"entityId"`

	// MediaType is type of the media, e.g. video, audio,
	// PDF, etc.
	MediaType string `json:"mediaType"`

	// Name is the user defined readable and meaningful name
	// for human users and AI agents. This is not the
	// filename, localPath is the filename.
	//
	// It must never be empty and should should be trimmed
	// of whitespace.
	Name string `json:"name"`

	// Description is the user defined detailed explanation
	// of the video for human users and AI agents.
	//
	// It should compliment the Name field but is also
	// intended for general notes. It may be may be empty and
	// should should be trimmed of whitespace.
	Description string `json:"description"`

	// LocalPath is the file path to the video file within
	// the local file system.
	LocalPath string `json:"localPath"`
}

// CleanMedia returns a new Media entity from an exsiting
// one ensuring all properties are valid. All values are
// trimmed and localPath is cleaned before being checked
// and assigned to the new media. If the EntityId is empty
// then a new one is allocated.
//
// If checks pass then the media is returned else an empty
// Media and an error. An error will occur if the Name or
// LocalPath are empty, or if LocalPath is not a valid
// absolute filepath.
//
// The existence of the file or file type are not checked
// but this may change in the future.
func CleanMedia(m Media) (Media, error) {
	empty := Media{}

	m.EntityId = strings.TrimSpace(m.EntityId)
	m.Name = strings.TrimSpace(m.Name)
	m.Description = strings.TrimSpace(m.Description)
	m.LocalPath = strings.TrimSpace(m.LocalPath)

	if m.EntityId == "" {
		m.EntityId = randomEntityId()
	}

	if m.Name == "" {
		return empty, errors.New("Name must not be empty")
	}

	if m.LocalPath == "" {
		return empty, errors.New("LocalPath must not be empty")
	}

	m.LocalPath = filepath.Clean(m.LocalPath)
	if !filepath.IsAbs(m.LocalPath) {
		return empty, errors.New("LocalPath must be absolute")
	}

	return m, nil
}

// GetEntityId returns the unique entity ID of the media.
func (m Media) GetEntityId() string {
	return m.EntityId
}

// GetMediaType returns the type of the media, e.g. video,
// audio, PDF, etc.
func (m Media) GetMediaType() string {
	return m.MediaType
}

// GetName returns the user defined readable and meaningful
// name for humans users and AI agents. This is not the
// filename, LocalPath returns the filename. It should
// never be empty and should be trimmed of whitespace.
func (m Media) GetName() string {
	return m.Name
}

// GetDescription is the user defined detailed explanation
// of the video for human users and AI agents It
// compliments the media name but may also hold general
// notes. It may be empty and should be trimmed of
// whitespace.
func (m Media) GetDescription() string {
	return m.Description
}

// GetLocalPath is the path to the media file within the
// local file system. There is no guarantee that the
// file exists, as it may have been deleted or moved
// since it was added to project. This may also happen
// if the project is moved to a new computer without
// copying the media files to matching locations in
// the new file system.
func (m Media) GetLocalPath() string {
	return m.LocalPath
}

func createMediaTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS media (
			id INTEGER PRIMARY KEY,
			entity_id TEXT UNIQUE NOT NULL ON CONFLICT ROLLBACK,
			media_type TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			local_path TEXT NOT NULL
		)
	`

	_, e := db.Exec(query)

	if e != nil {
		return fmt.Errorf("Failed to create media table: %w", e)
	}

	return nil
}

func insertMedia(db *sql.DB, media Media) (Media, error) {
	empty := Media{}

	media, e := CleanMedia(media)
	if e != nil {
		return empty, e
	}

	query := `
		INSERT INTO media (
			entity_id,
			media_type,
			name,
			description,
			local_path
		) VALUES (?, ?, ?, ?, ?)
	`

	_, e = db.Exec(
		query,
		media.EntityId,
		media.MediaType,
		media.Name,
		media.Description,
		media.LocalPath,
	)

	if e != nil {
		return empty, fmt.Errorf("Failed to insert media: %w", e)
	}

	return media, nil
}

func listMedia(db *sql.DB) ([]Media, error) {
	query := `
		SELECT
			entity_id,
			media_type,
			name,
			description,
			local_path
		FROM
			media
	`

	rows, e := db.Query(query)
	if e != nil {
		return nil, fmt.Errorf("Failed to query media table (all): %w", e)
	}

	defer rows.Close()
	return parseMediaRows(rows)
}

func getMediaByEntityId(
	db *sql.DB,
	entityId string,
) (Media, error) {
	query := `
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

	empty := Media{}

	rows, e := db.Query(query, entityId)
	if e != nil {
		return empty, fmt.Errorf("Failed to query media table by EntityId: %w", e)
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

func parseMediaRows(rows *sql.Rows) ([]Media, error) {
	var result []Media

	for rows.Next() {
		var m Media

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
