package sourcery

import (
	"fmt"
	"net/http"
)

// ServeHandler adds an endpoint to the HTTP server. Panic
// ensues if the path is already assigned a handler.
func ServeHandler(path string, handler http.Handler) {
	if _, ok := app.httpHandlers[path]; ok {
		err := fmt.Sprintf(
			"HTTP handler alreaddy exists for path: %s",
			path,
		)
		panic(err)
	}

	app.httpHandlers[path] = handler
}

// Serve is the same as [ServeHandler] except it accepts a
// http.HandlerFunc.
func Serve(path string, handler http.HandlerFunc) {
	ServeHandler(path, http.HandlerFunc(handler))
}
