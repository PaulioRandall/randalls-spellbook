package spellbook

import (
	"fmt"
)

// Incantation is a single step when casting a spell. It
// may transform data, create a side effect (e.g. write to
// storage), check data and generate an error, or a
// combination of these.
//
// The first thing most incantations will do is parse and
// check the input data. This is the price paid for such a
// decoupled architecture. However, since most Incantations
// expect a single specific ctx type and data type the
// Demysitfy and Demystifyf functions can be used to
// quickly cast to type; if type casting fails then panic
// ensues.
type Incantation func(ctx any, data any) Effect

// Spell is a series of incantations that are to be invoked
// in order.
type Spell []Incantation

// Demystify casts the value to type T. If casting fails
// panic ensues.
func Demystify[T any](value any) T {
	if tv, ok := any(value).(T); ok {
		return tv
	}

	panic("Unexpected type")
}

// Demystify casts the value to type T. If casting fails
// panic ensues with the formatted message.
func Demystifyf[T any](
	value any,
	message string,
	args ...any,
) T {
	tv, ok := value.(T)
	if ok {
		return tv
	}

	errMsg := fmt.Sprintf(message, args...)
	panic(errMsg)
}
