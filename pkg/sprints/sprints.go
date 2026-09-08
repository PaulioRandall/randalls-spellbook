package sprints

import (
	"fmt"
)

const (
	OptionShowUnexported = "OptionShowUnexported"
)

var (
	ErrBadIndentCount = fmt.Errorf(
		"Integer must follow OptionIndent",
	)
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

// String formats an objects into a string form similar
// to the object's definition or instantiation. If a
// non-struct kind is passed then panic ensues.
func String(object any, options ...string) string {
	return stringifyObject(
		object,
		newFmtCtx(options),
	)
}
