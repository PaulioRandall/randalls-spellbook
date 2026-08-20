package mediaserver

import (
	"log"
	"net/http"
	"os"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
)

// TODO: Improve error descriptions to include names and
//       paths so users have the info to fix their own file
//       issues.
// TODO: Return 404 if media doesn't exist.

type MediaSource interface {
	GetMediaById(id data.EntityId) (data.Media, error)
}

type MediaServer struct {
	mediaSource MediaSource
}

func New(mediaSource MediaSource) http.Handler {
	return MediaServer{mediaSource}
}

func (ms MediaServer) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	media := lookupMedia(w, r, ms.mediaSource)
	if media != (data.Media{}) {
		serveMedia(w, r, media)
	}
}

func lookupMedia(
	w http.ResponseWriter,
	r *http.Request,
	mediaSource MediaSource,
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
	media, e := mediaSource.GetMediaById(entityId)

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
