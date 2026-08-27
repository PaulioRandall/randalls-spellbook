package sourcery

import (
	"errors"
	"log"
	"net/http"

	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

// server represents a HTTP handler and its path.
type server struct {
	path    string
	handler http.Handler
}

// LifecycleHandler is a function called when the
// application lifecycle change occurs, i.e. the Realm is
// opened or destroyed.
type LifecycleHandler[T any] func(rm *Realm[T]) error

// Realm is the main application structure.
type Realm[T any] struct {
	Inventory    T
	spellbook    *spellbook.Spellbook
	debug        bool
	title        string
	width        int
	height       int
	servers      []server
	webview      glaze.WebView
	afterOpening LifecycleHandler[T]
	afterClosing LifecycleHandler[T]
}

// NewRealm returns a new Realm with default values.
func NewRealm[T any]() *Realm[T] {
	return &Realm[T]{
		spellbook: spellbook.Conjure(),
		title:     "Technotelicomnicon",
		width:     800,
		height:    600,
	}
}

// Spellbook returns the realm's Spellbook.
func (rm *Realm[T]) Spellbook() *spellbook.Spellbook {
	return rm.spellbook
}

// Debug sets the debug state. Default is false.
func (rm *Realm[T]) Debug(state bool) {
	rm.debug = state
}

// Title sets the Realm title. Default is
// "Technotelicomnicon".
func (rm *Realm[T]) Title(title string) {
	rm.title = title
}

// Size sets the initial realm window size if not available
// through configuration. Default is 800 by 600.
func (rm *Realm[T]) Size(width, height int) {
	rm.width = width
	rm.height = height
}

// Serve adds a HTTP handler.
func (rm *Realm[T]) Serve(
	path string,
	handler http.Handler,
) {
	rm.servers = append(
		rm.servers,
		server{
			path:    path,
			handler: handler,
		},
	)
}

// AfterOpening sets a handler that is called when the
// Realm has started.
func (rm *Realm[T]) AfterOpening(f LifecycleHandler[T]) {
	rm.afterOpening = f
}

// AfterClosing sets a handler that is called after the
// Realm has closed.
func (rm *Realm[T]) AfterClosing(f LifecycleHandler[T]) {
	rm.afterClosing = f
}

// lifecycleChange invokes a LifecycleHandler function if
// it's not nil.
func (rm *Realm[T]) onLifecycleEvent(
	f LifecycleHandler[T],
) error {
	if f != nil {
		return f(rm)
	}
	return nil
}

// OpenPortal creates and enters the Realm, i.e. starts the
// application and blocks until the application exits.
func (rm *Realm[T]) OpenPortal() error {
	handler := rm.createMuxServer()
	appOptions := rm.createRealmOptions(handler)
	return rm.runWebview(appOptions)
}

// WebView returns the Glaze WebView or nil if not set.
func (rm *Realm[T]) WebView() glaze.WebView {
	return rm.webview
}

// createMuxServer creates a mux handler for all servers
// currently set in the Realm.
func (rm *Realm[T]) createMuxServer() *http.ServeMux {
	mux := http.NewServeMux()

	for _, server := range rm.servers {
		mux.Handle(server.path, server.handler)
	}

	return mux
}

// createRealmOptions creates the options for Glaze
// WebView.
func (rm *Realm[T]) createRealmOptions(
	handler http.Handler,
) AppOptions {
	return AppOptions{
		Debug:          rm.debug,
		Title:          rm.title,
		Width:          rm.width,
		Height:         rm.height,
		Hint:           glaze.HintNone,
		Handler:        handler,
		OnWebViewReady: rm.onWebViewReady,
	}
}

// onWebViewReady is a callback for when the Realm is
// opened and ready for functions to be bound.
func (rm *Realm[T]) onWebViewReady(
	w glaze.WebView,
) error {
	rm.webview = w

	e := rm.onLifecycleEvent(rm.afterOpening)
	if e != nil {
		return e
	}

	return w.Bind("CastSpell", rm.castSpell)
}

// castSpell is a bound function called by the UI to
// execute backend functionality.
func (rm *Realm[T]) castSpell(
	spellname string,
	data any,
) (result any, e error) {
	defer func() {
		var ok bool
		r := recover()

		if r == nil {
			return
		}

		if e, ok = r.(error); ok {
			return
		}

		if s, ok := r.(string); ok {
			e = errors.New(s)
			return
		}

		println("FOOBAR: panic with unknown type")
		panic(r)
	}()

	effect := rm.spellbook.Cast(spellname, rm, data)

	if effect.Cursed() {
		return nil, effect.Error()
	}

	return effect.Result(), nil
}

// runWebview runs the WebView (blocking).
func (rm *Realm[T]) runWebview(
	appOptions AppOptions,
) (e error) {
	defer func() {
		closeErr := rm.onLifecycleEvent(rm.afterClosing)

		if closeErr == nil {
			return
		}

		if e == nil {
			e = closeErr
			return
		}

		log.Println(closeErr)
	}()

	return AppWindow(appOptions)
}
