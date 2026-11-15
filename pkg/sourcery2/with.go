package sourcery

// WithDebug enables or disables debugging.
func WithDebug(enable bool) {
	app.debug = enable
}

// WithTitle sets application title.
func WithTitle(title string) {
	app.title = title
}

// WithSize sets the initial size of the application
// window.
func WithSize(width, height int) {
	app.width = width
	app.height = height
}

// WithEventBuffer sets the buffer size for the event
// queue. Default size is 8. If you notice calls to [Post]
// are blocking then try redesigning your app or doubling
// the buffer size.
func WithEventBuffer(size int) {
	app.eventQueue = make(chan Event, size)
}

// WithEventErrorBuffer sets the buffer size for the event
// error queue. Default size is 4.
func WithEventErrorBuffer(size int) {
	app.eventErrorQueue = make(chan Event, size)
}
