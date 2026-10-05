package app

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
	"github.com/PaulioRandall/randalls-spellbook/pkg/stormy"
)

type Datastore struct {
	path string
	w    *sourcery.World
	db   *stormy.Stormy
}

func (ds *Datastore) Init(w *sourcery.World) func() {
	var e error

	// TODO: Allow user to specify DB path.
	ds.w = w
	ds.db, e = stormy.Open("./testproject/data.sqlite")
	if e != nil {
		panic(e)
	}

	e = ds.db.Create(Media{})
	if e != nil {
		panic(e)
	}

	return ds.free
}

func (ds *Datastore) free() {
	if ds.db != nil {
		ds.db.Close()
		ds.db = nil
	}
}
