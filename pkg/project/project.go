package project

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/datastore"
	"github.com/PaulioRandall/randalls-spellbook/pkg/datastore/sqlite"
	"github.com/PaulioRandall/randalls-spellbook/pkg/entity"
)

// Project is the access to project data including
// project configuration.
type Project struct {
	ds datastore.Datastore
}

// New creates and returns a new empty project.
func New(path string) *Project {
	return &Project{
		ds: sqlite.New(path),
	}
}

// Path returns the file path to the project file, i.e.
// its datastore.
func (p *Project) Path() string {
	return p.ds.Path()
}

// Open opens the project. The project datastore is opened,
// being created first if it doesn't exist. Project
// configuration is loaded.
func (p *Project) Open() error {
	return p.ds.Open()
	// TODO: Load configuration.
}

// Close closes the datastore and cleans up resources.
func (p *Project) Close() error {
	return p.ds.Close()
}

// AddMedia creates a new Media entity and adds it to the
// list of media. See entity.MakeMedia for more
// information.
func (p *Project) AddMedia(
	mediaType entity.MediaType,
	name string,
	description string,
	localPath string,
) (entity.Media, error) {
	empty := entity.Media{}

	m, e := entity.MakeMedia(
		mediaType,
		name,
		description,
		localPath,
	)

	if e != nil {
		return empty, e
	}

	e = p.ds.InsertMedia(m)
	if e != nil {
		return empty, e
	}

	return m, nil
}

// GetAllMedia returns all media.
func (p Project) GetAllMedia() ([]entity.Media, error) {
	return p.ds.GetAllMedia()
}

// GetMediaById returns the media with the given entityId
// or nil if it doesn't exist.
func (p Project) GetMediaById(
	entityId entity.EntityId,
) (entity.Media, error) {
	return p.ds.GetMediaById(entityId)
}
