package ui

import (
	"embed"

	"github.com/crgimenes/glaze"
)

//go:embed testdata/*.js
var ui embed.FS

func addContentToWebView(w glaze.WebView) error {

	initJs, e := ui.ReadFile("testdata/init.js")

	if e != nil {
		return e
	}

	w.SetTitle("Glaze")
	w.SetSize(800, 600, glaze.HintNone)
	w.SetHtml(`<h1 id="heading">Hello from Glaze</h1>`)
	w.Init(string(initJs))

	return nil
}
