package projectdata

import (
	"errors"
	//"fmt"
)

// EntityTypeVideo is not used for video media.
// Joking, of course it is :)
const EntityTypeVideo EntityType = "video"

// Video holds core information specific to videos.
type Video struct {
	Entity

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

	// FilePath is the path to the video file.
	FilePath SystemFile
}

func (v Video) _entity() {}

// addVideo adds a new video to the project data.
//
// An ID will be assigned, ovewriting any existing value.
// The entity's content is assumed to be valid.
func (pd *ProjectData) addVideo(video Video) {
	video.EntityId = newEntityId()
	pd.videos = append(pd.videos, video)
}

// getVideo returns an existing video entity given an
// id.
//
// If not found, an empty Video is returned.
func (pd *ProjectData) getVideo(id EntityId) Video {
	for _, v := range pd.videos {
		if v.EntityId == id {
			return v
		}
	}

	return Video{}
}

// updateVideo updates an existing video with the
// passed one.
//
// The passed video's EntityId is used to lookup the
// existing video. If not found, an error is returned.
func (pd *ProjectData) updateVideo(video Video) error {
	index := -1

	for i, v := range pd.videos {
		if v.EntityId == video.EntityId {
			index = i
		}
	}

	if index < 0 {
		// TODO: add ID to message.
		return errors.New("Video not found")
	}

	pd.videos[index] = video
	return nil
}

// deleteVideo removes the video with the specified
// id.
//
// If no matching video is found, an error is returned.
func (pd *ProjectData) deleteVideo(id EntityId) error {
	// AIDO

	return nil
}
