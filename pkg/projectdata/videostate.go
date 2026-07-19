package projectdata

import (
	"errors"
	//"fmt"
)

// videoState holds information about the state of a
// single video.
type videoState struct {
	// id is unique and primary key to the entity. It
	// must be unique, never empty, and never change.
	id EntityId

	// mediaId is the EntityId of the media. This field must
	// never be empty or change.
	mediaId EntityId

	// TODO: last known time position
}

// EntityId returns the id of the video.
func (vs videoState) EntityId() EntityId {
	return vs.id
}

// MediaId returns the EntityId of the associated video.
func (vs videoState) MediaId() EntityId {
	return vs.mediaId
}

// AddVideoState adds a new video state to the project data.
//
// An ID will be assigned, ovewriting any existing value.
// The entity's content is assumed to be valid.
func (pd *ProjectData) AddVideoState(videoState videoState) {
	videoState.id = newEntityId()
	pd.videoStates = append(pd.videoStates, videoState)
}

// GetVideoState returns an existing video state entity given
// an id.
//
// If not found, an empty videoState is returned.
func (pd *ProjectData) GetVideoState(id EntityId) videoState {
	return findEntityById(pd.videoStates, id, videoState{})
}

// DeleteVideoState removes the video state with the
// specified id.
//
// If no matching video state is found, an error is returned.
func (pd *ProjectData) DeleteVideoState(id EntityId) error {
	i := findEntityIndexById(pd.videoStates, id)

	if i < 0 {
		return newVideoStateNotFoundError(id)
	}

	pd.videoStates = deleteFromSlice(pd.videoStates, i)
	return nil
}

func newVideoStateNotFoundError(id EntityId) error {
	// TODO: add ID to message.
	return errors.New("Video state not found")
}

var _ = _mediaSupportEntity(videoState{})
