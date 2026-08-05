package entity

// Media is a generic interface covering all forms of media
// stored on the local filesystem, e.g. video files, audio
// files, PDF files, etc. Some media types may have
// additional functions specific to them.
type Media interface {
	// EntityId returns the unique entity ID of the media.
	EntityId() EntityId

	// MediaType returns the type of the media, e.g. video,
	// audio, PDF, etc.
	MediaType() MediaType

	// Name returns the user defined readable and meaningful
	// name for humans users and AI agents. This is not the
	// filename, LocalPath returns the filename. It should
	// never be empty and should be trimmed of whitespace.
	Name() string

	// Description is the user defined detailed explanation
	// of the video for human users and AI agents It
	// compliments the media name but may also hold general
	// notes. It may be empty and should be trimmed of
	// whitespace.
	Description() string

	// LocalPath is the path to the media file within the
	// local file system. There is no guarantee that the
	// file exists, as it may have been deleted or moved
	// since it was added to project. This may also happen
	// if the project is moved to a new computer without
	// copying the media files to matching locations in
	// the new file system.
	LocalPath() string

	// IsEmpty returns true if the media instance is empty.
	IsEmpty() bool
}

// FindMediaById looks up a media in the list of media for
// a given entityId. If not found, nil is returned.
func FindMediaById(media []Media, entityId EntityId) Media {
	for _, m := range media {
		if m.EntityId() == entityId {
			return m
		}
	}

	return nil
}
