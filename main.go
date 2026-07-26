package main

import (
	"log"

	"randalls-spellbook/view"
)

func main() {
	v, err := view.New(true)

	if err != nil {
		log.Fatal(err)
	}

	defer v.Destroy()

	w.SetTitle("Glaze")
	w.SetSize(800, 600, glaze.HintNone)
	w.SetHtml("<h1>Hello from Glaze</h1>")

	w.Run()
}
