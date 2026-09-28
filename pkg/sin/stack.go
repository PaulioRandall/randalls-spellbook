package sin

import (
	"errors"
	"slices"
	"strings"
)

// AsStack stringifies the error chain as formatted
// lines with the passed error first and the error at the
// end of the chain last. Each error message is trimmed of
// all content from the first colon ":" onwards, and the
// all but the first error messsage is indented and
// prefixed with corner arrows for ease of reading. If
// reverse is true then the root cause will be first and
// the passed error will last.
//
// The corner arrows always point from the symptom to the
// cause so following the arrows will lead you to the root
// cause.
func AsStack(e error, reverse bool) string {
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
		return strings.Join(st, "\n\t⤤ ")
	}

	return strings.Join(st, "\n\t⤷ ")
}

// AsRawStack is the same as [AsStack] except
// messages are left untrimmed.
func AsRawStack(e error, reverse bool) string {
	var st []string

	for e != nil {
		st = append(st, e.Error())
		e = errors.Unwrap(e)
	}

	if reverse {
		slices.Reverse(st)
		return strings.Join(st, "\n\t⤤ ")
	}

	return strings.Join(st, "\n\t⤷ ")
}

// StackError is returned by [Stack] and
// produces an error string according to [AsStack] or
// [AsRawStack] when [StackError.Error] is
// called. StackError is designed to summarise the chain
// for easy read output.
type StackError struct {
	err     error
	raw     bool
	reverse bool
}

// Stack returns an error that will return the error
// chain like a stack trace, according to [AsStack]
// (default), when [StackError.Error] is called. If raw is
// true then [AsRawStack] is used instead of AsStack. If
// reverse is true then the stack is reversed so the root
// cause appears first. If the error is nil then nil is
// returned.
func Stack(err error, raw bool, reverse bool) error {
	if err == nil {
		return nil
	}

	return StackError{
		err:     err,
		raw:     raw,
		reverse: reverse,
	}
}

// Error returns the error chain like a stack trace from
// calling [AsStack] or [AsRawStack].
func (se StackError) Error() string {
	if se.raw {
		return AsRawStack(se.err, se.reverse)
	}
	return AsStack(se.err, se.reverse)
}

// Is performs standard equality test against the
// receiver.
func (se StackError) Is(target error) bool {
	return target == se
}

// Unwrap returns the wrapped error.
func (se StackError) Unwrap() error {
	return se.err
}

var _ error = StackError{}
