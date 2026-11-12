package sourcery

import (
	"net/http"
)

type MyEntity struct {
	app *App
}

func (me *MyEntity) OnEvent(e Event) {
	switch e.Type() {
	case EventTypeLoad: // Intiialise.
		me.app = e.App()
	case EventTypeUnload: // Clean up.
		me.app = nil
	}
}

func (me *MyEntity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (me *MyEntity) CallMe() {
	println("Thanks for calling me!")
}

func ExampleApp() {
	entity := &MyEntity{}

	app := NewApp(
		OptionDebug(true),
		OptionTitle("Example App"),
		OptionSize(1000, 600),
		OptionListener(entity),
		OptionServer("/", entity),
		OptionEntity(entity),
		OptionHtml(`
<!DOCTYPE html>
<html>
	<body>
		<p>Sourcery</p>
	</body>
</html>
`),
	)

	_ = app
}

func ExampleApp_package() {
	entity := &MyEntity{}

	Append(
		OptionDebug(true),
		OptionTitle("Example App"),
		OptionSize(1000, 600),
		OptionListener(entity),
		OptionServer("/", entity),
		OptionEntity(entity),
		OptionHtml(`
<!DOCTYPE html>
<html>
	<body>
		<p>Sourcery</p>
	</body>
</html>
`),
	)

	_ = DefaultApp()
}
