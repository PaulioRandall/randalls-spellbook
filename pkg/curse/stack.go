package curse

import (
	"errors"
	"slices"
	"strings"
)

// Stack stringifies the error chain as formatted lines
// with the passed error first and the most wrapped error
// last. Each error message is trimmed of all content from
// the first colon ":" onwards, and the all but the first
// error messsage is indented and prefixed with angled
// arrows for ease of reading. If reverse is true then the
// root cause will be first and the passed error will last.
// The angled arrows always point towards the causes so you
// follow the arrows to learn the root cause.
func Stack(e error, reverse bool) string {
	var st []string

	for e != nil {
		msg := e.Error()
		causeIdx := strings.Index(msg, ":")

		if causeIdx > -1 {
			msg = msg[:causeIdx]
		}

		st = append(st, msg)
		e = errors.Unwrap(e)
	}

	if reverse {
		slices.Reverse(st)
		return strings.Join(st, "\n\t⮤ ")
	}

	return strings.Join(st, "\n\t⤷ ")
}

// RawStack is the same as [Stack] except messages are left
// untrimmed.
func RawStack(e error, reverse bool) string {
	var st []string

	for e != nil {
		st = append(st, e.Error())
		e = errors.Unwrap(e)
	}

	if reverse {
		slices.Reverse(st)
		return strings.Join(st, "\n\t⮤ ")
	}

	return strings.Join(st, "\n\t⤷ ")
}

type StackError struct {
	err     error
	raw     bool
	reverse bool
}

// StackErr returns an error that will return the error
// chain as a stack, according to [Stack] (default) or
// [RawStack].
func StackErr(err error) StackError {
	return StackError{
		err:     err,
		raw:     false,
		reverse: false,
	}
}

// Raw configure the returned stack error to call
// [RawStack] rather than [Stack]. Toggles so calling twice
// will revert back to using [Stack].
func (se StackError) Raw() StackError {
	se.raw = !se.raw
	return se
}

// Reverse reverses the error chain when calling Error so
// the root cause will be at the top of the stack. Toggles
// so calling twice will revert the change.
func (se StackError) Reverse() StackError {
	se.reverse = !se.reverse
	return se
}

// Error returns the error chain from calling [Stack] or
// [RawStack].
func (se StackError) Error() string {
	if se.raw {
		return RawStack(se.err, se.reverse)
	}
	return Stack(se.err, se.reverse)
}

var _ error = StackError{}
