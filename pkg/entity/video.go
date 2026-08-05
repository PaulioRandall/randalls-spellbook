package entity

import (
	"errors"
	"path/filepath"
	"strings"
)

// MediaTypeVideo is only used by video implementations
// of Media.
const MediaTypeVideo MediaType = "video"

// Video holds information about a user defined video
// entity.
type Video struct {
	// entityId of the video.
	//
	// It must be unique within the project, never empty, and
	// never change.
	entityId EntityId

	// name is the user defined readable and meaningful name
	// for human users and AI agents. This is not the
	// filename, localPath is the filename.
	//
	// It must never be empty and should should be trimmed
	// of whitespace.
	name string

	// description is the user defined detailed explanation
	// of the video for human users and AI agents.
	//
	// It should compliment the Name field but is also
	// intended for general notes. It may be may be empty and
	// should should be trimmed of whitespace.
	description string

	// localPath is the file path to the video file within
	// the local file system.
	localPath string
}

// MakeVideo creates and returns a new video entity
// ensuring passed properties are valid. All values are
// trimmed and localPath cleaned before being checked and
// assigned to the new video.
//
// If checks pass then the new Video is returned else an
// empty Video and an error. An error will occur if the
// name or localPath are empty, or if localPath is not a
// valid absolute filepath.
//
// The exsitance of the file or file type are not checked
// but this may change in the future.
func MakeVideo(
	name, description, localPath string,
) (Video, error) {
	emptyVideo := Video{}

	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	localPath = strings.TrimSpace(localPath)

	if name == "" {
		return emptyVideo, errors.New("name must not be empty")
	}

	if localPath == "" {
		return emptyVideo, errors.New("localPath must not be empty")
	}

	localPath = filepath.Clean(localPath)
	if !filepath.IsAbs(localPath) {
		return emptyVideo, errors.New("localPath must be absolute")
	}

	video := Video{
		entityId:    randomEntityId(),
		name:        name,
		description: description,
		localPath:   localPath,
	}

	return video, nil
}

// EntityId returns the unique entity ID of the media.
func (v Video) EntityId() EntityId {
	return v.entityId
}

// MediaType always returns MediaTypeVideo.
func (v Video) MediaType() MediaType {
	return MediaTypeVideo
}

// Name returns the user defined readable and meaningful
// name for humans users and AI agents. This is not the
// filename, LocalPath returns the filename. It should
// never be empty and should be trimmed of whitespace.
func (v Video) Name() string {
	return v.name
}

// Description is the user defined detailed explanation
// of the video for human users and AI agents It
// compliments the video name but may also hold general
// notes. It may be empty and should be trimmed of
// whitespace.
func (v Video) Description() string {
	return v.description
}

// LocalPath is the path to the video file within the local
// file system. There is no guarantee that the file exists,
// as it may have been deleted or moved since it was
// added to application. This may also happen if the
// project is moved to a new computer without copying the
// media files to the same location in the new file system.
func (v Video) LocalPath() string {
	return v.localPath
}

// IsEmpty returns true if the instance is an empty Video.
func (v Video) IsEmpty() bool {
	return v == Video{}
}

// Compile time type check.
var _ Media = Media(Video{})
