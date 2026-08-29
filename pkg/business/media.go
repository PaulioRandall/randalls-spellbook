package business

import (
	"encoding/json"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

func ListMedia(ctx any, _ any) spellbook.Effect {
	inventory := getInventory(ctx)
	mediaList, e := inventory.ListMedia()
	return spellbook.Choose(mediaList, e)
}

func GetMediaById(ctx any, input any) spellbook.Effect {
	inventory := getInventory(ctx)
	entityId := spellbook.Demystify[string](input)

	media, e := inventory.GetMediaById(entityId)
	if e != nil {
		return spellbook.Cursed[data.Media](e)
	}

	if media == (data.Media{}) {
		return spellbook.Curse[data.Media]("Unable to find media")
	}

	return spellbook.Bless(media)
}

func JsonToMedia(_ any, input any) spellbook.Effect {
	str := spellbook.Demystifyf[string](
		input,
		"Unable to demystify JSON into a Media",
	)
	media := data.Media{}
	e := json.Unmarshal([]byte(str), &media)
	return spellbook.Choose(media, e)
}

func AddMedia(ctx any, input any) spellbook.Effect {
	inventory := getInventory(ctx)
	media := spellbook.Demystify[data.Media](input)
	insertedMedia, e := inventory.InsertMedia(media)
	return spellbook.Choose(insertedMedia, e)
}
