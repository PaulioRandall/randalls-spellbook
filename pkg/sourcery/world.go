package sourcery

import (
	"fmt"

	"github.com/crgimenes/glaze"
)

type World struct {
	options   AppOptions
	stones    Litholog
	spellbook Codex
	webview   glaze.WebView
}

func buildWorld(
	options AppOptions,
	stones Litholog,
	spellbook Codex,
) *World {
	return &World{
		options:   options,
		stones:    stones,
		spellbook: spellbook,
	}
}

func (w *World) WebView() glaze.WebView {
	return w.webview
}

func (w *World) SeekSpell(spellName string) Spell {
	return w.spellbook[spellName]
}

func (w *World) OpenPortal() error {
	w.options.OnWebViewReady = w.onWebViewReady
	e := AppWindow(w.options)

	for name, rs := range w.stones {
		w.Log("Freeing rune stone: %s", name)
		rs.Free()
	}

	return e
}

func (w *World) onWebViewReady(wv glaze.WebView) error {
	w.webview = wv

	e := wv.Bind("Go", w.Go)
	if e != nil {
		return e
	}

	for name, rs := range w.stones {
		w.Log("Binding rune stone: %s", name)
		rs.Bind(w)
	}

	return nil
}

func (w *World) Go(
	spellName string,
	args ...any,
) (any, error) {
	if spell, ok := w.spellbook[spellName]; ok {
		w.Log("Invoking spell: %s", spellName)
		return spell.Invoke(args...)
	}

	w.Log("Unknown spell '%s'", spellName)
	return nil, fmt.Errorf("Unknown spell '%s'", spellName)
}

func (w *World) Log(msg string, args ...any) {
	if !w.options.Debug {
		return
	}

	if len(args) == 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	fmt.Print("[Sourcery] ", msg)
}
