package sourcery

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithDebug_1(t *testing.T) {
	// Enables debug mode.
	Reset()

	WithDebug(true)
	require.Equal(t, true, app.debug)
}

func Test_WithTitle_1(t *testing.T) {
	// Sets title.
	Reset()

	WithTitle("Abc")
	require.Equal(t, "Abc", app.title)
}

func Test_WithSize_1(t *testing.T) {
	// Sets size.
	Reset()

	WithSize(10, 20)
	require.Equal(t, 10, app.width)
	require.Equal(t, 20, app.height)
}

func Test_Serve_1(t *testing.T) {
	// Implicitly tests ServeHandler function too.
	// Registers http handler.
	Reset()

	Serve("/", func(w http.ResponseWriter, r *http.Request) {
		// Meh
	})

	require.Equal(t, 1, len(app.httpHandlers))
	_, found := app.httpHandlers["/"]
	require.Equal(t, true, found)
}

func Test_Serve_2(t *testing.T) {
	// Implicitly tests ServeHandler function too.
	// Panics when path used more than once.
	Reset()

	Serve("/", nil)
	require.Panics(t, func() {
		Serve("/", nil)
	})
}

func Test_Receive_1(t *testing.T) {
	// Registers message handler.
	Reset()

	Receive("abc", func(...any) (any, error) {
		// Meh
		return nil, nil
	})

	require.Equal(t, 1, len(app.msgHandlers))
	_, found := app.msgHandlers["abc"]
	require.Equal(t, true, found)
}

func Test_Receive_2(t *testing.T) {
	// Panics when address used more than once.
	Reset()

	Receive("abc", nil)
	require.Panics(t, func() {
		Receive("abc", nil)
	})
}

func Test_Resign_1(t *testing.T) {
	// Registers message handler.
	Reset()

	app.msgHandlers["abc"] = func(...any) (any, error) {
		// Meh
		return nil, nil
	}

	Resign("abc")

	require.Equal(t, 0, len(app.msgHandlers))
}

func Test_On_1(t *testing.T) {
	// Registers event listener.
	Reset()

	On("ns", "et", func(_ Event) {
		// Meh
	})

	require.Equal(t, 1, len(app.eventListeners))

	key := eventListenerKey{
		nameSpace: "ns",
		eventType: "et",
	}
	entry, found := app.eventListeners[key]

	require.Equal(t, true, found)
	require.NotNil(t, entry.listener)
	require.Equal(t, false, entry.removeOnDispatch)
}

func Test_On_2(t *testing.T) {
	// Panics when a key is used more than once.
	Reset()

	On("ns", "et", nil)
	require.Panics(t, func() {
		On("ns", "et", nil)
	})
}

func Test_Once_1(t *testing.T) {
	// Registers event listener.
	Reset()

	Once("ns", "et", func(_ Event) {
		// Meh
	})

	require.Equal(t, 1, len(app.eventListeners))

	key := eventListenerKey{
		nameSpace: "ns",
		eventType: "et",
	}
	entry, found := app.eventListeners[key]

	require.Equal(t, true, found)
	require.NotNil(t, entry.listener)
	require.Equal(t, true, entry.removeOnDispatch)
}

func Test_Off_1(t *testing.T) {
	// Removes event listener.
	Reset()

	key := eventListenerKey{
		nameSpace: "ns",
		eventType: "et",
	}
	app.eventListeners[key] = eventListenerValue{}

	Off("ns", "et")

	require.Equal(t, 0, len(app.eventListeners))
}

func Test_Off_2(t *testing.T) {
	// Removes event listeners.
	Reset()

	key1 := eventListenerKey{
		nameSpace: "ns",
		eventType: "et1",
	}
	app.eventListeners[key1] = eventListenerValue{}

	key2 := eventListenerKey{
		nameSpace: "ns",
		eventType: "et2",
	}
	app.eventListeners[key2] = eventListenerValue{}

	Off("ns", "")

	require.Equal(t, 0, len(app.eventListeners))
}

func Test_Off_3(t *testing.T) {
	// Removes event listeners.
	Reset()

	key1 := eventListenerKey{
		nameSpace: "ns1",
		eventType: "et",
	}
	app.eventListeners[key1] = eventListenerValue{}

	key2 := eventListenerKey{
		nameSpace: "ns2",
		eventType: "et",
	}
	app.eventListeners[key2] = eventListenerValue{}

	Off("", "et")

	require.Equal(t, 0, len(app.eventListeners))
}
