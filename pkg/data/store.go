package data

// Store is an interface for accessing data within the
// project. It may be implemented using any form of data
// storage, local or remote.
type Store interface {
	// Path returns the path or URL to the data store as a
	// string.
	Path() string

	// Open opens the store, creating it if it doesn't
	// already exist.
	Open() error

	// IsOpen returns true if the store is open.
	IsOpen() bool

	// Close closes the store and cleans up all resources.
	Close() error

	// GetProject returns the project data.
	GetProject() (Project, error)

	// InsertMedia inserts a media. The media is assumed to
	// be valid.
	InsertMedia(Media) error

	// ListMedia returns all the media entities.
	ListMedia() ([]Media, error)

	// GetMediaById returns the media entity with the given
	// entity ID.
	GetMediaById(string) (Media, error)
}
