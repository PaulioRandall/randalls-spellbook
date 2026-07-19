package projectdata

import (
	"errors"
	//"fmt"
)

// videoSection holds information about a single
// video section.
//
// videoSection allows users to partition or define
// sections of their video so it's easier to analyse.
type videoSection struct {
	// id is unique and primary key to the entity. It
	// must be unique, never empty, and never change.
	id EntityId

	// mediaId is the EntityId of the media. This field must
	// never be empty or change.
	mediaId EntityId

	// name is the section's meaningful label for human
	// readers, assuming it has been thoughtfully chosen.
	//
	// Name does not need to be unique and can be empty.
	name string

	// desc of the section, as provided by the user,
	// to aid users in understanding and using it.
	//
	// The description should compliment the Name field but
	// is also intended for general notes about the media
	// section. Description may be empty.
	desc string

	// TODO: start time
	// TODO: end time
}

// EntityId returns the id of the video.
func (vs videoSection) EntityId() EntityId {
	return vs.id
}

// MediaId returns the EntityId of the associated video.
func (vs videoSection) MediaId() EntityId {
	return vs.mediaId
}

// name returns the user defined name for the section.
func (vs videoSection) Name() string {
	return vs.name
}

// Description returns the user defined section description.
func (vs videoSection) Description() string {
	return vs.desc
}

// AddVideoSection adds a new video section to the project
// data.
//
// An ID will be assigned, ovewriting any existing value.
// The entity's content is assumed to be valid.
func (pd *ProjectData) AddVideoSection(videoSection videoSection) {
	videoSection.id = newEntityId()
	pd.videoSections = append(pd.videoSections, videoSection)
}

// GetVideoSection returns an existing video section entity
// given an id.
//
// If not found, an empty videoSection is returned.
func (pd *ProjectData) GetVideoSection(id EntityId) videoSection {
	return findEntityById(pd.videoSections, id, videoSection{})
}

// DeleteVideoSection removes the video section with the
// specified id.
//
// If no matching video section is found, an error is
// returned.
func (pd *ProjectData) DeleteVideoSection(id EntityId) error {
	i := findEntityIndexById(pd.videoSections, id)

	if i < 0 {
		return newVideoSectionNotFoundError(id)
	}

	pd.videoSections = deleteFromSlice(pd.videoSections, i)
	return nil
}

func newVideoSectionNotFoundError(id EntityId) error {
	// TODO: add ID to message.
	return errors.New("Video section not found")
}

var _ = _mediaSupportEntity(videoSection{})
