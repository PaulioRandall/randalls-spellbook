package example

import (
	"fmt"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sourcery2"
)

func Example() {
	sourcery.WithDebug(true)
	sourcery.WithTitle("Example App")
	sourcery.WithSize(800, 600)

	sourcery.On("main", "sourcery:start", func(_ sourcery.Event) {
		println("Application started")
	})

	sourcery.On("main", "sourcery:exit", func(_ sourcery.Event) {
		println("Application exiting")
	})

	sourcery.On("main", "sourcery:error", func(ev sourcery.Event) {
		err, _ := sourcery.EventArg[error](ev, 0)
		fmt.Println(err)
	})

	sourcery.Run()
}
