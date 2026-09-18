package api

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/data2"
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
		data2.Media{},
	)
}

func (ds Datastore) Free() {
	if ds.db != nil {
		ds.db.Close()
		ds.db = nil
	}
}

func (ds Datastore) ListMedia() ([]data2.Media, error) {
	return ds.db.Select(data2.Media{})
}

var _ sourcery.RuneStone = Datastore{}
