package sourcery

import (
	"net/http"
	"strings"

	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

// Sourcerer is the main application.
type Sourcerer struct {
	debug     bool
	title     string
	width     int
	height    int
	servers   []server
	spellbook spellbook.Spellbook
}

type server struct {
	path    string
	handler http.Handler
}

// New returns a new empty Sourcerer.
func New() *Sourcerer {
	return &Sourcerer{}
}

// Debug sets the debug state. Default is false.
func (sour *Sourcerer) Debug(state bool) {
	sour.debug = state
}

// Title sets the application title. Default is
// "Technotelicomnicon".
func (sour *Sourcerer) Title(title string) {
	sour.title = title
}

// Size sets the initial application window size if not
// available through configuration. Default is
// 800 by 600.
func (sour *Sourcerer) Size(width, height int) {
	sour.width = width
	sour.height = height
}

// Serve adds a HTTP handler.
func (sour *Sourcerer) Serve(
	path string,
	handler http.Handler,
) {
	sour.servers = append(
		sour.servers,
		server{
			path:    path,
			handler: handler,
		},
	)
}

// Spellbook returns the Spellbook used for creating
// spells that are invoked by the frontend.
func (sour *Sourcerer) Spellbook() spellbook.Spellbook {
	return sour.spellbook
}

// Begin
func (sour *Sourcerer) Begin() {
	sour.setDefaults()

	// TODO: add servers, only if at least one exists
	// TODO: create app
	// TODO: configure frontend to invoke spells
}

// setDefaults applies defaults to unset user options.
func (sour *Sourcerer) setDefaults() {
	if sour.width <= 0 {
		sour.width = 800
	}

	if sour.height <= 0 {
		sour.width = 600
	}

	if strings.TrimSpace(sour.title) == "" {
		sour.title = "Technotelicomnicon"
	}
}
