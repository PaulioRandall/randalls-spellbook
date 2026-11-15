package sourcery

// Reset resets the application to its initial empty state
// removing all event handlers, message receivers, HTTP
// handlers, and options. Useful for testing.
func Reset() {
	if app.running {
		panic("Application can't be reset while running")
	}

	app = makeApplication()
}

// Run starts the application and blocks until it exits.
// Panic ensues if the application is already running.
func Run() error {
	// TODO: Copy code from other sourcery package.
	return nil
}

// Exit terminates the application causing [Run] to return.
// Does nothing if the application is not running.
func Exit() {
	// TODO
}
