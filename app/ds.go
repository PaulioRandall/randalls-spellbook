package app

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
	"github.com/PaulioRandall/randalls-spellbook/pkg/storm"
)

type Datastore struct {
	path string
	w    *sourcery.World
	db   *storm.Storm
}

func (ds *Datastore) Open(w *sourcery.World) {
	// TODO: Allow user to specify DB path.
	ds.w = w
	ds.db = storm.New("./testproject/data.sqlite")
	ds.db.Create(
		Media{},
	)
}

func (ds *Datastore) Close() {
	if ds.db != nil {
		ds.db.Close()
		ds.db = nil
	}
}
