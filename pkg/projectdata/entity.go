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

func newEntityId() EntityId {
	return EntityId(uuid.New().String())
}

type _entity interface {
	EntityId() EntityId
}

type _mediaEntity interface {
	_entity
	Name() string
	Description() string
}

type _mediaSupportEntity interface {
	MediaId() EntityId
}

type _observationEntity interface {
	_entity
	_mediaSupportEntity
	Description() string
}
