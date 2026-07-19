package projectdata

// videoObservation holds information about a single
// video observation.
type videoObservation struct {
	Entity
	MediaRef Entity

	// Description of the observation as provided by the
	// analyst (user).
	Description string
}

func (vo videoObservation) _entity() {}
