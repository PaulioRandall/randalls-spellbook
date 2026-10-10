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
		New(true).
		SetTitle("Randall's Spellbook").
		SetSize(800, 600).
		AddEntity(&app.Util{}).
		AddEntityServer("/media/", &app.Datastore{}).
		AddServer("/", http.FileServerFS(uiFiles)).
		CreateWorld().
		Start() // Blocks until WebView closes.

	if e != nil {
		log.Fatal(e)
	}
}
