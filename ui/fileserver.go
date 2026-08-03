package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/PaulioRandall/randalls-spellbook/pkg/application"
	"github.com/PaulioRandall/randalls-spellbook/pkg/videostreamer"
)

//go:embed build/*
var buildFS embed.FS

func newFileServer() (*http.ServeMux, error) {
	svelteFS, e := fs.Sub(buildFS, "build")
	if e != nil {
		return nil, e
	}

	app := application.New()
	mux := http.NewServeMux()

	// Handle requests for media, e.g. videos.
	mux.Handle("/media/", videostreamer.New(app))

	// Handle requests for statically built Svelte files.
	mux.Handle("/", http.FileServerFS(svelteFS))

	return mux, nil
}
