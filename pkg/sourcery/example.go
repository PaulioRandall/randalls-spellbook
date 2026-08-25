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
	sour := NewSourcerer()

	sour.Debug(true)
	sour.Title("Example")
	sour.Size(420, 300)
	sour.Serve("/testdata/", http.FileServerFS(testdata))

	sour.Transcribe(
		"time",
		getTime,
		fmtTime,
	)

	e := sour.Conjure()

	if e != nil {
		log.Fatal(e)
	}
}

func getTime(_ Effect) Effect {
	return Bless(time.Now())
}

func fmtTime(input Effect) Effect {
	currTime, ok := input.Result().(time.Time)
	if !ok {
		return Curse("Wrong type, expected time.Time")
	}

	formattedTime := currTime.Format("13:14:15")
	return Bless(formattedTime)
}
