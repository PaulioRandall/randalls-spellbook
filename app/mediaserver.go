package app

import (
	"log"
	"net/http"
	"os"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
)

type MediaServer struct {
	w *sourcery.World
}

type Res = http.ResponseWriter
type PtrMS = *MediaServer

func NewMediaServer(w *sourcery.World) PtrMS {
	return &MediaServer{}
}

func (ms PtrMS) Open(w *sourcery.World) {
	ms.w = w
}

func (ms PtrMS) Close() {
	ms.w = nil
}

func (ms PtrMS) ServeHTTP(w Res, r *http.Request) {
	entityId := r.URL.Query().Get("entity_id")
	if entityId == "" {
		httpErrBadIdParam(w)
		return
	}

	mediaUntyped, e := ms.w.Go("GetMediaById", entityId)
	if e != nil {
		httpErrMediaLookup(w, e)
		return
	}

	media, ok := mediaUntyped.(Media)
	if !ok {
		httpErrNotMediaObject(w)
		return
	}

	if media == (Media{}) {
		httpErrMediaNotFound(w)
		return
	}

	file, e := os.Open(media.LocalPath)
	if e != nil {
		// TODO: Create different response based on Not Found
		//       error and access error.
		httpErrMediaFileNotFound(w, e)
		return
	}

	defer file.Close()

	info, e := file.Stat()
	if e != nil {
		httpErrMediaFileAccess(w, e)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}

func httpErrBadIdParam(w Res) {
	http.Error(
		w,
		"Missing or invalid entity ID parameter",
		http.StatusBadRequest,
	)
}

func httpErrMediaLookup(w Res, e error) {
	log.Println(e)
	http.Error(
		w,
		"Error looking up media",
		http.StatusInternalServerError,
	)
}

func httpErrNotMediaObject(w Res) {
	http.Error(
		w,
		`Sanity check failed! Expected Media{} to be retuned by 'Go("GetMediaById", entityId)'`,
		http.StatusInternalServerError,
	)
}

func httpErrMediaNotFound(w Res) {
	http.Error(
		w,
		"Could not find media by ID",
		http.StatusNotFound,
	)
}

func httpErrMediaFileNotFound(w Res, e error) {
	// TODO: Create different response based on Not Found
	//       error and access error.
	log.Println(e)
	http.Error(
		w,
		"Media in database, but could not find media file",
		http.StatusNotFound,
	)
}

func httpErrMediaFileAccess(w Res, e error) {
	log.Println(e)
	http.Error(
		w,
		"Could not read media file stats",
		http.StatusInternalServerError,
	)
}

var _ sourcery.HttpPortal = &MediaServer{}
