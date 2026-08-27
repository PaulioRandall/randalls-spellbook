package ui

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"

	"github.com/crgimenes/glaze"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/mediaserver"
	"github.com/PaulioRandall/randalls-spellbook/pkg/project"
)

//go:embed build/*
var BuildFS embed.FS

// TODO: Spell to open database that adds a closer to the
//       sourcerer.
// TODO: Spell GetMediaById
// TODO: Spell ListMedia
// TODO: Spell AddMedia
// TODO: Select SelectLocalFile

// project stores project data.
var proj *project.Project

// webview is set when app is run.
var webview glaze.WebView

// Run starts (blocking) the UI by creating a WebView
// window and starting the file server.
func Run(debug bool) error {
	proj = project.New("./testproject/data.sqlite")

	// TEMP START
	// TODO: Implement project selection in UI.
	e := proj.Open()
	if e != nil {
		return e
	}
	// TEMP END

	fileServer, e := createFileServer()
	if e != nil {
		return e
	}

	return startUi(fileServer, debug)
}

func createFileServer() (*http.ServeMux, error) {
	svelteFS, e := fs.Sub(BuildFS, "build")
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

func startUi(handler http.Handler, debug bool) (e error) {
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

			e = w.Bind("selectLocalMediaFile", selectLocalMediaFile)
			if e != nil {
				return e
			}

			e = w.Bind("addMediaToProject", addMediaToProject)
			if e != nil {
				return e
			}

			e = w.Bind("getAllMedia", getAllMedia)
			if e != nil {
				return e
			}

			e = w.Bind("getMediaById", getMediaById)
			if e != nil {
				return e
			}

			return nil
		},
	}

	defer func() {
		closeErr := proj.Close()

		// Don't hide original error.
		if e == nil {
			log.Println(closeErr)
			e = closeErr
		}
	}()

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
	m, e := proj.AddMedia(
		mediaType,
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

func makeMediaResult(media data.Media) MediaResult {
	return MediaResult{
		EntityId:    media.EntityId.String(),
		MediaType:   media.MediaType,
		Name:        media.Name,
		Description: media.Description,
		LocalPath:   media.LocalPath,
	}
}

func getAllMedia() ([]MediaResult, error) {
	mediaList, e := proj.ListMedia()

	if e != nil {
		return nil, e
	}

	result := []MediaResult{}

	for _, media := range mediaList {
		result = append(result, makeMediaResult(media))
	}

	return result, nil
}

func getMediaById(entityId string) (MediaResult, error) {
	empty := MediaResult{}
	media, e := proj.GetMediaById(data.EntityId(entityId))

	if e != nil {
		return empty, e
	}

	if media == (data.Media{}) {
		return empty, errors.New("Unable to find media")
	}

	return makeMediaResult(media), nil
}
