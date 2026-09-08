package sprints

import (
	"fmt"
)

const (
	// OptionShowUnexported shows unexported field names but
	// won't show their values.
	OptionShowUnexported = "OptionShowUnexported"
)

// Print stringifies the object and prints it to terminal.
func Print(object any, options ...string) {
	s := stringifyObject(
		object,
		newFmtCtx(options),
	)
	fmt.Print(s)
}

// Println stringifies the object and prints it to
// terminal followed by a linefeed.
func Println(object any, options ...string) {
	s := stringifyObject(
		object,
		newFmtCtx(options),
	)
	fmt.Println(s)
}

// String returns the stringified object.
func String(object any, options ...string) string {
	return stringifyObject(
		object,
		newFmtCtx(options),
	)
}
