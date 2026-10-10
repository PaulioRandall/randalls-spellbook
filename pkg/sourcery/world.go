package sourcery

import (
	"fmt"

	"github.com/crgimenes/glaze"
)

type World struct {
	options Options
}

func newWorld(options Options) *World {
	return &World{
		options: options,
	}
}

func (w *World) Start() (e error) {
	var baseUrl string

	if w.options.serveMux != nil {
		cs, e := createContentServer(w.options.serveMux)
		if e != nil {
			return e
		}

		go func() {
			err := cs.serve()
			if e == nil {
				e = err
			}
		}()

		defer func() {
			err := cs.close()
			if e == nil {
				e = err
			}
		}()

		baseUrl = cs.baseUrl
	}

	defer func() {
		w.options.webview.Destroy()
		w.options.webview = nil
	}()

	if baseUrl != "" {
		w.options.webview.Navigate(baseUrl)
	}

	e = w.bindGoRaw()
	if e != nil {
		return e
	}

	defer w.unitialise()
	w.initialise()

	// Starts the webview and blocks until it exits.
	w.options.webview.Run()

	return nil
}

func (w *World) initialise() {
	for _, v := range w.options.initialisers {
		v.Init(w)
	}
}

func (w *World) unitialise() {
	for _, v := range w.options.unitialisers {
		v.Unit()
	}
}

func (w *World) Options() Options {
	return w.options
}

func (w *World) WebView() glaze.WebView {
	return w.options.webview
}

func (w *World) Exit() {
	w.WebView().Terminate()
}

func (w *World) bindGoRaw() error {
	e := w.options.webview.Bind("GoRaw", w.GoRaw)
	if e != nil {
		return e
	}

	w.options.webview.Init(`
		var Go = function(funcName, ...args) {
			return GoRaw(
				funcName,
				JSON.stringify(args),
			)
		}
	`)

	return nil
}

func (w *World) GoRaw(
	funcName string,
	jsonArgs string,
) (any, error) {
	if f, ok := w.options.functions[funcName]; ok {
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
	if !w.options.debug {
		return
	}

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Printf("[Sourcery] %s\n", msg)
}
