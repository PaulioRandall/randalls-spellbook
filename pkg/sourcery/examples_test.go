package sourcery

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

type Counter struct {
	Value float64
}

func (c *Counter) Init(a *App) func() {
	// Initialisation code called on app start.
	c.Value = 0

	return func() {
		// Clean up code called on app exit.
		c.Value = -1
	}
}

func (c *Counter) Set(n float64) {
	c.Value = n
}

func (c *Counter) MultiplyBy(n float64) float64 {
	return c.Value * n
}

func (c *Counter) DivideBy(n float64) (float64, error) {
	if n == 0 {
		return 0, fmt.Errorf("You can't divide %v by %v", c.Value, n)
	}
	return c.Value / n, nil
}

func (c *Counter) Print() {
	fmt.Printf("Counter.Value == %v\n", c.Value)
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
		Width(400).
		Height(320).
		AddEntity(&Counter{}).
		AddServer("/", http.FileServerFS(webpage)).
		Start() // Blocks until WebView exits.

	if e != nil {
		log.Fatal(e)
	}
}
