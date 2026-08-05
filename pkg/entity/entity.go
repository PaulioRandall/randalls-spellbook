package entity

import (
	"github.com/google/uuid"
)

// EntityId is the core ID type for every stored entity.
type EntityId string

// String returns the string representation of the
// EntityId.
func (id EntityId) String() string {
	return string(id)
}

// MediaType is a string representing the type of media,
// e.g. video, audio, PDF, etc.
//
// Media type constants are defined within the Go file
// holding its implementation.
type MediaType string

// randomEntityId randomly generates a new EntityId in the
// form of a UUID.
func randomEntityId() EntityId {
	return EntityId(uuid.New().String())
}
