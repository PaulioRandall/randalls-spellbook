package main

import (
	"log"

	"github.com/PaulioRandall/randalls-spellbook/ui"
)

func main() {
	e := ui.Run(true)

	if e != nil {
		log.Fatal(e)
	}
}
