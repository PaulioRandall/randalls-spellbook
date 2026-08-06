package ui

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"

	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/entity"
	"github.com/PaulioRandall/randalls-spellbook/pkg/mediaserver"
	"github.com/PaulioRandall/randalls-spellbook/pkg/project"
)

//go:embed build/*
var buildFS embed.FS

// project stores project data.
var proj *project.Project

// webview is set when app is run.
var webview glaze.WebView

// Run starts (blocking) the UI by creating a WebView
// window and starting the file server.
func Run(debug bool) error {
	proj = project.New()
	fileServer, err := createFileServer()

	if err != nil {
		return err
	}

	return startUi(fileServer, debug)
}

func createFileServer() (*http.ServeMux, error) {
	svelteFS, err := fs.Sub(buildFS, "build")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	// Handle requests for media, e.g. videos.
	mux.Handle("/media/", mediaserver.New(proj))

	// Handle requests for statically built Svelte files.
	mux.Handle("/", http.FileServerFS(svelteFS))

	return mux, nil
}

func startUi(handler http.Handler, debug bool) error {
	options := AppOptions{
		Debug:   debug,
		Title:   "Randall's Spellbook",
		Width:   800,
		Height:  600,
		Hint:    glaze.HintNone,
		Handler: handler,
		OnWebViewReady: func(w glaze.WebView) error {
			webview = w

			// TODO: Create struct and use BindMethods.

			err := w.Bind("selectLocalMediaFile", selectLocalMediaFile)
			if err != nil {
				return err
			}

			err = w.Bind("addVideoToProject", addVideoToProject)
			if err != nil {
				return err
			}

			err = w.Bind("getAllMedia", getAllMedia)
			if err != nil {
				return err
			}

			err = w.Bind("getMediaById", getMediaById)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return AppWindow(options)
}

func selectLocalMediaFile() (string, error) {
	return webview.OpenFile(glaze.FileDialogOptions{
		Title: "Select media file",
	})
}

func addVideoToProject(
	name, description, localPath string,
) (string, error) {
	media, err := proj.AddVideo(name, description, localPath)
	return string(media.EntityId()), err
}

type MediaResult struct {
	EntityId    string `json:"entityId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LocalPath   string `json:"lcoalPath"`
}

func makeMediaResult(media entity.Media) MediaResult {
	return MediaResult{
		EntityId:    media.EntityId().String(),
		Name:        media.Name(),
		Description: media.Description(),
		LocalPath:   media.LocalPath(),
	}
}

func getAllMedia() []MediaResult {
	mediaList := proj.GetAllMedia()
	result := []MediaResult{}

	for _, media := range mediaList {
		result = append(result, makeMediaResult(media))
	}

	return result
}

func getMediaById(entityId string) (MediaResult, error) {
	media := proj.GetMediaById(entity.EntityId(entityId))

	if media == nil {
		return MediaResult{}, errors.New("Unable to find media")
	}

	return makeMediaResult(media), nil
}
