package projectdata

// VideoObservation holds information about a single
// video observation.
type VideoObservation struct {
	Entity
	MediaRef Entity

	// Description of the observation as provided by the
	// analyst (user).
	Description string
}

func (vo VideoObservation) _entity() {}
