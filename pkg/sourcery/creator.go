package sourcery

import (
	"maps"
	"net/http"
	"reflect"

	"github.com/crgimenes/glaze"
)

type Portal interface {
	Open(w *World)
	Close()
}
type PortalMap = map[string]Portal
type HandlerMap = map[string]http.Handler
type Spellbook = map[string]Spell

// ********************************************************

type httpOnlyPortal struct {
	handler http.Handler
}

func (httpOnlyPortal) Open(w *World) {}
func (httpOnlyPortal) Close()        {}
func (hop httpOnlyPortal) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	hop.handler.ServeHTTP(w, r)
}

// ********************************************************

type Creator struct {
	options AppOptions
	pm      PortalMap
	hm      HandlerMap
}

func NewCreator() *Creator {
	return &Creator{
		options: AppOptions{
			Hint: glaze.HintNone,
		},
		pm: PortalMap{},
		hm: HandlerMap{},
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

func (wb *Creator) AddPortal(
	name string,
	portal Portal,
) *Creator {
	wb.pm[name] = portal
	return wb
}

func (wb *Creator) AddServer(
	path string,
	handler http.Handler,
) *Creator {
	if _, ok := handler.(Portal); ok {
		wb.hm[path] = handler
	} else {
		wb.hm[path] = httpOnlyPortal{
			handler: handler,
		}
	}
	return wb
}

func (wb *Creator) BuildWorld() *World {
	op := wb.options
	op.Handler = marryHandlers(wb.hm)
	spellbook := marrySpells(wb.pm)
	return buildWorld(op, wb.pm, spellbook)
}

func marryHandlers(hm HandlerMap) *http.ServeMux {
	mux := http.NewServeMux()

	for path, handler := range hm {
		mux.Handle(path, handler)
	}

	return mux
}

func marrySpells(pm PortalMap) Spellbook {
	book := Spellbook{}

	for _, port := range pm {
		spells := deriveSpellsFromPortal(port)
		maps.Copy(book, spells)
	}

	return book
}

func deriveSpellsFromPortal(port Portal) Spellbook {
	portVal := reflect.ValueOf(port)
	portTyp := portVal.Type()
	book := Spellbook{}

	for i := 0; i < portVal.NumMethod(); i++ {
		funcVal := portVal.Method(i)
		funcTyp := portTyp.Method(i)

		if !funcTyp.IsExported() {
			continue
		}

		// Ignore Portal & Handler functions.
		n := funcTyp.Name
		if n == "Open" || n == "Close" || n == "ServeHttp" {
			continue
		}

		name := portTyp.Elem().Name() + "." + funcTyp.Name
		book[name] = NewSpell(
			name,
			funcVal.Interface(),
		)
	}

	return book
}
