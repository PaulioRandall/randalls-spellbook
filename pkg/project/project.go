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

// AddVideo creates a new video entity and adds it to the
// list of videos. See entity.MakeVideo for more
// information.
func (p *Project) AddVideo(
	name, description, localPath string,
) (entity.Media, error) {
	v, e := entity.MakeVideo(name, description, localPath)

	if e == nil {
		p.media = append(p.media, v)
	}

	return v, e
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
