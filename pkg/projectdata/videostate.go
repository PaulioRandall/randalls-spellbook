package projectdata

// videoState holds information about the state of a
// single video.
type videoState struct {
	// id is unique and primary key to the entity. It
	// must be unique, never empty, and never change.
	id EntityId

	// TODO
	MediaRef Entity

	// TODO: last known time position
}

func (vs videoState) _entity() {}
