package sourcery

import (
	"net/http"
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

type Initable interface {
	Init(w *World) func()
}

type App struct {
	world        World
	servers      map[string]http.Handler
	functions    map[string]any
	initialisers []Initable
}

func New() *App {
	return &App{
		world: World{
			title:   "Technotelicomnicon",
			handler: http.NewServeMux(),
		},
		servers:   map[string]http.Handler{},
		functions: map[string]any{},
	}
}

func (a *App) Debug() *App {
	a.world.debug = true
	return a
}

func (a *App) Title(title string) *App {
	a.world.title = title
	return a
}

func (a *App) Width(width int) *App {
	a.world.width = width
	return a
}

func (a *App) Height(height int) *App {
	a.world.height = height
	return a
}

func (a *App) AddEntity(entity any) *App {
	a.parseEntityFunctions(entity)
	return a
}

func (a *App) AddServer(path string, server http.Handler) *App {
	a.world.handler.Handle(path, server)
	return a
}

func (a *App) AddEntityServer(path string, entityServer http.Handler) *App {
	a.world.handler.Handle(path, entityServer)
	a.AddEntity(any(entityServer))
	return a
}

func (a *App) parseEntityFunctions(entity any) {
	val := reflect.ValueOf(entity)
	typ := val.Type()
	entityName := typ.Elem().Name()

	if initable, ok := entity.(Initable); ok {
		a.world.initialisers = append(a.world.initialisers, initable)
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

		a.world.functions[name] = f
	}
}

func (a *App) Start() error {
	return a.world.run()
}
