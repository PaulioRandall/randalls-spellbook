package sourcery

import (
	"net/http"

	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

// Sourcerer is the main application.
type Sourcerer struct {
	spellbook.Spellbook
	debug   bool
	title   string
	width   int
	height  int
	servers []server
}

type server struct {
	path    string
	handler http.Handler
}

// NewSourcerer returns a new empty Sourcerer that can be
// used to conjure an application.
func NewSourcerer() *Sourcerer {
	return &Sourcerer{
		title:  "Technotelicomnicon",
		width:  800,
		height: 600,
	}
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

// Conjure starts the application.
func (sour *Sourcerer) Conjure() error {

	// TODO: add servers, only if at least one exists
	// TODO: create app
	// TODO: configure frontend to invoke spells

	return nil
}
