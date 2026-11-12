package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	_ "github.com/PaulioRandall/randalls-spellbook/business"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
)

//go:embed ui/build/*
var webFiles embed.FS

func main() {
	uiFiles, e := fs.Sub(webFiles, "ui/build")
	if e != nil {
		log.Fatal(e)
	}

	sourcery.Append(
		sourcery.OptionDebug(true),
		sourcery.OptionTitle("Randall's Spellbook"),
		sourcery.OptionSize(1000, 600),
		sourcery.OptionServer("/", http.FileServerFS(uiFiles)),
	)

	app := sourcery.DefaultApp()

	e = app.Run()
	if e != nil {
		log.Fatal(e)
	}
}
