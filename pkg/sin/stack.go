package sin

import (
	"errors"
	"slices"
	"strings"
)

// ErrorChain recursively unwraps the error and returns the
// set of errors as a slice with the passed error at the
// front. If the error is nil then nil is returned.
func ErrorChain(err error) []string {
	var chain []string

	for err != nil {
		var msg string

		if cu, ok := err.(Curse); ok {
			msg = cu.Message
		} else {
			msg = err.Error()
		}

		chain = append(chain, msg)
		err = errors.Unwrap(err)
	}

	return chain
}

// StackString stringifies the error chain as formatted
// lines with the passed error first. All but the first
// error messsage is indented and prefixed with arrows
// pointing towards the cause for ease of reading. If the
// error is nil then an empty string is returned.
func StackString(err error) string {
	if err == nil {
		return ""
	}
	chain := ErrorChain(err)
	return strings.Join(chain, "\n\t⤷ ")
}

// ReverseStackString is the same as [StackString] except
// the chain is reversed when stringified.
func ReverseStackString(err error) string {
	if err == nil {
		return ""
	}
	chain := ErrorChain(err)
	slices.Reverse(chain)
	return strings.Join(chain, "\n\t⮤ ")
}

// StackError is returned by [Stack] and
// produces an error string according to [StackString] or
// [AsRawStack] when [StackError.Error] is
// called. StackError is designed to summarise the chain
// for easy read output.
type StackError struct {
	err     error
	reverse bool
}

// Stack returns an error that will return the error
// chain like a stack trace, according to [StackString],
// when the Error method is called.
func Stack(err error) error {
	if err == nil {
		return nil
	}

	return StackError{
		err:     err,
		reverse: false,
	}
}

// Stack returns an error that will return the error
// chain like a stack trace, according to
// [ReverseStackString], when the Error method is called.
func ReverseStack(err error) error {
	if err == nil {
		return nil
	}

	return StackError{
		err:     err,
		reverse: true,
	}
}

// Error returns the error chain like a stack trace from
// calling [StackString] or [ReverseStackString].
func (se StackError) Error() string {
	if se.reverse {
		return ReverseStackString(se.err)
	}
	return StackString(se.err)
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
