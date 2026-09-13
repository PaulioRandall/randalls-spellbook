package curse

import (
	"fmt"

	"github.com/google/uuid"
)

// Curse represents an error that may have additional
// information attached to it.
type Curse interface {
	error
	Attached() []string
	Attach(string) Curse
}

// Attach attaches an info string to err if err can be cast
// to a [Curse]. If not, a curse is created using err's
// error message and the information attached to that,
// which is then returned with err as the cause.
func Attach(err error, info string) Curse {
	var cu Curse
	var ok bool

	if err == nil {
		panic("Cannot attach info to a nil error")
	}

	if cu, ok = err.(Curse); !ok {
		cu = Wrap(err, err.Error())
	}

	return cu.Attach(info)
}

// Hex is a simple implementation of [Curse].
type Hex struct {
	errId       string
	msg         string
	cause       error
	attachments []string
}

// Err returns a new [Curse].
func Err(msg string, args ...any) Hex {
	return newHex(nil, msg, args...)
}

// Wrap returns a new [Curse] with a known cause error.
func Wrap(cause error, msg string, args ...any) Hex {
	return newHex(cause, msg, args...)
}

func newHex(cause error, msg string, args ...any) Hex {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return Hex{
		errId: uuid.New().String(),
		msg:   msg,
	}
}

// Attached returns the list off attached inforrmation.
func (h Hex) Attached() []string {
	return h.attachments
}

// Attach attaches new info to the curse.
func (h Hex) Attach(info string) Curse {
	h.attachments = append(h.attachments, info)
	return h
}

// Wrap wraps the cause error replacing any existing cause.
// This is useful when creating package level exported
// curses (errors). The API producer can call Wrap on the
// exported error, returning the result as a new instance
// of the error that will return true when one is compared
// to the other using the Is function.
func (h Hex) Wrap(cause error) Curse {
	h.cause = cause
	return h
}

// Unwrap returns the cause of the error, i.e. the wrapped
// error.
func (h Hex) Unwrap() error {
	return h.cause
}

// Is returns true if the target is the same type and has
// the same ID as the receiving error.
func (h Hex) Is(target error) bool {
	if h2, ok := target.(Hex); ok {
		return h.errId == h2.errId
	}
	return false
}

// Error returns the error message along with any attached
// information. It satisfies Go's error interface.
func (h Hex) Error() string {
	const prefix string = "\n\t+ "
	s := h.msg

	for _, info := range h.attachments {
		s += prefix + info
	}

	if h.cause != nil {
		s += "\nCaused by: " + h.cause.Error()
	}

	return s
}

func newErrLine(msg string, args ...any) string {
	return "\n\t+ " + fmt.Sprintf(msg, args...)
}

var _ Curse = Hex{}
