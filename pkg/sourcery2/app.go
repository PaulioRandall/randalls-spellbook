package sourcery

import (
	"net/http"
)

// TODO
//
//	- Panic when http handlers are added when the system
//		is running.

type application struct {
	// Config
	debug  bool
	title  string
	width  int
	height int

	// Handlers
	httpHandlers   map[string]http.Handler
	msgHandlers    map[string]MessageHandler
	eventListeners map[eventListenerKey]eventListenerValue

	// Runtime
	running         bool
	doneChan        chan bool
	eventQueue      chan Event
	eventErrorQueue chan Event
}

var app = makeApplication()

func makeApplication() application {
	return application{
		// Options
		debug:  false,
		title:  "Technotelicomnicon",
		width:  800,
		height: 600,

		// Handlers
		httpHandlers:   map[string]http.Handler{},
		msgHandlers:    map[string]MessageHandler{},
		eventListeners: map[eventListenerKey]eventListenerValue{},

		// Runtime
		running:         false,
		doneChan:        make(chan bool),
		eventQueue:      make(chan Event, 8),
		eventErrorQueue: make(chan Event, 4),
	}
}
