package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/PaulioRandall/randalls-spellbook/pkg/videostreamer"
)

//go:embed build/*
var buildFS embed.FS

// Temp
type tokenPool struct{}

func (tp tokenPool) Validate(token string) bool {
	// TODO: Replace with actual checker.
	return true
}

func newFileServer() (*http.ServeMux, error) {
	svelteFS, e := fs.Sub(buildFS, "build")
	if e != nil {
		return nil, e
	}

	mux := http.NewServeMux()

	// Handle requests for media, e.g. videos.
	mux.Handle("/media/", videostreamer.New(tokenPool{}))

	// Handle requests for statically built Svelte files.
	mux.Handle("/", http.FileServerFS(svelteFS))

	return mux, nil
}
