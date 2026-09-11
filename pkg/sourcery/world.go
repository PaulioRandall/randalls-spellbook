package sourcery

import (
	"fmt"
	"net/http"

	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/wizzard"
)

/*
// server represents a HTTP handler and its path.
type server struct {
	path    string
	handler http.Handler
}
*/

type World struct {
	Debug   bool
	Title   string
	Width   int
	Height  int
	spells  map[string]wizzard.Spell
	servers []server
	webview glaze.WebView
}

func NewWorld() *World {
	return &World{
		Title:  "Technotelicomnicon",
		Width:  800,
		Height: 600,
		spells: map[string]wizzard.Spell{},
	}
}

// Inscribe adds a new spell to the world. The spell may be
// called from the front end using the 'Go' global
// function, e.g. Go("name", arg1, arg2, etc).
func (w *World) Inscribe(name string, fn any) {
	w.spells[name] = wizzard.NewSpell(name, fn)
}

// Serve adds a HTTP handler.
func (w *World) Serve(
	path string,
	handler http.Handler,
) {
	w.servers = append(
		w.servers,
		server{
			path:    path,
			handler: handler,
		},
	)
}

// WebView returns the Glaze WebView or nil if not set.
func (w *World) WebView() glaze.WebView {
	return w.webview
}

// OpenPortal creates and enters the Realm, i.e. starts the
// application and blocks until the application exits.
func (w *World) OpenPortal() error {
	handler := w.createMuxServer()
	options := w.createAppOptions(handler)
	return AppWindow(options)
}

// createMuxServer creates a mux handler for all servers
// currently set in the Realm.
func (w *World) createMuxServer() *http.ServeMux {
	mux := http.NewServeMux()

	for _, server := range w.servers {
		mux.Handle(server.path, server.handler)
	}

	return mux
}

// createAppOptions creates the options for Glaze WebView.
func (w *World) createAppOptions(
	handler http.Handler,
) AppOptions {
	return AppOptions{
		Debug:          w.Debug,
		Title:          w.Title,
		Width:          w.Width,
		Height:         w.Height,
		Hint:           glaze.HintNone,
		Handler:        handler,
		OnWebViewReady: w.onWebViewReady,
	}
}

// onWebViewReady is a callback for when the Realm is
// opened and ready for functions to be bound.
func (w *World) onWebViewReady(
	wv glaze.WebView,
) error {
	w.webview = wv
	return wv.Bind("Go", w.Go)
}

// Go is the entry point for synchronus requests from the
// front end.
func (w *World) Go(
	cmd string,
	args ...any,
) (any, error) {
	println("Go: " + cmd)
	for _, v := range args {
		println(fmt.Sprintf("\t%v", v))
	}
	return nil, fmt.Errorf("%s", cmd)
}
