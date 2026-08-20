package data

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

// randomEntityId randomly generates a new EntityId in the
// form of a UUID.
func randomEntityId() EntityId {
	return EntityId(uuid.New().String())
}
