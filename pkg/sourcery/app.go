package sourcery

import (
	"net/http"
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

func New() *World {
	return &World{
		title:     "Technotelicomnicon",
		handler:   http.NewServeMux(),
		functions: map[string]any{},
	}
}

func (w *World) Debug() *World {
	w.debug = true
	return w
}

func (w *World) Title(title string) *World {
	w.title = title
	return w
}

func (w *World) Width(width int) *World {
	w.width = width
	return w
}

func (w *World) Height(height int) *World {
	w.height = height
	return w
}

func (w *World) AddEntity(entity any) *World {
	w.parseEntityFunctions(entity)
	return w
}

func (w *World) AddServer(path string, server http.Handler) *World {
	w.handler.Handle(path, server)
	return w
}

func (w *World) AddEntityServer(path string, entityServer http.Handler) *World {
	w.handler.Handle(path, entityServer)
	w.AddEntity(any(entityServer))
	return w
}

func (w *World) parseEntityFunctions(entity any) {
	val := reflect.ValueOf(entity)
	typ := val.Type()
	entityName := typ.Elem().Name()

	if initable, ok := entity.(Initable); ok {
		w.initialisers = append(w.initialisers, initable)
	}

	for i := 0; i < val.NumMethod(); i++ {
		funcTyp := typ.Method(i)

		if !funcTyp.IsExported() {
			continue
		}

		n := funcTyp.Name
		if n == "Init" || n == "ServeHttp" {
			continue
		}

		name := entityName + "." + funcTyp.Name
		f := val.Method(i).Interface()

		e := ValidateValErrFunc(f)
		if e != nil {
			e = sin.Fmt(
				"'%s' has invalid function signature",
				name,
			).Wrap(e)
			panic(e)
		}

		w.functions[name] = f
	}
}
