package main

import (
	"log"

	"github.com/PaulioRandall/randalls-spellbook/ui"
)

func main() {
	// TODO: Determine using cmd option -d.
	debug := true

	// Blocks!
	e := ui.Run(debug)

	if e != nil {
		log.Fatal(e)
	}
}
