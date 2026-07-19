package projectdata

import (
	"errors"
	//"fmt"
)

// EntityTypeVideo is not used for video media.
// Joking, of course it is :)
const EntityTypeVideo EntityType = "video"

// video holds core information specific to videos.
type video struct {
	// id is unique and primary key to the entity. It
	// must be unique, never empty, and never change.
	id EntityId

	// Name is the readable and meaningful name for human
	// and AI agent users. It must never be empty.
	name string

	// Description is a detailed explanation of the media for
	// human and AI users.
	//
	// It should compliment the Name field but is also
	// intended for general notes about the media.
	// Description may be empty.
	desc string

	// path is the file path to the video file.
	path SystemFile
}

func (v video) _entity() {}

// EntityId returns the id of the video.
func (v video) EntityId() EntityId {
	return v.id
}

// EntityType always returns EntityTypeVideo, i.e. 'video'.
func (v video) EntityType() EntityType {
	return EntityTypeVideo
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

// addVideo adds a new video to the project data.
//
// An ID will be assigned, ovewriting any existing value.
// The entity's content is assumed to be valid.
func (pd *ProjectData) addVideo(video video) {
	video.id = newEntityId()
	pd.videos = append(pd.videos, video)
}

// getVideo returns an existing video entity given an
// id.
//
// If not found, an empty video is returned.
func (pd *ProjectData) getVideo(id EntityId) video {
	return findEntityById(pd.videos, id, video{})
}

// updateVideo updates an existing video with the
// passed one.
//
// The passed video's EntityId is used to lookup the
// existing video. If not found, an error is returned.
func (pd *ProjectData) updateVideo(video video) error {
	if updateEntity(pd.videos, video) {
		return nil
	}

	return newVideoNotFoundError(video.id)
}

// deleteVideo removes the video with the specified
// id.
//
// If no matching video is found, an error is returned.
func (pd *ProjectData) deleteVideo(id EntityId) error {
	i := findEntityIndexById(pd.videos, id)

	if i < 0 {
		return newVideoNotFoundError(id)
	}

	deleteFromSlice(pd.videos, i)
	return nil
}

func newVideoNotFoundError(id EntityId) error {
	// TODO: add ID to message.
	return errors.New("Video not found")
}
