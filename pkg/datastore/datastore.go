package datastore

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/entity"
)

// Datastore is an interface for accessing data within the
// project. It may be backed by any form of data storage,
// local or remote.
type Datastore interface {
	// Path returns the URL to the datastore.
	Path() string

	// Open opens the datastore, creating it if it doesn't
	// already exist.
	Open() error

	// IsOpen returns true if the datastore is open.
	IsOpen() bool

	// Close closes the datastore and cleans up all
	// resources.
	Close() error

	// InsertMedia inserts a media. The media is assumed to
	// be valid.
	InsertMedia(entity.Media) error

	// GetAllMedia gets all the media records.
	GetAllMedia() ([]entity.Media, error)

	// GetMedia gets the media with the given EntityID.
	GetMediaById(entityId entity.EntityId) (entity.Media, error)
}
