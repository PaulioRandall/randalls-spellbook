package data2

import (
	"errors"
	"path/filepath"
	"strings"
)

// MediaTypeVideo is only used by video implementations
// of Media.
const (
	MediaTypeVideo string = "video"
)

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
func (m Media) Clean() (Media, error) {
	m.EntityId = strings.TrimSpace(m.EntityId)
	m.Name = strings.TrimSpace(m.Name)
	m.Description = strings.TrimSpace(m.Description)
	m.LocalPath = strings.TrimSpace(m.LocalPath)

	if m.EntityId == "" {
		m.EntityId = randomEntityId()
	}

	if m.Name == "" {
		return Media{}, errors.New("Name must not be empty")
	}

	if m.LocalPath == "" {
		return Media{}, errors.New("LocalPath must not be empty")
	}

	m.LocalPath = filepath.Clean(m.LocalPath)
	if !filepath.IsAbs(m.LocalPath) {
		return Media{}, errors.New("LocalPath must be absolute")
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
