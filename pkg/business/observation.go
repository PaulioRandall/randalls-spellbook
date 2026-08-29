package business

import (
	"encoding/json"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

func JsonToObservation(_ any, input any) spellbook.Effect {
	str := spellbook.Demystify[string](input)
	ob := data.Observation{}
	e := json.Unmarshal([]byte(str), &ob)
	return spellbook.Choose(ob, e)
}

func AddObservation(ctx any, input any) spellbook.Effect {
	inventory := getInventory(ctx)
	ob := spellbook.Demystifyf[data.Observation](
		input,
		"Unable to demystify JSON into an Observation",
	)
	e := inventory.InsertObservation(ob)
	return spellbook.Choose(ob, e)
}

func ListObservationsByMediaId(
	ctx any,
	input any,
) spellbook.Effect {
	inventory := getInventory(ctx)
	mediaId := spellbook.Demystify[string](input)
	obs, e := inventory.ListObservationsByMediaId(mediaId)
	return spellbook.Choose(obs, e)
}
