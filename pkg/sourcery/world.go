package sourcery

import (
	"fmt"
	"strings"

	"github.com/crgimenes/glaze"
)

type World struct {
	options        AppOptions
	functions      map[string]any
	initialisers   []Initable
	uninitialisers []func()
	webview        glaze.WebView
}

func (w *World) run() error {
	w.options.OnWebViewReady = w.onWebViewReady
	e := AppWindow(w.options)

	for _, unit := range w.uninitialisers {
		unit()
	}

	return e
}

func (w *World) WebView() glaze.WebView {
	return w.webview
}

func (w *World) Exit() {
	w.WebView().Terminate()
}

func (w *World) GoRaw(
	funcName string,
	jsonArgs string,
) (any, error) {
	if f, ok := w.functions[funcName]; ok {
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

	for _, initable := range w.initialisers {
		if unit := initable.Init(w); unit != nil {
			w.uninitialisers = append(w.uninitialisers, unit)
		}
	}

	return nil
}

func prefixLines(s, pre string) string {
	lines := strings.Split(s, "\n")
	for i, _ := range lines {
		lines[i] = pre + lines[i]
	}
	return strings.Join(lines, "\n")
}
