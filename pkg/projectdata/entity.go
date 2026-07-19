package projectdata

import (
	"github.com/google/uuid"
)

// EntityId is the core ID type that every entity, within
// the project, will have a unique value for.
type EntityId string

// EntityType refers to the type of an entity, e.g. media,
// video, tag, etc.
type EntityType string

// Entity is the base structure that all entities embed.
type Entity struct {
	// id and primary key to the entity. It must be unique,
	// never empty, and never change.
	entityId EntityId
}

type _entity interface{ _entity() }

func newEntityId() EntityId {
	return EntityId(uuid.New().String())
}
