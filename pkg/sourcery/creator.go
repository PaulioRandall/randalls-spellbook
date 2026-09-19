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

type HttpPortal interface {
	Portal
	http.Handler
}
type HttpPortalMap = map[string]HttpPortal

type GoPortal interface{ Portal }
type GoPortalMap = map[string]GoPortal
type Spellbook = map[string]Spell

type statelessHttpPortal struct {
	handler http.Handler
}

func (hp statelessHttpPortal) Open(w *World) {
	type opener interface{ Open(w *World) }
	if portal, ok := hp.handler.(opener); ok {
		portal.Open(w)
	}
}
func (hp statelessHttpPortal) Close() {
	type closer interface{ Close() }
	if portal, ok := hp.handler.(closer); ok {
		portal.Close()
	}
}
func (hp statelessHttpPortal) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	hp.handler.ServeHTTP(w, r)
}

type Creator struct {
	options AppOptions
	gpm     GoPortalMap
	hpm     HttpPortalMap
}

func New() *Creator {
	return &Creator{
		options: AppOptions{
			Hint: glaze.HintNone,
		},
		gpm: GoPortalMap{},
		hpm: HttpPortalMap{},
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

func (wb *Creator) GoPortal(
	name string,
	portal GoPortal,
) *Creator {
	wb.gpm[name] = portal
	return wb
}

func (wb *Creator) HttpPortal(
	path string,
	portal http.Handler,
) *Creator {
	if hp, ok := portal.(HttpPortal); ok {
		wb.hpm[path] = hp
	} else {
		wb.hpm[path] = statelessHttpPortal{
			handler: portal,
		}
	}
	return wb
}

func (wb *Creator) BuildWorld() *World {
	op := wb.options
	op.Handler = httpPortalMux(wb.hpm)
	portals := marryPortals(wb.hpm, wb.gpm)
	spellbook := conjureSpellbook(wb.gpm)
	return buildWorld(op, portals, spellbook)
}

func httpPortalMux(hpm HttpPortalMap) *http.ServeMux {
	mux := http.NewServeMux()

	for path, handler := range hpm {
		mux.Handle(path, handler)
	}

	return mux
}

func conjureSpellbook(gpm GoPortalMap) Spellbook {
	book := Spellbook{}

	for _, gp := range gpm {
		spells := deriveSpells(gp)
		maps.Copy(book, spells)
	}

	return book
}

func deriveSpells(gp GoPortal) Spellbook {
	gpVal := reflect.ValueOf(gp)
	gpTyp := gpVal.Type()

	chapter := Spellbook{}

	for i := 0; i < gpVal.NumMethod(); i++ {
		funcVal := gpVal.Method(i)
		funcTyp := gpTyp.Method(i)

		if !funcTyp.IsExported() {
			continue
		}

		if funcTyp.Name == "Open" || funcTyp.Name == "Close" {
			continue
		}

		name := gpTyp.Elem().Name() + "." + funcTyp.Name
		chapter[name] = NewSpell(
			name,
			funcVal.Interface(),
		)
	}

	return chapter
}

func marryPortals(
	hpm HttpPortalMap,
	gpm GoPortalMap,
) []Portal {
	size := len(hpm) + len(gpm)
	portals := make([]Portal, size, size)
	index := 0

	for _, hp := range hpm {
		portals[index] = hp
		index++
	}

	for _, gp := range gpm {
		portals[index] = gp
		index++
	}

	return portals
}
