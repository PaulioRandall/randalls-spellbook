package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
	"github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

type MediaServer struct {
	realm *sourcery.Realm
}

func (ms MediaServer) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	media := ms.lookupMedia(w, r)
	if media != (data.Media{}) {
		serveMedia(w, r, media)
	}
}

func (ms MediaServer) lookupMedia(
	w http.ResponseWriter,
	r *http.Request,
) data.Media {
	entityIdParam := r.URL.Query().Get("entity_id")

	if entityIdParam == "" {
		http.Error(
			w,
			"Missing or invalid entity ID parameter",
			http.StatusBadRequest,
		)

		return data.Media{}
	}

	entityId := data.EntityId(entityIdParam)
	effect := ms.realm.Cast("GetMediaById", ms.realm, entityId)
	media, e := spellbook.TransmuteEffect[data.Media](effect)

	if e != nil {
		log.Println(e)
		http.Error(
			w,
			"Error looking up media",
			http.StatusInternalServerError,
		)
	}

	if media == (data.Media{}) {
		http.Error(
			w,
			"Could not find media by ID",
			http.StatusBadRequest,
		)
	}

	return media
}

func serveMedia(
	w http.ResponseWriter,
	r *http.Request,
	media data.Media,
) {
	file, err := os.Open(media.LocalPath)
	if err != nil {
		http.Error(
			w,
			"Could not find media file",
			http.StatusNotFound,
		)
		return
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(
			w,
			"Could not read media file stats",
			http.StatusInternalServerError,
		)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}
