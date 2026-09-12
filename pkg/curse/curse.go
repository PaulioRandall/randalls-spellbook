package curse

import (
	"fmt"
)

// Curse represents an error and satisfies the error
// interface. It implements functions needed to work with
// the standard errors package.
//
// Additional information maybe added to the Curse after it
// has been created which means standard equality checks
// may not work. Always use [Curse.Is] to test for equality
// as it accounts for those changes.
//
// TODO: Line number: runtime.Caller(1)
// TODO: File path: runtime.Caller(1)
// TODO: [Curse] prefix
type Curse struct {
	// Message is the original error message created on
	// construction.
	Message string

	// Cause is the wrapped error returned by [Curse.Unwrap].
	Cause error

	// original is the original error set by constructor
	// functions. It must remain unchanged for the Is
	// function to work correctly.
	original error
}

// New creates a new [Curse].
func New(message string) Curse {
	cu := Curse{
		Message: message,
	}

	cu.original = cu
	return cu
}

// Fmt creates a new [Curse] formatted using fmt.Sprintf.
func Fmt(message string, args ...any) Curse {
	cu := Curse{
		Message: fmt.Sprintf(message, args...),
	}

	cu.original = cu
	return cu
}

// Err creates a new [Curse]. Some may find the name
// preferable for readability.
func Err(message string) Curse {
	return New(message)
}

// Errf is an alias for [Fmt]. Some may find the name
// preferable for readability.
func Errf(message string, args ...any) Curse {
	return Fmt(message, args...)
}

// Error returns the full error message and satisfies the
// error interface.
func (cu Curse) Error() string {
	s := cu.Message

	if cu.Cause != nil {
		s += ":\n" + cu.Cause.Error()
	}

	return s
}

// Wrap sets the Cause field but also available using
// [Curse.Unwrap].
func (cu Curse) Wrap(e error) Curse {
	cu.Cause = e
	return cu
}

// Unwrap returns the Cause field, i.e. the cause of the
// curse. It staisfies the interface used by the functions
// in the errors package.
func (cu Curse) Unwrap() error {
	return cu.Cause
}

// Is returns true if the passed error is a Curse and both
// this and the passed Curses have the same original error;
// tested using '=='. It staisfies the interface used by
// the functions in the errors package.
func (cu Curse) Is(target error) bool {
	if cuss, ok := target.(Curse); ok {
		return cuss.original == cu.original
	}
	return false
}

func info(msg string, args ...any) string {
	return ":\n\t+ " + fmt.Sprintf(msg, args...)
}

func causedBy(e error) string {
	return ":\n\t+ " + e.Error()
}

var _ error = Curse{}
