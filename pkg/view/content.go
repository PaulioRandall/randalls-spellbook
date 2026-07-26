package view

import (
	"github.com/crgimenes/glaze"
)

func addContentToWebView(w glaze.WebView) {
	w.SetTitle("Glaze")
	w.SetSize(800, 600, glaze.HintNone)
	w.SetHtml("<h1>Hello from Glaze</h1>")
}
