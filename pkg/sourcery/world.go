package sourcery

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/crgimenes/glaze"
)

type contentServer struct {
	server   *http.Server
	listener net.Listener
	baseUrl  string
}

func createContentServer(handler http.Handler) (contentServer, error) {
	server := &http.Server{
		// TODO: Needs optimising.
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       20 * time.Second,
		MaxHeaderBytes:    16 << 10, // 16 KiB
	}

	// Create listener on a random loopback port.
	listener, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return contentServer{}, e
	}

	baseUrl := "http://" + listener.Addr().(*net.TCPAddr).String()

	cs := contentServer{
		server:   server,
		listener: listener,
		baseUrl:  baseUrl,
	}

	return cs, nil
}

func (cs *contentServer) serve() error {
	return cs.server.Serve(cs.listener)
}

func (cs *contentServer) close() error {
	return cs.server.Close()
}

type World struct {
	debug          bool
	title          string
	width          int
	height         int
	handler        *http.ServeMux
	functions      map[string]any
	initialisers   []Initable
	uninitialisers []func()
	webview        glaze.WebView
}

func (w *World) run() (e error) {
	cs, e := createContentServer(w.handler)
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

	defer w.callUninitialisers()

	// Create the webview window.
	w.webview, e = glaze.New(w.debug)
	if e != nil {
		return e
	}

	defer func() {
		w.webview.Destroy()
		w.webview = nil
	}()

	w.webview.SetTitle(w.title)
	w.webview.SetSize(w.width, w.height, glaze.HintNone)
	w.webview.Navigate(cs.baseUrl)
	w.bindGoRaw()
	w.callInitialisers()

	// Starts the webview and blocks until it exits.
	w.webview.Run()

	return nil
}

func (w *World) bindGoRaw() error {
	w.webview.Init(`
		var Go = function(funcName, ...args) {
			return GoRaw(
				funcName,
				JSON.stringify(args),
			)
		}
	`)
	return w.webview.Bind("GoRaw", w.GoRaw)
}

func (w *World) callInitialisers() {
	for _, initable := range w.initialisers {
		if unit := initable.Init(w); unit != nil {
			w.uninitialisers = append(w.uninitialisers, unit)
		}
	}
}

func (w *World) callUninitialisers() {
	for _, unit := range w.uninitialisers {
		unit()
	}
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
	if !w.debug {
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
