package app

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MediaTypeVideo is only used by video implementations
// of Media.
const (
	MediaTypeVideo string = "video"
)

// Media holds information about a user defined media
// entity. It maps directly to a single table in the
// database.
type Media struct {
	// EntityId of the video.
	//
	// It must be unique within the project, never empty, and
	// never change.
	EntityId string

	// MediaType is type of the media, e.g. video, audio,
	// PDF, etc.
	MediaType string

	// Name is the user defined readable and meaningful name
	// for human users and AI agents. This is not the
	// filename, localPath is the filename.
	//
	// It must never be empty and should should be trimmed
	// of whitespace.
	Name string

	// Description is the user defined detailed explanation
	// of the video for human users and AI agents.
	//
	// It should compliment the Name field but is also
	// intended for general notes. It may be may be empty and
	// should should be trimmed of whitespace.
	Description string

	// LocalPath is the file path to the video file within
	// the local file system.
	LocalPath string
}

// Clean returns a new Media entity from an exsiting
// one ensuring all properties are valid. All values are
// trimmed and localPath is cleaned before being checked
// and assigned to the new media. If the EntityId is empty
// then a new one is allocated.
//
// If checks pass then the media is returned else an empty
// Media and an error. An error will occur if the Name or
// LocalPath are empty, or if LocalPath is not a valid
// absolute filepath.
//
// The existence of the file or file type are not checked
// but this may change in the future.
func (m Media) Clean() (Media, error) {
	m.EntityId = strings.TrimSpace(m.EntityId)
	m.Name = strings.TrimSpace(m.Name)
	m.Description = strings.TrimSpace(m.Description)
	m.LocalPath = strings.TrimSpace(m.LocalPath)

	if m.EntityId == "" {
		m.EntityId = randomEntityId()
	}

	if m.Name == "" {
		return Media{}, errors.New("Name must not be empty")
	}

	if m.LocalPath == "" {
		return Media{}, errors.New("LocalPath must not be empty")
	}

	m.LocalPath = filepath.Clean(m.LocalPath)
	if !filepath.IsAbs(m.LocalPath) {
		return Media{}, errors.New("LocalPath must be absolute")
	}

	return m, nil
}

// GetEntityId returns the unique entity ID of the media.
func (m Media) GetEntityId() string {
	return m.EntityId
}

// GetMediaType returns the type of the media, e.g. video,
// audio, PDF, etc.
func (m Media) GetMediaType() string {
	return m.MediaType
}

// GetName returns the user defined readable and meaningful
// name for humans users and AI agents. This is not the
// filename, LocalPath returns the filename. It should
// never be empty and should be trimmed of whitespace.
func (m Media) GetName() string {
	return m.Name
}

// GetDescription is the user defined detailed explanation
// of the video for human users and AI agents It
// compliments the media name but may also hold general
// notes. It may be empty and should be trimmed of
// whitespace.
func (m Media) GetDescription() string {
	return m.Description
}

// GetLocalPath is the path to the media file within the
// local file system. There is no guarantee that the
// file exists, as it may have been deleted or moved
// since it was added to project. This may also happen
// if the project is moved to a new computer without
// copying the media files to matching locations in
// the new file system.
func (m Media) GetLocalPath() string {
	return m.LocalPath
}

func (ds *Datastore) ListMedia() ([]Media, error) {
	return ds.db.Select(Media{})
}

func (ds *Datastore) AddMedia(media Media) (Media, error) {
	media, e := media.Clean()
	if e == nil {
		e = ds.db.Insert(media)
	}
	return media, e
}

func (ds *Datastore) GetMediaById(id string) (Media, error) {
	return ds.db.SelectById(Media{}, id)
}

func (ds *Datastore) DeleteMediaById(id string) error {
	return ds.db.DeleteById(Media{}, id)
}

type Res = http.ResponseWriter

func (ds *Datastore) ServeHTTP(w Res, r *http.Request) {
	entityId := r.URL.Query().Get("entity_id")
	if entityId == "" {
		httpErrBadIdParam(w)
		return
	}

	media, e := ds.GetMediaById(entityId)
	if e != nil {
		httpErrMediaLookup(w, e)
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

var _ http.Handler = &Datastore{}
