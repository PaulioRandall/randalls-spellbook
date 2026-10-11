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

func (ds *Datastore) Init(w *sourcery.World) error {
	var e error

	// TODO: Allow user to specify DB path.
	ds.w = w
	ds.db, e = stormy.Open("./testproject/data.sqlite")
	if e != nil {
		return e
	}

	return ds.db.Create(Media{})
}

func (ds *Datastore) Unit() error {
	var e error

	if ds.db != nil {
		e = ds.db.Close()
		ds.db = nil
	}

	return e
}
