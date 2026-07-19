package projectdata

import (
	"errors"
	//"fmt"
)

// videoObservation holds information about a single
// user observation for specific video.
type videoObservation struct {
	// id is unique and primary key to the entity. It
	// must be unique, never empty, and never change.
	id EntityId

	// mediaId is the EntityId of the media. This field must
	// never be empty or change.
	mediaId EntityId

	// desc of the observation as provided by the analyst
	// (user).
	desc string
}

// EntityId returns the id of the video.
func (vo videoObservation) EntityId() EntityId {
	return vo.id
}

// MediaId returns the EntityId of the associated video.
func (vo videoObservation) MediaId() EntityId {
	return vo.mediaId
}

// Description returns the user defined video observation
// description.
func (vo videoObservation) Description() string {
	return vo.desc
}

// AddVideoObservation adds a new video observation to the
// project data.
//
// An ID will be assigned, ovewriting any existing value.
// The entity's content is assumed to be valid.
func (pd *ProjectData) AddVideoObservation(videoObservation videoObservation) {
	videoObservation.id = newEntityId()
	pd.videoObservations = append(pd.videoObservations, videoObservation)
}

// GetVideoObservation returns an existing video observation
// entity given an id.
//
// If not found, an empty videoObservation is returned.
func (pd *ProjectData) GetVideoObservation(id EntityId) videoObservation {
	return findEntityById(pd.videoObservations, id, videoObservation{})
}

// DeleteVideoObservation removes the video observation with
// the specified id.
//
// If no matching video observation is found, an error is
// returned.
func (pd *ProjectData) DeleteVideoObservation(id EntityId) error {
	i := findEntityIndexById(pd.videoObservations, id)

	if i < 0 {
		return newVideoObservationNotFoundError(id)
	}

	pd.videoObservations = deleteFromSlice(pd.videoObservations, i)
	return nil
}

func newVideoObservationNotFoundError(id EntityId) error {
	// TODO: add ID to message.
	return errors.New("Video observation not found")
}

var _ = _mediaSupportEntity(videoObservation{})
