package main

import (
	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

type MediaResult struct {
	EntityId    string `json:"entityId"`
	MediaType   string `json:"mediaType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LocalPath   string `json:"lcoalPath"`
}

func makeMediaResult(media data.Media) MediaResult {
	return MediaResult{
		EntityId:    media.EntityId.String(),
		MediaType:   media.MediaType,
		Name:        media.Name,
		Description: media.Description,
		LocalPath:   media.LocalPath,
	}
}

var store data.Store

func afterOpening(rm *sourcery.Realm) error {
	store = data.NewStore("./testproject/data.sqlite")

	e := store.Open()
	if e != nil {
		return e
	}

	rm.Enscribe("SelectLocalFile", SelectLocalFile)
	rm.Enscribe("ListMedia", ListMedia)
	rm.Enscribe("GetMediaById", GetMediaById)
	rm.Enscribe("AddMedia", AddMedia)

	return nil
}

func afterClosing(rm *sourcery.Realm) error {
	var e error

	if store != nil {
		e = store.Close()
		store = nil
	}

	return e
}

/*
	// InsertMedia inserts a media. The media is assumed to
	// be valid.
	InsertMedia(Media) error

	// ListMedia returns all the media entities.
	ListMedia() ([]Media, error)

	// GetMediaById returns the media entity with the given
	// EntityID.
	GetMediaById(EntityId) (Media, error)
*/

func SelectLocalFile(
	ctxAny any,
	input spellbook.Effect,
) spellbook.Effect {
	rm, ok := ctxAny.(*sourcery.Realm)
	if !ok {
		return spellbook.Curse("Wrong context, expected a *Realm")
	}

	// Blocks!
	path, e := rm.WebView().OpenFile(glaze.FileDialogOptions{
		Title: "Select media file",
	})

	return spellbook.Choose(path, e)
}

func GetProject(
	_ any,
	_ spellbook.Effect,
) spellbook.Effect {
	return spellbook.Choose(store.GetProject())
}

func ListMedia(
	_ any,
	input spellbook.Effect,
) spellbook.Effect {
	mediaList, e := store.ListMedia()

	if e != nil {
		return spellbook.Cursed(e)
	}

	result := []MediaResult{}

	for _, media := range mediaList {
		result = append(result, makeMediaResult(media))
	}

	return spellbook.Bless(result)
}

func GetMediaById(
	_ any,
	input spellbook.Effect,
) spellbook.Effect {
	entityId, ok := input.Result().(string)
	if !ok {
		return spellbook.Curse("Wrong type, expected a string")
	}

	mediaId := data.EntityId(entityId)
	media, e := store.GetMediaById(mediaId)

	if e != nil {
		return spellbook.Cursed(e)
	}

	if media == (data.Media{}) {
		return spellbook.Curse("Unable to find media")
	}

	result := makeMediaResult(media)
	return spellbook.Bless(result)
}

func AddMedia(
	_ any,
	input spellbook.Effect,
) spellbook.Effect {
	type protoMedia struct {
		MediaType   string `json:"mediaType"`
		Name        string `json:"name"`
		Description string `json:"description"`
		LocalPath   string `json:"lcoalPath"`
	}

	var p protoMedia
	protoDemystifyer := spellbook.JsonDemystifyer(p)
	protoEffect := protoDemystifyer(nil, input)

	if protoEffect.Cursed() {
		return protoEffect
	}

	p, _ = spellbook.TransmuteEffect[protoMedia](protoEffect)

	m, e := data.MakeMedia(
		p.MediaType,
		p.Name,
		p.Description,
		p.LocalPath,
	)

	if e != nil {
		return spellbook.Cursed(e)
	}

	e = store.InsertMedia(m)
	if e != nil {
		return spellbook.Cursed(e)
	}

	return spellbook.Bless(m)
}
