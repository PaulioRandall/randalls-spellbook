package main

import (
	"log"
	
	"github.com/crgimenes/glaze"
)

func main() {
	w, err := glaze.New(true)
	
	if err != nil {
		log.Fatal(err)
	}

	defer w.Destroy()

	w.SetTitle("Glaze")
	w.SetSize(800, 600, glaze.HintNone)
	w.SetHtml("<h1>Hello from Glaze</h1>")
	
	w.Run()
}
