package sourcery

import (
	"maps"
	"net/http"
	"reflect"
	"slices"

	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

type Initable interface {
	Init(w *World) func()
}

type App struct {
	options      AppOptions
	servers      map[string]http.Handler
	functions    map[string]any
	initialisers []Initable
}

func New() *App {
	return &App{
		options: AppOptions{
			Hint: glaze.HintNone,
		},
		servers:   map[string]http.Handler{},
		functions: map[string]any{},
	}
}

func (a *App) Debug() *App {
	a.options.Debug = true
	return a
}

func (a *App) Name(name string) *App {
	a.options.Title = name
	return a
}

func (a *App) Width(width int) *App {
	a.options.Width = width
	return a
}

func (a *App) Height(height int) *App {
	a.options.Height = height
	return a
}

func (a *App) AddEntity(entity any) *App {
	a.parseEntityFunctions(entity)
	return a
}

func (a *App) AddServer(path string, server http.Handler) *App {
	a.servers[path] = server
	return a
}

func (a *App) AddEntityServer(path string, entityServer http.Handler) *App {
	a.servers[path] = entityServer
	return a.AddEntity(any(entityServer))
}

func (a *App) parseEntityFunctions(entity any) {
	val := reflect.ValueOf(entity)
	typ := val.Type()
	entityName := typ.Elem().Name()

	if initable, ok := entity.(Initable); ok {
		a.initialisers = append(a.initialisers, initable)
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

		a.functions[name] = f
	}
}

func (a *App) Start() error {
	options := a.options
	options.Handler = a.marryServers()

	w := &World{
		options:      options,
		functions:    maps.Clone(a.functions),
		initialisers: slices.Clone(a.initialisers),
	}

	return w.run()
}

func (a *App) marryServers() *http.ServeMux {
	mux := http.NewServeMux()

	for path, server := range a.servers {
		mux.Handle(path, server)
	}

	return mux
}
