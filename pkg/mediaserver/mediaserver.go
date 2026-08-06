package mediaserver

import (
	"net/http"
	"os"

	"github.com/PaulioRandall/randalls-spellbook/pkg/entity"
)

// TODO: Improve error descriptions to include names and
//       paths so users have the info to fix their own file
//       issues.
// TODO: Return 404 if media doesn't exist.

type MediaSource interface {
	GetMediaById(id entity.EntityId) entity.Media
}

type MediaServer struct {
	mediaSource MediaSource
}

func New(mediaSource MediaSource) http.Handler {
	return MediaServer{mediaSource}
}

func (ms MediaServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	media := lookupMedia(w, r, ms.mediaSource)
	if media != nil {
		serveMedia(w, r, media)
	}
}

func lookupMedia(w http.ResponseWriter, r *http.Request, mediaSource MediaSource) entity.Media {
	entityIdParam := r.URL.Query().Get("entity_id")

	if entityIdParam == "" {
		http.Error(w, "Missing or invalid entity ID parameter", http.StatusBadRequest)
		return nil
	}

	entityId := entity.EntityId(entityIdParam)
	media := mediaSource.GetMediaById(entityId)

	if media == nil {
		http.Error(w, "Could not find media by ID", http.StatusBadRequest)
		return nil
	}

	return media
}

func serveMedia(w http.ResponseWriter, r *http.Request, media entity.Media) {
	file, err := os.Open(media.LocalPath())
	if err != nil {
		http.Error(w, "Could not find media file", http.StatusNotFound)
		return
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not read media file stats", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}
