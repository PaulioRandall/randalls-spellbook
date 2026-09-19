package sourcery

import (
	"fmt"
	"strings"

	"github.com/crgimenes/glaze"
)

type World struct {
	options   AppOptions
	portals   PortalMap
	spellbook Spellbook
	webview   glaze.WebView
}

func buildWorld(
	options AppOptions,
	portals PortalMap,
	spellbook Spellbook,
) *World {
	return &World{
		options:   options,
		portals:   portals,
		spellbook: spellbook,
	}
}

func (w *World) WebView() glaze.WebView {
	return w.webview
}

func (w *World) Enter() error {
	w.options.OnWebViewReady = w.onWebViewReady
	e := AppWindow(w.options)

	w.Log("Closing portals:")
	for name, port := range w.portals {
		w.Log("\t%s{}", name)
		port.WorldExit()
	}

	return e
}

func (w *World) Exit() {
	w.WebView().Terminate()
}

func (w *World) onWebViewReady(wv glaze.WebView) error {
	w.webview = wv

	e := wv.Bind("Go", w.Go)
	if e != nil {
		return e
	}

	w.Log("Inscribing spells:")
	for name, _ := range w.spellbook {
		w.Log("\t%s", name)
	}

	w.Log("Opening portals:")
	for name, port := range w.portals {
		w.Log("\t%s{}", name)
		port.WorldEnter(w)
	}

	return nil
}

func (w *World) Go(
	spellName string,
	args ...any,
) (any, error) {
	if spell, ok := w.spellbook[spellName]; ok {
		w.Log("Invoking: %s", spellName)
		return spell.Invoke(args...)
	}

	w.Log("Unknown spell: %s", spellName)
	return nil, fmt.Errorf("Unknown spell: %s", spellName)
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
