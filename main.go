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

	e = sourcery.
		NewWorld(true).
		SetTitle("Randall's Spellbook").
		SetSize(800, 600).
		Entity(&app.Util{}).
		Server("/", http.FileServerFS(uiFiles)).
		EntityServer("/media/", &app.Datastore{}).
		Start() // Blocks until WebView closes.

	if e != nil {
		log.Fatal(e)
	}
}
