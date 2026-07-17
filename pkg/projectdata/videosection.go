package projectdata

// VideoSection holds information about a single
// video section.
//
// VideoSection allows users to partition or define
// sections of their video so it's easier to analyse.
type VideoSection struct {
	Entity
	MediaRef Entity

	// Name is the section's meaningful label for human
	// readers, assuming it has been thoughtfully chosen.
	//
	// Name does not need to be unique and can be empty.
	Name string

	// Description of the section, as provided by the user,
	// to aid users in understanding and using it.
	//
	// Description should compliment the Name field but is
	// also intended for general notes about the media
	// section. Description may be empty.
	Description string

	// TODO: start time
	// TODO: end time
}

func (vs VideoSection) _entity() {}
