package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/PaulioRandall/randalls-spellbook/pkg/business"
	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery"
)

//go:embed ui/build/*
var webFiles embed.FS

func main() {
	rm := sourcery.NewRealm[data.Store]()

	rm.Debug(true)
	rm.Title("Randall's Spellbook")
	rm.Size(1200, 800)
	rm.Serve("/media/", business.NewMediaServer(rm))
	rm.Serve("/", createHtmlServer())

	rm.Transcribe(
		"SelectLocalFile",
		business.SelectLocalFile,
	)
	rm.Transcribe(
		"ListMedia",
		business.ListMedia,
	)
	rm.Transcribe(
		"GetMediaById",
		business.GetMediaById,
	)
	rm.Transcribe(
		"AddMedia",
		business.JsonToMedia,
		business.AddMedia,
	)

	rm.AfterOpening(afterOpening)
	rm.AfterClosing(afterClosing)

	// Blocks!
	e := rm.OpenPortal()

	if e != nil {
		log.Fatal(e)
	}
}

func createHtmlServer() http.Handler {
	fs, e := fs.Sub(webFiles, "ui/build")
	if e != nil {
		panic(e)
	}
	return http.FileServerFS(fs)
}

func afterOpening(rm business.Realm) error {
	rm.Inventory = data.NewStore("./testproject/data.sqlite")

	e := rm.Inventory.Open()
	if e != nil {
		return e
	}

	return nil
}

func afterClosing(rm business.Realm) error {
	var e error

	if rm.Inventory != nil {
		e = rm.Inventory.Close()
		rm.Inventory = nil
	}

	return e
}
