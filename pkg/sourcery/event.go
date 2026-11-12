package sourcery

// EventType is the type of an application event. Used by
// [Event].
type EventType = string

const (
	// EventTypeLoad is emitted once just before the
	// application is started.
	EventTypeLoad EventType = "load"

	// EventTypeUnload is emitted once straight after the
	// application exits.
	EventTypeUnload EventType = "unload"
)

// Event provides access to event information. Passed to
// [Listener] functions.
type Event interface {
	// Type returns the event type.
	Type() EventType

	// App returns the application that triggered the event.
	App() *App
}

// Listener implementations receive notification for
// application events.
type Listener interface {
	OnEvent(Event)
}

type lifecycleEvent struct {
	eventType EventType
	app       *App
}

func (ev *lifecycleEvent) Type() EventType {
	return ev.eventType
}

func (ev *lifecycleEvent) App() *App {
	return ev.app
}
