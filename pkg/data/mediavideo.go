package data

// MediaTypeVideo used for video media files.
const MediaTypeVideo MediaType = "video"

// MediaSectionVideo represents a section video media,
// usually defined explicitly by a user.
type MediaSectionVideo struct {
	// EntityID and Primary Key unique to the video media
	// section. It must never be empty and never change.
	EntityID

	// MediaRef links to the media. It must never be
	// an empty structure and must never change.
	MediaRef

	// Name is the section's meaningful label for human
	// readers, assuming it has been thoughtfully chosen.
	//
	// Name does not need to be unique and can be empty.
	Name string

	// Description of the section, as provided by the user,
	// to aid users in understanding and using it.
	//
	// Description should compliment the Name field but is
	// also intended for general notes about the media
	// section. Description may be empty.
	Description string

	// TODO: start time
	// TODO: end time
}

// MediaStateVideo represents the state of a specific video
// media.
type MediaStateVideo struct {
	// EntityID and Primary Key unique to the media state. It
	// must never be empty and never change.
	EntityID

	// MediaRef links to the media. It must never be
	// an empty structure and must never change.
	MediaRef

	// TODO: last known time position
}

type MediaObservationVideo struct {
	// EntityID and Primary Key unique to the video media
	// observation. It must never be empty and never change.
	EntityID

	// MediaId
	MediaId EntityID
}
