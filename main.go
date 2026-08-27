package main

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/PaulioRandall/randalls-spellbook/ui"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"

	_ "github.com/PaulioRandall/randalls-spellbook/pkg/effect"
)

func main() {
	rm := sourcery.NewRealm()

	rm.Debug(true)
	rm.Title("Randall's Spellbook")
	rm.Size(800, 600)
	rm.Serve("/media/", MediaServer{rm})
	rm.Serve("/", createWebServer())

	rm.AfterOpening(afterOpening)
	rm.AfterClosing(afterClosing)

	// Blocks!
	e := rm.OpenPortal()

	if e != nil {
		log.Fatal(e)
	}
}

func createWebServer() http.Handler {
	webFiles, e := fs.Sub(ui.BuildFS, "build")
	if e != nil {
		panic(e)
	}
	return http.FileServerFS(webFiles)
}
