package view

import (
	"github.com/crgimenes/glaze"
)

type View struct {
	webView glaze.WebView
}

func New(debug bool) (*View, error) {
	wv, e := glaze.New(debug)

	if e != nil {
		return nil, e
	}

	addContentToWebView(wv)
	v := View{webView: wv}
	return &v, nil
}

func (v *View) Alive() bool {
	return v.webView != nil
}

func (v *View) Run() {
	if v.Alive() {
		v.webView.Run()
	}
}

func (v *View) Destroy() {
	if v.Alive() {
		v.webView.Destroy()
		v.webView = nil
	}
}
