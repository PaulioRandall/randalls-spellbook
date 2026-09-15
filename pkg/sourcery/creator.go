package sourcery

import (
	"maps"
	"net/http"
	"reflect"

	"github.com/crgimenes/glaze"
)

/*
// server represents a HTTP handler and its path.
type server struct {
	path    string
	handler http.Handler
}
*/

type RuneStone interface {
	Bind(w *World)
	Free()
}

type Creator struct {
	options AppOptions
	stones  map[string]RuneStone
	portals map[string]http.Handler
}

func New() *Creator {
	return &Creator{
		options: AppOptions{
			Hint: glaze.HintNone,
		},
		stones:  map[string]RuneStone{},
		portals: map[string]http.Handler{},
	}
}

func (wb *Creator) Debug() *Creator {
	wb.options.Debug = true
	return wb
}

func (wb *Creator) Name(
	name string,
) *Creator {
	wb.options.Title = name
	return wb
}

func (wb *Creator) Size(
	width, height int,
) *Creator {
	wb.options.Width = width
	wb.options.Height = height
	return wb
}

func (wb *Creator) Bind(
	name string,
	rs RuneStone,
) *Creator {
	wb.stones[name] = rs
	return wb
}

func (wb *Creator) Serve(
	path string,
	handler http.Handler,
) *Creator {
	wb.portals[path] = handler
	return wb
}

func (wb *Creator) BuildWorld() *World {
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
) map[string]Spell {
	book := map[string]Spell{}

	for _, rs := range stones {
		spells := deriveSpells(rs)
		maps.Copy(book, spells)
	}

	return book
}

func deriveSpells(rs RuneStone) map[string]Spell {
	rsVal := reflect.ValueOf(rs)
	rsTyp := rsVal.Type()

	chapter := map[string]Spell{}

	for i := 0; i < rsVal.NumMethod(); i++ {
		funcVal := rsVal.Method(i)
		funcTyp := rsTyp.Method(i)

		if !funcTyp.IsExported() {
			continue
		}

		if funcTyp.Name == "Bind" || funcTyp.Name == "Free" {
			continue
		}

		name := rsTyp.Elem().Name() + "." + funcTyp.Name
		chapter[name] = NewSpell(
			name,
			funcVal.Interface(),
		)
	}

	return chapter
}
