package app

import (
	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
)

type App struct {
	w      *sourcery.World
	dbPath string
}

func (app *App) Open(w *sourcery.World) {
	app.w = w
	// TODO: Allow user to specify DB path.
	app.dbPath = "./testproject/data.sqlite"
}

func (app *App) Close() {
	app.w = nil
}

func (app *App) SelectLocalFile(
	title string,
) (string, error) {
	// Blocks!
	return app.w.WebView().OpenFile(glaze.FileDialogOptions{
		Title: title,
	})
}
