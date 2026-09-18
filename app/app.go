package app

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm"
)

type Datastore struct {
	w  *sourcery.World
	db *storm.Storm
}

func (ds Datastore) Bind(w *sourcery.World) {
	// TODO: Allow user to specify DB path.
	ds.w = w
	ds.db = storm.New("./testproject/data.sqlite")
	ds.db.Create(
		Media{},
	)
}

func (ds Datastore) Free() {
	if ds.db != nil {
		ds.db.Close()
		ds.db = nil
	}
}

func (ds Datastore) ListMedia() ([]Media, error) {
	return ds.db.Select(Media{})
}

func (ds Datastore) AddMedia(media Media) (Media, error) {
	media, e := media.Clean()
	if e == nil {
		e = ds.db.Insert(media)
	}
	return media, e
}

func (ds Datastore) GetMediaById(id string) (Media, error) {
	return ds.db.SelectById(Media{}, id)
}

func (ds Datastore) DeleteMediaById(id string) error {
	return ds.db.DeleteById(Media{}, id)
}

var _ sourcery.RuneStone = Datastore{}
