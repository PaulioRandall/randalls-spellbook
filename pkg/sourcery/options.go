package sourcery

import (
	"maps"
	"net/http"
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

type Initable interface {
	Init(w *World) func()
}

type Options struct {
	Debug     bool
	Title     string
	Width     int
	Height    int
	ServeMux  *http.ServeMux
	Functions map[string]any
}

func NewOptions() Options {
	return Options{
		Title:     "Technotelicomnicon",
		Width:     800,
		Height:    600,
		ServeMux:  http.NewServeMux(),
		Functions: map[string]any{},
	}
}

func (ops Options) EnableDebug() Options {
	ops.Debug = true
	return ops
}

func (ops Options) SetTitle(title string) Options {
	ops.Title = title
	return ops
}

func (ops Options) SetWidth(width int) Options {
	ops.Width = width
	return ops
}

func (ops Options) SetHeight(height int) Options {
	ops.Height = height
	return ops
}

func (ops Options) AddEntity(entity any) Options {
	funcs := parseEntityFunctions(entity)
	maps.Copy(ops.Functions, funcs)
	return ops
}

func (ops Options) AddServer(path string, server http.Handler) Options {
	ops.ServeMux.Handle(path, server)
	return ops
}

func (ops Options) AddEntityServer(path string, entityServer http.Handler) Options {
	ops.ServeMux.Handle(path, entityServer)
	return ops.AddEntity(any(entityServer))
}

func parseEntityFunctions(entity any) map[string]any {
	val := reflect.ValueOf(entity)
	typ := val.Type()
	entityName := typ.Elem().Name()

	results := map[string]any{}

	for i := 0; i < val.NumMethod(); i++ {
		funcTyp, allowed := getEntityMethodAt(typ, i)
		if !allowed {
			continue
		}

		name := entityName + "." + funcTyp.Name
		f := val.Method(i).Interface()
		checkEntityMethod(f, name)

		results[name] = f
	}

	return results
}

func getEntityMethodAt(typ reflect.Type, i int) (reflect.Method, bool) {
	funcTyp := typ.Method(i)

	if !funcTyp.IsExported() {
		// Ignore unexported because we can only call exported
		// methods when external to the object.
		return reflect.Method{}, false
	}

	n := funcTyp.Name
	if n == "Init" || n == "ServeHttp" {
		// Ignore Init and ServeHttp named functions as they
		// serve special purposes.
		return reflect.Method{}, false
	}

	return funcTyp, true
}

func checkEntityMethod(f any, name string) {
	e := ValidateValErrFunc(f)
	if e == nil {
		return
	}

	sin.Fmt(
		"'%s' has invalid function signature",
		name,
	).
		Wrap(e).
		Panic()
}
