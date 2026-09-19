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

func (ds *Datastore) WorldEnter(w *sourcery.World) {
	// TODO: Allow user to specify DB path.
	ds.w = w
	ds.db = storm.New("./testproject/data.sqlite")

	e := ds.db.Open()
	if e != nil {
		panic(e)
	}

	e = ds.db.Create(Media{})
	if e != nil {
		panic(e)
	}
}

func (ds *Datastore) WorldExit() {
	if ds.db != nil {
		ds.db.Close()
		ds.db = nil
	}
}

var _ sourcery.Portal = &Datastore{}
