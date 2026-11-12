package business

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
	"github.com/PaulioRandall/randalls-spellbook/pkg/stormy"
)

type Project struct {
	path           string
	db             *stormy.Stormy
	currentMediaId string
}

func init() {
	p := &Project{}

	sourcery.Append(
		sourcery.OptionServer("/media/", p),
		sourcery.OptionNamedEntity("Project", p),
		sourcery.OptionListener(p),
	)
}

func (p *Project) OnEvent(ev sourcery.Event) {
	switch ev.Type() {
	case sourcery.EventTypeUnload: // Clean up.
		if e := p.Close(); e != nil {
			panic(e)
		}
	}
}

func (p *Project) Close() error {
	if p.db == nil {
		return nil
	}

	defer func() {
		p.path = ""
		p.db = nil
	}()

	return p.db.Close()
}

func (p *Project) Open(path string) error {
	if p.db != nil {
		p.Close()
	}

	db, e := stormy.Open(path)

	if e != nil {
		return e
	}

	p.path = path
	p.db = db
	return nil
}

func (p *Project) Path() string {
	return p.path
}

var _ sourcery.Listener = &Project{}
