package sourcery

import (
	"maps"
	"net/http"
	"reflect"

	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

type Initialiser interface {
	Init(w *World) func()
}

type Unitialiser interface {
	Unit() func()
}

type Options struct {
	debug        bool
	webview      glaze.WebView
	functions    map[string]any
	serveMux     *http.ServeMux
	initialisers []Initialiser
	unitialisers []Unitialiser
}

func New(debug bool) Options {
	return Options{
		debug:        debug,
		webview:      createWebView(debug),
		functions:    map[string]any{},
		serveMux:     nil,
		initialisers: nil,
		unitialisers: nil,
	}
}

func createWebView(debug bool) glaze.WebView {
	webview, e := glaze.New(debug)
	if e != nil {
		panic(e)
	}

	webview.SetTitle("Technotelicomnicon")
	webview.SetSize(
		800,
		600,
		glaze.HintNone,
	)

	return webview
}

func (ops Options) SetTitle(title string) Options {
	ops.webview.SetTitle(title)
	return ops
}

func (ops Options) SetSize(width, height int) Options {
	ops.webview.SetSize(
		width,
		height,
		glaze.HintNone,
	)
	return ops
}

func (ops Options) SetHtml(html string) Options {
	ops.webview.SetHtml(html)
	return ops
}

func (ops Options) AddEntity(entity any) Options {
	funcs := parseEntityFunctions(entity)
	maps.Copy(ops.functions, funcs)

	if init, ok := entity.(Initialiser); ok {
		ops.initialisers = append(ops.initialisers, init)
	}

	if unit, ok := entity.(Unitialiser); ok {
		ops.unitialisers = append(ops.unitialisers, unit)
	}

	return ops
}

func (ops Options) AddServer(path string, server http.Handler) Options {
	if ops.serveMux == nil {
		ops.serveMux = http.NewServeMux()
	}

	ops.serveMux.Handle(path, server)
	return ops
}

func (ops Options) AddEntityServer(path string, entityServer http.Handler) Options {
	ops = ops.AddServer(path, entityServer)
	return ops.AddEntity(any(entityServer))
}

func (ops Options) CreateWorld() *World {
	return newWorld(ops)
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
