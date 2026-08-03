package videostreamer

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type TokenPool interface {
	Validate(token string) bool
}

type VideoStreamer struct {
	tp TokenPool
}

// TODO: Update to accept an interface with a function that
//
//	accepts an ID and returns a local filepath.
func New(tp TokenPool) http.Handler {
	return VideoStreamer{tp}
}

func (vs VideoStreamer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Update to accept an ID to a video that is looked
	//       up in the project data to get the path.
	path := getAndCheckParams(w, r, vs.tp)
	serveVideo(w, r, path)
}

func getAndCheckParams(w http.ResponseWriter, r *http.Request, tp TokenPool) string {
	token := r.URL.Query().Get("token")
	if isMissingOrInvalidToken(token, tp) {
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

func isMissingOrInvalidToken(token string, tokenPool TokenPool) bool {
	if token == "" || !tokenPool.Validate(token) {
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
