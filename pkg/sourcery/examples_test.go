package sourcery

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

type Box struct {
	w     *World
	value int
}

func (b *Box) WorldEnter(w *World) {
	b.w = w
	fmt.Printf("Initial box value is %d\n", b.value)
}

func (b *Box) WorldExit() {
	fmt.Printf("Final box value is %d\n", b.value)
}

func (b *Box) GetValue() int {
	return b.value
}

func (b *Box) SetValue(n int) {
	b.value = n
	fmt.Printf("New box value is %d\n", b.value)

	if b.value >= 5 {
		b.w.Exit()
	}
}

//go:embed testdata
var testdata embed.FS

func Example() {
	/*
		type Box struct {
			w     *World
			value int
		}

		func (b *Box) WorldEnter(w *World) {
			b.w = w
			fmt.Printf("Initial box value is %d\n", b.value)
		}

		func (b *Box) WorldExit() {
			fmt.Printf("Final box value is %d\n", b.value)
		}

		func (b *Box) GetValue() int {
			return b.value
		}

		func (b *Box) SetValue(n int) {
			b.value = n
			fmt.Printf("New box value is %d\n", b.value)

			if b.value >= 5 {
				b.w.Exit()
			}
		}

		//go:embed testdata
		var testdata embed.FS

		<!DOCTYPE html>
		<html>
			<body>
				<p>Count: <span id="value">0</span><p>
				<p id="error" color="red"></p>

				<script>
					function onError(e) {
						document.getElementById('error').innerHTML = e
						throw e
					}

					function incrementValue() {
						Go("Box.GetValue")
							.then(value => {
								value++
								document.getElementById('value').innerHTML = value
								Go("Box.SetValue", value).catch(onError)
								setTimeout(incrementValue, delayMS)
							})
							.catch(onError)
					}

					const delayMS = 250
					setTimeout(incrementValue, delayMS)
				</script>
			</body>
		</html>
	*/

	webpage, e := fs.Sub(testdata, "testdata")
	if e != nil {
		log.Fatal(e)
	}

	e = NewCreator().
		Name("Example App").
		Size(400, 320).
		AddPortal(&Box{}).
		AddServer("/", http.FileServerFS(webpage)).
		BuildWorld().
		Enter() // Blocks until WebView closes.

	if e != nil {
		log.Fatal(e)
	}
	// Output:
	// Initial box value is 0
	// New box value is 1
	// New box value is 2
	// New box value is 3
	// New box value is 4
	// New box value is 5
	// Final box value is 5
}
