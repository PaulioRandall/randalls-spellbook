package sourcery

import (
	"fmt"
	"strings"

	"github.com/crgimenes/glaze"
)

type World struct {
	options   AppOptions
	stones    map[string]RuneStone
	spellbook map[string]Spell
	webview   glaze.WebView
}

func buildWorld(
	options AppOptions,
	stones map[string]RuneStone,
	spellbook map[string]Spell,
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

func (w *World) Enter() error {
	w.options.OnWebViewReady = w.onWebViewReady
	e := AppWindow(w.options)

	w.Log("Freeing rune stones")
	for name, rs := range w.stones {
		w.Log("\t%s{}", name)
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

	w.Log("Inscribing spells:")
	for name, _ := range w.spellbook {
		w.Log("\t%s", name)
	}

	w.Log("Binding rune stones:")
	for name, rs := range w.stones {
		w.Log("\t%s{}", name)
		rs.Bind(w)
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
