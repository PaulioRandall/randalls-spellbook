package curse

import (
	"fmt"
	"runtime"
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

	// Extra contains a list of additional information
	// messages to be printed with after the main message,
	// but before the cause error is printed.
	Extra []string

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

// Newf creates a new [Curse] formatted using fmt.Sprintf.
func Newf(message string, args ...any) Curse {
	cu := Curse{
		Message: fmt.Sprintf(message, args...),
	}

	cu.original = cu
	return cu
}

// Err creates a new [Curse] that includes file and line
// information by default.
func Err(message string) Curse {
	cu := Curse{
		Message: message,
	}

	cu.original = cu
	return cu.Pos()
}

// Fmt creates a new [Curse] formatted using fmt.Sprintf
// and includes file and line information by default.
func Fmt(message string, args ...any) Curse {
	cu := Curse{
		Message: fmt.Sprintf(message, args...),
	}

	cu.original = cu
	return cu.Pos()
}

// Error returns the full error message and satisfies the
// error interface.
func (cu Curse) Error() string {
	msg := cu.Message

	for _, ex := range cu.Extra {
		msg += "\n\t" + ex
	}

	return msg + causedBy(cu.Cause)
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

// Pos adds file and line information about the call to
// this method. If called straight after the error was
// created then it serves as an approximate location where
// thee error was created. Curses created using [Err] or
// [Fmt] will already provide this information. It's
// possible the file and line information cannot be
// obtained, in which case the error string will note this.
func (cu Curse) Pos() Curse {
	cu.Extra = append(cu.Extra, lineInfo(2))
	return cu
}

func causedBy(e error) string {
	if e != nil {
		return "\nCaused by: " + e.Error()
	}
	return ""
}

func lineInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)

	if ok {
		return fmt.Sprintf("%d: %s", line, file)
	} else {
		return "?: <unknown>"
	}
}

type fullError interface {
	error
	Is(error) bool
	Unwrap() error
}

var _ fullError = Curse{}
