package sourcery

import (
	"fmt"
)

// EventArg returns the argument at the given index and
// true if the event contains an argument, else T's zero
// value is retuned with false.
func EventArg[T any](ev Event, idx int) (T, bool) {
	var zero T

	v, ok := ev.Arg(idx)
	if !ok {
		return zero, false
	}

	arg, ok := v.(T)
	if !ok {
		return zero, false
	}

	return arg, true
}

// Event provides an interface to event details and
// functionalty for influencing events.
type Event interface {
	// Type returns the event's type.
	Type() string

	// Arg returns the argument at the given index and true.
	// Else nil is retuned with false. Use [EventArg] as a
	// convenience for returning a typed value.
	Arg(idx int) (any, bool)
}

type stdEvent struct {
	eventType string
	args      []any
}

func (ev stdEvent) Type() string {
	return ev.eventType
}

func (ev stdEvent) Arg(idx int) (any, bool) {
	if len(ev.args) >= idx {
		return nil, false
	}
	return ev.args[idx], true
}

// EventListener is registered through [On] and [Once] to
// receive events with a specific type.
type EventListener = func(Event)

type eventListenerKey struct {
	nameSpace string
	eventType string
}

type eventListenerValue struct {
	listener         EventListener
	removeOnDispatch bool
}

// On registers a listener for the specified nameSpace
// and eventType. Panic ensues if either nameSpace or
// eventType is empty or the combination of
// nameSpace and eventType is already assigned.
func On(nameSpace, eventType string, listener EventListener) {
	onEvent(nameSpace, eventType, listener, false)
}

// Once is the same as [On] except the listener is removed
// after the first dispatch, i.e. called only once.
func Once(nameSpace, eventType string, listener EventListener) {
	onEvent(nameSpace, eventType, listener, true)
}

func onEvent(
	nameSpace, eventType string,
	listener EventListener,
	once bool,
) {
	if nameSpace == "" {
		panic("Event listener name space cannot be empty")
	}

	if eventType == "" {
		panic("Event listener event type cannot be empty")
	}

	key := eventListenerKey{
		nameSpace: nameSpace,
		eventType: eventType,
	}

	if _, ok := app.eventListeners[key]; ok {
		err := fmt.Sprintf(
			"Event listener already exists for name space and event type combination: '%s' and '%s'",
			nameSpace,
			eventType,
		)
		panic(err)
	}

	app.eventListeners[key] = eventListenerValue{
		listener:         listener,
		removeOnDispatch: once,
	}
}

// Off unregisters the listener with the specified
// nameSpace and eventType combination.
func Off(nameSpace, eventType string) {
	if nameSpace == "" && eventType == "" {
		// Remove all entries.
		clear(app.eventListeners)
		return
	}

	if nameSpace != "" && eventType != "" {
		// Remove specific entry.
		delete(app.eventListeners, eventListenerKey{
			nameSpace: nameSpace,
			eventType: eventType,
		})
		return
	}

	if nameSpace == "" {
		// Remove eventType.
		for key := range app.eventListeners {
			if key.eventType == eventType {
				delete(app.eventListeners, key)
			}
		}
		return
	}

	// Remove nameSpace.
	for key := range app.eventListeners {
		if key.nameSpace == nameSpace {
			delete(app.eventListeners, key)
		}
	}
}

// Post dispatches an event to all handlers registered to
// the eventType. Nothing happens if there are no handlers
// registered.
//
// The event system should be viewed as an independent
// system. It runs on its own Go routine in an attempt to
// follow the principles of the Reactive Manifesto.
//
// Calling Post adds a request to the dispatch queue then
// returns. The event system will process the queue in its
// own time dispatching the posted events to all handlers
// one-by-one. This means if an event handler posts an
// event it will not be dispatched until all handlers of
// the current event have been called. The only exception
// is when a 'sourcery:error' event is posted; error
// events are processed immediately after the current
// handler returns, i.e. they have there own queue which is
// prioritised. If any event handler panics then an event
// is dispatched and processed immediately, with the
// recovered value as the only argument to the error
// handlers. Of course the only exception is if a
// 'sourcery:error' event handler panics, then all that
// happens is the recovered value is printed.
func Post(eventType string, args ...any) {
	app.eventQueue <- stdEvent{
		eventType,
		args,
	}
}

func handleEvents() {
	for {
		select {
		case ev := <-app.eventQueue:
			// TODO: Impl as according to [Post] func docs
			_ = ev
		case <-app.doneChan:
			return
		}
	}
}
