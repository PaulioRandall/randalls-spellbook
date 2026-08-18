package entity

import (
	"errors"
	"path/filepath"
	"strings"
)

// MediaType is a string representing the type of media,
// e.g. video, audio, PDF, etc.
type MediaType string

// ToMediaType returns the MediaType given its string
// representation. If there is no matching MediaType then
// an error is returned.
func ToMediaType(s string) (MediaType, error) {
	mt := MediaType(s)

	if SupportsMediaType(mt) {
		return mt, nil
	}

	// TODO: Fmt with the media type for debugging.
	return "", errors.New("Invalid media type")
}

// SupportsMediaType returns true if the given MediaType is
// supported.
func SupportsMediaType(mt MediaType) bool {
	switch mt {
	case MediaTypeVideo:
		return true
	default:
		return false
	}
}

func (mt MediaType) String() string {
	return string(mt)
}

// MediaTypeVideo is only used by video implementations
// of Media.
const (
	MediaTypeUnknown MediaType = "unknown"
	MediaTypeVideo   MediaType = "video"
)

// Media holds information about a user defined media
// entity.
type Media struct {
	// EntityId of the video.
	//
	// It must be unique within the project, never empty, and
	// never change.
	EntityId EntityId

	// MediaType is type of the media, e.g. video, audio,
	// PDF, etc.
	MediaType MediaType

	// Name is the user defined readable and meaningful name
	// for human users and AI agents. This is not the
	// filename, localPath is the filename.
	//
	// It must never be empty and should should be trimmed
	// of whitespace.
	Name string

	// Description is the user defined detailed explanation
	// of the video for human users and AI agents.
	//
	// It should compliment the Name field but is also
	// intended for general notes. It may be may be empty and
	// should should be trimmed of whitespace.
	Description string

	// LocalPath is the file path to the video file within
	// the local file system.
	LocalPath string
}

// MakeMedia creates and returns a new Media entity
// ensuring passed properties are valid. All values are
// trimmed and localPath is cleaned before being checked
// and assigned to the new media.
//
// If checks pass then the new media is returned else an
// empty Media and an error. An error will occur if the
// name or localPath are empty, or if localPath is not a
// valid absolute filepath.
//
// The exsitance of the file or file type are not checked
// but this may change in the future.
func MakeMedia(
	mediaType MediaType,
	name string,
	description string,
	localPath string,
) (Media, error) {
	empty := Media{}

	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	localPath = strings.TrimSpace(localPath)

	if name == "" {
		return empty, errors.New("name must not be empty")
	}

	if localPath == "" {
		return empty, errors.New("localPath must not be empty")
	}

	localPath = filepath.Clean(localPath)
	if !filepath.IsAbs(localPath) {
		return empty, errors.New("localPath must be absolute")
	}

	m := Media{
		EntityId:    randomEntityId(),
		MediaType:   mediaType,
		Name:        name,
		Description: description,
		LocalPath:   localPath,
	}

	return m, nil
}

// GetEntityId returns the unique entity ID of the media.
func (m Media) GetEntityId() EntityId {
	return m.EntityId
}

// GetMediaType returns the type of the media, e.g. video,
// audio, PDF, etc.
func (m Media) GetMediaType() MediaType {
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

// IsEmpty returns true if the instance is an empty
// Media.
func (m Media) IsEmpty() bool {
	return m == Media{}
}

// FindMediaById looks up a specific media in mediaList
// given an entityId. If not found, an empty Media is
// returned.
func FindMediaById(
	mediaList []Media,
	entityId EntityId,
) Media {
	for _, m := range mediaList {
		if m.EntityId == entityId {
			return m
		}
	}

	return Media{}
}
