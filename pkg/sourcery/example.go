package sourcery

import (
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/PaulioRandall/randalls-spellbook/pkg/data"
	. "github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

//go:embed testdata
var testdata embed.FS

func sourceryExample() {
	rm := NewRealm[data.Store]()

	rm.Debug(true)
	rm.Title("Example")
	rm.Size(420, 300)
	rm.Serve("/testdata/", http.FileServerFS(testdata))

	rm.spellbook.Transcribe(
		"time",
		exampleGetTime,
		exampleFmtTime,
	)

	rm.AfterOpening(exampleAfterOpening)
	rm.AfterClosing(exampleAfterClosing)

	e := rm.OpenPortal()

	if e != nil {
		log.Fatal(e)
	}
}

func exampleGetTime(_ any, _ any) Effect {
	return Bless(time.Now())
}

func exampleFmtTime(_ any, data any) Effect {
	currTime := Demystify[time.Time](data)
	formattedTime := currTime.Format("13:14:15")
	return Bless(formattedTime)
}

func exampleAfterOpening(rm *Realm[data.Store]) error {
	println("Realm is open!")
	return nil
}

func exampleAfterClosing(rm *Realm[data.Store]) error {
	println("Realm is closed!")
	return nil
}
