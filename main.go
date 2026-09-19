package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/PaulioRandall/randalls-spellbook/app"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
)

//go:embed ui/build/*
var webFiles embed.FS

func main() {
	uiFiles, e := fs.Sub(webFiles, "ui/build")
	if e != nil {
		log.Fatal(e)
	}

	appPortal := &app.App{}
	dsPortal := &app.Datastore{}

	e = sourcery.NewCreator().
		Debug().
		Name("Randall's Spellbook").
		Size(800, 600).
		AddPortal(appPortal).
		AddPortal(dsPortal).
		AddServer("/media/", dsPortal).
		AddServer("/", http.FileServerFS(uiFiles)).
		BuildWorld().
		Enter() // Blocks until WebView closes.

	if e != nil {
		log.Fatal(e)
	}
}
