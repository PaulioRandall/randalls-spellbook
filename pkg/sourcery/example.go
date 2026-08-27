package sourcery

import (
	"embed"
	"log"
	"net/http"
	"time"

	. "github.com/PaulioRandall/randalls-spellbook/pkg/spellbook"
)

//go:embed testdata
var testdata embed.FS

func SourceryExample() {
	rm := NewRealm()

	rm.Debug(true)
	rm.Title("Example")
	rm.Size(420, 300)
	rm.Serve("/testdata/", http.FileServerFS(testdata))

	rm.Transcribe(
		"time",
		getTime,
		fmtTime,
	)

	rm.AfterOpening(afterOpening)
	rm.AfterClosing(afterClosing)

	e := rm.OpenPortal()

	if e != nil {
		log.Fatal(e)
	}
}

func getTime(_ any, _ Effect) Effect {
	return Bless(time.Now())
}

func fmtTime(_ any, input Effect) Effect {
	currTime, ok := input.Result().(time.Time)
	if !ok {
		return Curse("Wrong type, expected time.Time")
	}

	formattedTime := currTime.Format("13:14:15")
	return Bless(formattedTime)
}

func afterOpening(rm *Realm) error {
	println("Realm is open!")
	return nil
}

func afterClosing(rm *Realm) error {
	println("Realm is closed!")
	return nil
}
