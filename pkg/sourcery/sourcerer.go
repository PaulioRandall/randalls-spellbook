package sourcery

import (
	"maps"
	"net/http"

	"github.com/crgimenes/glaze"
)

type Codex = map[string]Spell
type Atlas = map[string]http.Handler
type Litholog = map[string]RuneStone

/*
// server represents a HTTP handler and its path.
type server struct {
	path    string
	handler http.Handler
}
*/

// RuneStone implementations add functionality to a [World]
// when it is built. They can be added via the
// [Sourcerer.Fuse] function.
type RuneStone interface {
	Bind(w *World)
	Free()
}

type WorldBuilder struct {
	options AppOptions
	stones  Litholog
	portals Atlas
}

func Summon() *WorldBuilder {
	return &WorldBuilder{
		options: AppOptions{
			Hint: glaze.HintNone,
		},
		stones:  Litholog{},
		portals: Atlas{},
	}
}

func (wb *WorldBuilder) Debug(
	state bool,
) *WorldBuilder {
	wb.options.Debug = state
	return wb
}

func (wb *WorldBuilder) Name(
	name string,
) *WorldBuilder {
	wb.options.Title = name
	return wb
}

func (wb *WorldBuilder) Resize(
	width, height int,
) *WorldBuilder {
	wb.options.Width = width
	wb.options.Height = height
	return wb
}

func (wb *WorldBuilder) AddEnchant(
	name string,
	rs RuneStone,
) *WorldBuilder {
	wb.stones[name] = rs
	return wb
}

func (wb *WorldBuilder) AddPortal(
	path string,
	handler http.Handler,
) *WorldBuilder {
	wb.portals[path] = handler
	return wb
}

func (wb *WorldBuilder) Conjure() *World {
	op := wb.options
	op.Handler = marryPortals(wb.portals)

	stones := maps.Clone(wb.stones)
	spellbook := compileSpellbook(stones)

	return buildWorld(op, stones, spellbook)
}

func marryPortals(
	portals map[string]http.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	for path, handler := range portals {
		mux.Handle(path, handler)
	}

	return mux
}

func compileSpellbook(
	stones map[string]RuneStone,
) Codex {
	var book Codex
	// TODO
	// + Get all methods
	// + Add methods as spells, except Bind and Free
	// + In the form 'RuneStoneName.Spellname'

	return book
}
