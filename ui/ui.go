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
	fileServer, e := createFileServer()

	if e != nil {
		return e
	}

	return startUi(fileServer, debug)
}

func createFileServer() (*http.ServeMux, error) {
	svelteFS, e := fs.Sub(buildFS, "build")
	if e != nil {
		return nil, e
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

			err = w.Bind("addMediaToProject", addMediaToProject)
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

func addMediaToProject(
	mediaType string,
	name string,
	description string,
	localPath string,
) (string, error) {
	mt, e := entity.ToMediaType(mediaType)
	if e != nil {
		return "", e
	}

	m, e := proj.AddMedia(
		mt,
		name,
		description,
		localPath,
	)
	return string(m.EntityId), e
}

type MediaResult struct {
	EntityId    string `json:"entityId"`
	MediaType   string `json:"mediaType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LocalPath   string `json:"lcoalPath"`
}

func makeMediaResult(media entity.Media) MediaResult {
	return MediaResult{
		EntityId:    media.EntityId.String(),
		MediaType:   media.MediaType.String(),
		Name:        media.Name,
		Description: media.Description,
		LocalPath:   media.LocalPath,
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

	if media.IsEmpty() {
		return MediaResult{}, errors.New("Unable to find media")
	}

	return makeMediaResult(media), nil
}
