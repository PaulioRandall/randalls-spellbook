package sourcery

import (
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
type LifecycleHandler func(rm *Realm) error

// Realm is the main application structure.
type Realm struct {
	spellbook.Spellbook
	debug        bool
	title        string
	width        int
	height       int
	servers      []server
	webview      glaze.WebView
	afterOpening LifecycleHandler
	afterClosing LifecycleHandler
}

// NewRealm returns a new Realm with default values.
func NewRealm() *Realm {
	return &Realm{
		title:  "Technotelicomnicon",
		width:  800,
		height: 600,
	}
}

// Debug sets the debug state. Default is false.
func (rm *Realm) Debug(state bool) {
	rm.debug = state
}

// Title sets the Realm title. Default is
// "Technotelicomnicon".
func (rm *Realm) Title(title string) {
	rm.title = title
}

// Size sets the initial realm window size if not available
// through configuration. Default is 800 by 600.
func (rm *Realm) Size(width, height int) {
	rm.width = width
	rm.height = height
}

// Serve adds a HTTP handler.
func (rm *Realm) Serve(
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
func (rm *Realm) AfterOpening(f LifecycleHandler) {
	rm.afterOpening = f
}

// AfterClosing sets a handler that is called after the
// Realm has closed.
func (rm *Realm) AfterClosing(f LifecycleHandler) {
	rm.afterClosing = f
}

// lifecycleChange invokes a LifecycleHandler function if
// it's not nil.
func (rm *Realm) onLifecycleEvent(
	f LifecycleHandler,
) error {
	if f != nil {
		return f(rm)
	}
	return nil
}

// OpenPortal creates and enters the Realm, i.e. starts the
// application and blocks until the application exits.
func (rm *Realm) OpenPortal() error {
	handler := rm.createMuxServer()
	appOptions := rm.createRealmOptions(handler)
	return rm.runWebview(appOptions)
}

// WebView returns the Glaze WebView or nil if not set.
func (rm *Realm) WebView() glaze.WebView {
	return rm.webview
}

// createMuxServer creates a mux handler for all servers
// currently set in the Realm.
func (rm *Realm) createMuxServer() *http.ServeMux {
	mux := http.NewServeMux()

	for _, server := range rm.servers {
		mux.Handle(server.path, server.handler)
	}

	return mux
}

// createRealmOptions creates the options for Glaze
// WebView.
func (rm *Realm) createRealmOptions(
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
func (rm *Realm) onWebViewReady(
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
func (rm *Realm) castSpell(
	spellname string,
	data any,
) (any, error) {
	effect := rm.Cast(spellname, rm, data)
	if effect.Cursed() {
		return nil, effect.Error()
	}
	return effect.Result(), nil
}

// runWebview runs the WebView (blocking).
func (rm *Realm) runWebview(
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
