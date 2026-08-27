package business

import (
	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

func SelectLocalFile(ctx any, _ any) spellbook.Effect {
	rm := getRealm(ctx)

	// Blocks!
	path, e := rm.WebView().OpenFile(glaze.FileDialogOptions{
		Title: "Select media file",
	})

	return spellbook.Choose(path, e)
}

func ListMedia(ctx any, _ any) spellbook.Effect {
	inventory := getInventory(ctx)
	mediaList, e := inventory.ListMedia()
	return spellbook.Choose(mediaList, e)
}

func GetMediaById(ctx any, input any) spellbook.Effect {
	inventory := getInventory(ctx)

	entityIdStr := spellbook.Demystify[string](input)
	entityId := data.EntityId(entityIdStr)
	media, e := inventory.GetMediaById(entityId)

	if e != nil {
		return spellbook.Cursed[data.Media](e)
	}

	if media == (data.Media{}) {
		return spellbook.Curse[data.Media]("Unable to find media")
	}

	return spellbook.Bless(media)
}

func AddMedia(ctx any, input any) spellbook.Effect {
	inventory := getInventory(ctx)

	genMedia := func() *data.Media {
		return &data.Media{}
	}
	mediaDemystifyer := spellbook.JsonDemystifyer(genMedia)

	ef := mediaDemystifyer(nil, input)
	if ef.Cursed() {
		return ef
	}

	m, _ := spellbook.DemystifyEffect[*data.Media](ef)
	e := inventory.InsertMedia(*m)
	return spellbook.Choose(m, e)
}
