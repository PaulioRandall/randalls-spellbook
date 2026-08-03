package ui

import (
	"io/ioutil"
	"net/http"

	"github.com/crgimenes/glaze"
)

// webview is set when app is run.
var webview glaze.WebView

// Run starts (blocking) the UI by creating a WebView
// window and starting the file server.
func Run(debug bool) error {
	handler, e := newFileServer()

	if e != nil {
		return e
	}

	return startUi(handler, debug)
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
