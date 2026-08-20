package project

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
)

// Project is the access to project data including
// project configuration.
type Project struct {
	store data.Store
}

// New creates and returns a new empty project.
func New(path string) *Project {
	return &Project{
		store: data.NewStore(path),
	}
}

// Path returns the file path to the project file, i.e.
// its datastore.
func (p *Project) Path() string {
	return p.store.Path()
}

// Open opens the project. The project datastore is opened,
// being created first if it doesn't exist. Project
// configuration is loaded.
func (p *Project) Open() error {
	return p.store.Open()
	// TODO: Load configuration.
}

// Close closes the datastore and cleans up resources.
func (p *Project) Close() error {
	return p.store.Close()
}

// AddMedia creates a new Media entity and adds it to the
// list of media. See data.MakeMedia for more
// information.
func (p *Project) AddMedia(
	mediaType string,
	name string,
	description string,
	localPath string,
) (data.Media, error) {
	empty := data.Media{}

	m, e := data.MakeMedia(
		mediaType,
		name,
		description,
		localPath,
	)

	if e != nil {
		return empty, e
	}

	e = p.store.InsertMedia(m)
	if e != nil {
		return empty, e
	}

	return m, nil
}

// ListMedia returns all media.
func (p Project) ListMedia() ([]data.Media, error) {
	return p.store.ListMedia()
}

// GetMediaById returns the media with the given entityId
// or nil if it doesn't exist.
func (p Project) GetMediaById(
	entityId data.EntityId,
) (data.Media, error) {
	return p.store.GetMediaById(entityId)
}
