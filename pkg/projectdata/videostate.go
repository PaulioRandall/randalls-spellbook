package projectdata

// VideoState holds information about the state of a
// single video.
type VideoState struct {
	Entity
	MediaRef Entity

	// TODO: last known time position
}

func (vs VideoState) _entity() {}
