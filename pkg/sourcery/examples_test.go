package sourcery

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

type Counter struct {
	w     *World
	value int
}

func (c *Counter) Bind(w *World) {
	c.w = w
	c.value = 0
}

func (c *Counter) Free() {
	fmt.Printf("Counter.value = %d\n", c.value)
}

func (c *Counter) Add(n int) {
	c.value += n
}

func (c *Counter) ExitApp() {
	c.w.WebView().Terminate()
}

//go:embed testdata
var testdata embed.FS

func Example() {
	webpage, e := fs.Sub(testdata, "testdata")
	if e != nil {
		log.Fatal(e)
	}

	e = New().
		Name("Example App").
		Size(400, 320).
		Bind("Counter", &Counter{}).
		Serve("/", http.FileServerFS(webpage)).
		BuildWorld().
		Enter() // Blocks until WebView closes.

	if e != nil {
		log.Fatal(e)
	}
	// Output:
	// Counter.value = 3
}
