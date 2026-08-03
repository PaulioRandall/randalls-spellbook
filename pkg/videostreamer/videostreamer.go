package videostreamer

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type App interface {
	GetMediaPath(id string) string

	// TODO: Remove
	ConsumeToken(token string) bool
}

type VideoStreamer struct {
	app App
}

func New(app App) http.Handler {
	return VideoStreamer{app}
}

func (vs VideoStreamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Update to accept an ID to a video that is looked
	//       up in the project data to get the path.
	path := vs.getAndCheckParams(w, r)
	serveVideo(w, r, path)
}

func (vs VideoStreamer) getAndCheckParams(w http.ResponseWriter, r *http.Request) string {
	token := r.URL.Query().Get("token")
	if isMissingOrInvalidToken(token, vs.app) {
		http.Error(w, "Missing or invalid token", http.StatusBadRequest)
		return ""
	}

	path := r.URL.Query().Get("path")
	if isMissingOrInvalidPath(path) {
		http.Error(w, "Missing or invalid path", http.StatusBadRequest)
		return ""
	}

	return filepath.Clean(path)
}

func isMissingOrInvalidToken(token string, app App) bool {
	if token == "" || !app.ConsumeToken(token) {
		return true
	}
	return false
}

func isMissingOrInvalidPath(path string) bool {
	if path == "" {
		return true
	}

	if _, e := url.Parse(path); e != nil {
		return true
	}

	return false
}

func serveVideo(w http.ResponseWriter, r *http.Request, path string) {
	f, e := os.Open(path)
	if e != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	defer f.Close()

	info, e := f.Stat()
	if e != nil {
		http.Error(w, "Cannot stat file", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, info.Name(), info.ModTime(), f)
}
