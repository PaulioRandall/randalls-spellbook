package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/crgimenes/glaze"
)

//go:embed build/*
var uiFS embed.FS

// Run starts (blocking) the UI by creating a WebView
// window and starting a file server for the build dir
// (output of Svelte's static adpater).
func Run(debug bool) error {
	buildDirFileServer, e := newBuildDirFileServer()

	if e != nil {
		return e
	}

	return startUi(buildDirFileServer, debug)
}

func newBuildDirFileServer() (http.Handler, error) {
	buildFS, e := fs.Sub(uiFS, "build")
	if e != nil {
		return nil, e
	}

	server := http.FileServer(http.FS(buildFS))
	return server, nil
}

func startUi(handler http.Handler, debug bool) error {
	options := glaze.AppOptions{
		Debug:   debug,
		Title:   "Randall's Spellbook",
		Width:   800,
		Height:  600,
		Hint:    glaze.HintNone,
		Handler: handler,
		OnReady: func(addr string) {
			println("WebView started and ready!")
		},
	}

	return glaze.AppWindow(options)
}
