package sourcery

import (
	"fmt"
	"strings"

	"github.com/crgimenes/glaze"
)

type World struct {
	options   AppOptions
	portalMap PortalMap
	funcMap   FuncMap
	webview   glaze.WebView
}

func buildWorld(
	options AppOptions,
	portalMap PortalMap,
	funcMap FuncMap,
) *World {
	return &World{
		options:   options,
		portalMap: portalMap,
		funcMap:   funcMap,
	}
}

func (w *World) WebView() glaze.WebView {
	return w.webview
}

func (w *World) Enter() error {
	w.options.OnWebViewReady = w.onWebViewReady
	e := AppWindow(w.options)

	w.Log("Freeing portals:")
	for name, port := range w.portalMap {
		w.Log("\t%s{}", name)
		port.Free()
	}

	return e
}

func (w *World) Exit() {
	w.WebView().Terminate()
}

func (w *World) onWebViewReady(wv glaze.WebView) error {
	w.webview = wv

	wv.Init(`
		var Go = function(funcName, ...args) {
			return GoRaw(
				funcName,
				JSON.stringify(args),
			)
		}
	`)

	e := wv.Bind("GoRaw", w.GoRaw)
	if e != nil {
		return e
	}

	w.Log("Functions registered:")
	for name, _ := range w.funcMap {
		w.Log("\t%s", name)
	}

	w.Log("Opening portals:")
	for name, port := range w.portalMap {
		w.Log("\t%s{}", name)
		port.Init(w)
	}

	return nil
}

func (w *World) GoRaw(
	funcName string,
	jsonArgs string,
) (any, error) {
	if f, ok := w.funcMap[funcName]; ok {
		w.Log("Go: %s", funcName)
		thunk, e := WithJsonArgs(f, jsonArgs)
		if e != nil {
			return nil, e
		}
		return thunk.Call()
	}

	w.Log("Unknown function: %s", funcName)
	return nil, fmt.Errorf("Unknown function: %s", funcName)
}

func (w *World) Log(msg string, args ...any) {
	if !w.options.Debug {
		return
	}

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Printf("%s\n", prefixLines(msg, "[Sourcery] "))
}

func prefixLines(s, pre string) string {
	lines := strings.Split(s, "\n")
	for i, _ := range lines {
		lines[i] = pre + lines[i]
	}
	return strings.Join(lines, "\n")
}
