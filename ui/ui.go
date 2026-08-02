package ui

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"net/http"

	"github.com/crgimenes/glaze"
)

//go:embed build/*
var uiFS embed.FS

// webview is set when app is run.
var webview glaze.WebView

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
	options := AppOptions{
		Debug:   debug,
		Title:   "Randall's Spellbook",
		Width:   800,
		Height:  600,
		Hint:    glaze.HintNone,
		Handler: handler,
		OnWebViewReady: func(w glaze.WebView) error {
			webview = w

			e := w.Bind("selectVideoFile", selectVideoFile)
			if e != nil {
				return e
			}

			return w.Bind("readVideoFile", readVideoFile)
		},
	}

	return AppWindow(options)
}

func selectVideoFile() (string, error) {
	return webview.OpenFile(glaze.FileDialogOptions{
		Title: "Select file to open",
		Filters: []glaze.FileFilter{
			glaze.FileFilter{
				Name: "Videos",
				Extensions: []string{
					"mp4",
				},
			},
		},
	})
}

func readVideoFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(filepath)
}
