package ui

import (
	"github.com/crgimenes/glaze"
)

func Run(debug bool) error {
	webview, e := glaze.New(debug)

	if e != nil {
		return e
	}

	defer webview.Destroy()

	e = addContentToWebView(webview)
	if e != nil {
		return e
	}

	webview.Run()
	return nil
}
