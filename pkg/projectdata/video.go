package projectdata

import (
	"errors"
	//"fmt"
)

// EntityTypeVideo is not used for video media.
// Joking, of course it is :)
const EntityTypeVideo EntityType = "video"

// ProtoVideo is for constructing new video media.
type ProtoVideo struct {
	// Name is the readable and meaningful name for human
	// and AI agent users. It must never be empty.
	Name string

	// Description is a detailed explanation of the media for
	// human and AI users.
	//
	// It should compliment the Name field but is also
	// intended for general notes about the media.
	// Description may be empty.
	Description string

	// Path is the file path to the video file.
	Path string
}

// video holds core information specific to videos.
type video struct {
	// id is unique and primary key to the entity. It
	// must be unique, never empty, and never change.
	id EntityId

	// name is the readable and meaningful name for human
	// and AI agent users. It must never be empty.
	name string

	// desc is a detailed explanation of the media for
	// human and AI users.
	//
	// It should compliment the Name field but is also
	// intended for general notes about the media.
	// desc may be empty.
	desc string

	// path is the file path to the video file.
	path string
}

// EntityId returns the id of the video.
func (v video) EntityId() EntityId {
	return v.id
}

// Name returns the user defined video name, not the
// filename.
func (v video) Name() string {
	return v.name
}

// Description returns the user defined video description.
func (v video) Description() string {
	return v.desc
}

// Path returns the URL to the video file.
func (v video) Path() string {
	return v.path
}

// AddVideo adds a new video (proto) to the project data.
//
// An ID will be assigned, ovewriting any existing value.
// The entity's content is assumed to be valid.
func (pd *ProjectData) AddVideo(proto ProtoVideo) {
	v := video{
		id:   newEntityId(),
		name: proto.Name,
		desc: proto.Description,
		path: proto.Path,
	}
	pd.videos = append(pd.videos, v)
}

// GetVideo returns an existing video entity given an
// id.
//
// If not found, an empty video is returned.
func (pd *ProjectData) GetVideo(id EntityId) video {
	return findEntityById(pd.videos, id, video{})
}

// DeleteVideo removes the video with the specified
// id.
//
// If no matching video is found, an error is returned.
func (pd *ProjectData) DeleteVideo(id EntityId) error {
	i := findEntityIndexById(pd.videos, id)

	if i < 0 {
		return newVideoNotFoundError(id)
	}

	pd.videos = deleteFromSlice(pd.videos, i)
	return nil
}

func newVideoNotFoundError(id EntityId) error {
	// TODO: add ID to message.
	return errors.New("Video not found")
}

var _ = _mediaEntity(video{})
