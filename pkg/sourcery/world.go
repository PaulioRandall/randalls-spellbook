package sourcery

import (
	"fmt"
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

type World struct {
	debug        bool
	webview      glaze.WebView
	functions    map[string]any
	serveMux     *http.ServeMux
	initialisers []Initialiser
	unitialisers []Unitialiser
}

func NewWorld(debug bool) *World {
	return &World{
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

func (w *World) SetTitle(title string) *World {
	w.webview.SetTitle(title)
	return w
}

func (w *World) SetSize(width, height int) *World {
	w.webview.SetSize(
		width,
		height,
		glaze.HintNone,
	)
	return w
}

func (w *World) SetHtml(html string) *World {
	w.webview.SetHtml(html)
	return w
}

func (w *World) AddEntity(entity any) *World {
	funcs := parseEntityFunctions(entity)
	maps.Copy(w.functions, funcs)

	if init, ok := entity.(Initialiser); ok {
		w.initialisers = append(w.initialisers, init)
	}

	if unit, ok := entity.(Unitialiser); ok {
		w.unitialisers = append(w.unitialisers, unit)
	}

	return w
}

func (w *World) AddServer(path string, server http.Handler) *World {
	if w.serveMux == nil {
		w.serveMux = http.NewServeMux()
	}

	w.serveMux.Handle(path, server)
	return w
}

func (w *World) AddEntityServer(path string, entityServer http.Handler) *World {
	w = w.AddServer(path, entityServer)
	return w.AddEntity(any(entityServer))
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

func (w *World) Start() (e error) {
	var baseUrl string

	if w.serveMux != nil {
		cs, e := createContentServer(w.serveMux)
		if e != nil {
			return e
		}

		go func() {
			err := cs.serve()
			if e == nil {
				e = err
			}
		}()

		defer func() {
			err := cs.close()
			if e == nil {
				e = err
			}
		}()

		baseUrl = cs.baseUrl
	}

	defer func() {
		w.webview.Destroy()
		w.webview = nil
	}()

	if baseUrl != "" {
		w.webview.Navigate(baseUrl)
	}

	e = w.bindGoRaw()
	if e != nil {
		return e
	}

	defer w.unitialise()
	w.initialise()

	// Starts the webview and blocks until it exits.
	w.webview.Run()

	return nil
}

func (w *World) initialise() {
	for _, v := range w.initialisers {
		v.Init(w)
	}
}

func (w *World) unitialise() {
	for _, v := range w.unitialisers {
		v.Unit()
	}
}

func (w *World) WebView() glaze.WebView {
	return w.webview
}

func (w *World) Exit() {
	w.WebView().Terminate()
}

func (w *World) bindGoRaw() error {
	e := w.webview.Bind("GoRaw", w.GoRaw)
	if e != nil {
		return e
	}

	w.webview.Init(`
		var Go = function(funcName, ...args) {
			return GoRaw(
				funcName,
				JSON.stringify(args),
			)
		}
	`)

	return nil
}

func (w *World) GoRaw(
	funcName string,
	jsonArgs string,
) (any, error) {
	if f, ok := w.functions[funcName]; ok {
		w.Log("Go: %s", funcName)
		thunk, e := WithJsonArgs(f, jsonArgs)

		if e != nil {
			return nil, e
		}

		return thunk.Call()
	}

	w.Log("Unknown function: %s", funcName)
	return nil, fmt.Errorf("Unknown function: %s", funcName)
}

func (w *World) Log(msg string, args ...any) {
	if !w.debug {
		return
	}

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Printf("[Sourcery] %s\n", msg)
}
