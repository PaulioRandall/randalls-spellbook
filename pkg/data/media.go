package data

// MediaType refers to the format of a media, e.g. video,
// audio, PDF document, etc.
type MediaType string

// Media holds information about a specific media item,
// e.g. video, audio, PDF document, etc.
//
// It excludes all additional information added by the user
// or through user interaction. Other structures such as
// MediaSection and MediaState are used to hold this
// information.
type Media struct {
	// EntityID and Primary Key unique to the media. It must
	// never be empty and never change.
	EntityID

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

	// Type is the media type, e.g. video, audio, PDF
	// document, etc.
	//
	// The media type is important because it determines how
	// user observations are stored and presented. Type must
	// never be empty and never change.
	Type MediaType

	// FilePath is the path to the media file.
	FilePath SystemFile
}

// MediaRef references a media and intended for embedding
// in structures designed for specific media types.
type MediaRef struct {
	// MediaId references the entity ID of the associated
	// media.
	MediaId EntityID

	// Type is the media type, e.g. video, audio, PDF
	// document, etc.
	//
	// Media type determines a lot about how related media
	// information is stored, this duplicate field avoids
	// needing to look up the media when using entities that
	// reference a specific media. Because a media's type
	// must never change, it's safe to duplicate.
	MediaType MediaType
}

// For each media:
// - Media sections allow users to partition and organise
//   their media so it's easier to analyse.
// - Media states allow users to quickly return to where
//   they left off in the prior session, e.g. the last time
//   position in a video.
