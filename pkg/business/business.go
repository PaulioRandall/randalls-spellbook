package business

import (
	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

type Realm = *sourcery.Realm[data.Store]

func getRealm(ctx any) Realm {
	return spellbook.Demystify[Realm](ctx)
}

func getInventory(ctx any) data.Store {
	return getRealm(ctx).Inventory
}

func SelectLocalFile(ctx any, _ any) spellbook.Effect {
	rm := getRealm(ctx)

	// Blocks!
	path, e := rm.WebView().OpenFile(glaze.FileDialogOptions{
		Title: "Select media file",
	})

	return spellbook.Choose(path, e)
}
