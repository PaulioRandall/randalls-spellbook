package app

import (
	"github.com/google/uuid"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
)

type Util struct {
	w      *sourcery.World
	dbPath string
}

func (u *Util) Init(w *sourcery.World) error {
	u.w = w
	// TODO: Allow user to specify DB path.
	u.dbPath = "./testproject/data.sqlite"
	return nil
}

func (u *Util) Unit() error {
	u.w = nil
	return nil
}

func (u *Util) SelectLocalFile(
	title string,
) (string, error) {
	return "", sin.Err("Not yet implemented!")
}

func randomEntityId() string {
	return uuid.New().String()
}

var (
	_ sourcery.Initialiser = &Util{}
	_ sourcery.Unitialiser = &Util{}
)
