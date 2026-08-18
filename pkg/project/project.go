package project

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/entity"
)

type Project struct {
	media []entity.Media
}

// New creates and returns a new empty project.
func New() *Project {
	return &Project{}
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
	m, e := entity.MakeMedia(
		mediaType,
		name,
		description,
		localPath,
	)

	if e == nil {
		p.media = append(p.media, m)
	}

	return m, e
}

// GetAllMedia returns all media.
func (p Project) GetAllMedia() []entity.Media {
	return p.media
}

// GetMediaById returns the media with the given EntityId
// or nil if it doesn't exist.
func (p Project) GetMediaById(entityId entity.EntityId) entity.Media {
	return entity.FindMediaById(p.media, entityId)
}
