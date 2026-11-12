package sourcery

import (
	"fmt"
	"net/http"
	"strings"

	engine "github.com/PaulioRandall/go-webview"
	_ "github.com/PaulioRandall/go-webview/embedded"
)

// TODO
//
//	- Do I need sync anywhere?

type entityFuncs = map[string]any

type entityEntry struct {
	entity any
	funcs  entityFuncs
}

// App is the core of an applicaiton. It starts, manages,
// and closes the webview, backend, and server.
type App struct {
	debug     bool
	running   bool
	webview   engine.WebView
	server    http.Handler
	entities  map[string]entityEntry
	listeners []Listener
}

// DefaultApp constructs a new App from the
// [DefaultOptions] package object which is appended to via
// the [Append] function.
func DefaultApp() *App {
	return NewApp(DefaultOptions...)
}

// NewApp constructs a new App with the specified options.
// Options are applied in order they are provided so some
// options may replace values set by a prior one.
func NewApp(options ...Option) *App {
	c := newConfig(options)

	return &App{
		debug:     c.debug,
		running:   false,
		webview:   newWebView(c),
		server:    newServeMux(c),
		entities:  c.entities,
		listeners: c.listeners,
	}
}

func newConfig(options []Option) config {
	c := config{
		debug:    false,
		title:    "Technotelicomnicon",
		width:    800,
		height:   600,
		handlers: map[string]http.Handler{},
		entities: map[string]entityEntry{},
	}

	for _, applyOption := range options {
		applyOption(&c)
	}

	return c
}

func newWebView(c config) engine.WebView {
	webview := engine.New(c.debug)
	webview.SetTitle(c.title)
	webview.SetSize(c.width, c.height, engine.HintNone)
	return webview
}

func newServeMux(c config) *http.ServeMux {
	if len(c.handlers) == 0 {
		return nil
	}

	mux := http.NewServeMux()

	for path, handler := range c.handlers {
		mux.Handle(path, handler)
	}

	return mux
}

// Debug returns true if debug mode is enabled.
func (app *App) Debug() bool {
	return app.debug
}

// Running returns true if the app is currently running.
func (app *App) Running() bool {
	return app.running
}

// WebView returns the applications webview instance.
func (app *App) WebView() engine.WebView {
	return app.webview
}

// HasEntity returns true if an entity with the given name
// exists.
func (app *App) HasEntity(name string) bool {
	_, ok := app.entities[name]
	return ok
}

// Entity returns the entity for the given name. If no such
// entity exists or it cannot be cast to T then false is
// returned.
func (app *App) Entity[T any](name string) (T, bool) {
	var zero T

	entry, ok := app.entities[name]
	if !ok {
		return zero, false
	}

	typeEntity, ok := entry.entity.(T)
	if !ok {
		return zero, false
	}

	return typeEntity, ok
}

// Run starts the application, blocking until the
// application has exited.
func (app *App) Run() (e error) {
	if app.running {
		panic("App must be in the created state to run")
	}
	app.running = true

	var baseUrl string

	if app.server != nil {
		cs, e := createContentServer(app.server)
		if e != nil {
			return e
		}

		// Start server in a separate Go routine.
		go func() {
			err := cs.serve()
			if e == nil {
				e = err
			}
		}()

		// Trigger server close when application exits.
		defer func() {
			err := cs.close()
			if e == nil {
				e = err
			}
		}()

		baseUrl = cs.baseUrl
	}

	// TODO: Only navigate if user hasn't passed non-empty
	// OptionHtml (non-whitespace rather).
	if baseUrl != "" {
		app.webview.Navigate(baseUrl)
	}

	// TODO: Create init that provides general backend
	//       details as an object accessible from window.
	app.webview.Init("var SourceryServer = \"" + baseUrl + "\"")

	// Bind the function that the passes backend function
	// invocations from JS to Go.
	e = app.bindGoRaw()
	if e != nil {
		return e
	}

	// Call event listeners with EventTypeLoad.
	ev := &lifecycleEvent{
		eventType: EventTypeLoad,
		app:       app,
	}
	for _, l := range app.listeners {
		l.OnEvent(ev)
	}

	// Trigger webview clean up when application exits.
	defer func() {
		app.webview.Destroy()
		app.webview = nil
	}()

	// Trigger EventTypeUnload when application exits.
	defer func() {
		ev := &lifecycleEvent{
			eventType: EventTypeUnload,
			app:       app,
		}
		for _, l := range app.listeners {
			l.OnEvent(ev)
		}
	}()

	// Starts the webview and blocks until it exits.
	app.webview.Run()

	return nil
}

func (app *App) bindGoRaw() error {
	e := app.webview.Bind("GoRaw", app.goRaw)
	if e != nil {
		return e
	}

	app.webview.Init(`
		var Go = function(funcName, ...args) {
			return GoRaw(
				funcName,
				JSON.stringify(args),
			)
		}
	`)

	return nil
}

// goRaw is bound to the webview and the entry point for
// users JS function calls to the Go backend. The user
// actually calls 'Go(...)' in JS which Jsonifies the
// arguments before calling goRaw.
func (app *App) goRaw(
	funcPath string,
	jsonArgs string,
) (any, error) {
	returnErr := func(msg string, args ...any) (any, error) {
		fmt.Printf(msg+"\n", args...)
		return nil, fmt.Errorf(msg, args...)
	}

	parts := strings.Split(funcPath, ".")
	if len(parts) != 2 {
		return returnErr("Unknown function: %s", funcPath)
	}

	entityName := parts[0]
	funcName := parts[1]

	entry, ok := app.entities[entityName]
	if !ok {
		return returnErr("Unknown entity: %s", funcPath)
	}

	f, ok := entry.funcs[funcName]
	if !ok {
		return returnErr("Unknown function: %s", funcPath)
	}

	fmt.Printf("Go: %s\n", funcName)
	return call(f, jsonArgs)
}
