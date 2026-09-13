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
		cu = HexWrap(err, err.Error())
	}

	return cu.Attach(info)
}

type hexCurse struct {
	errId       string
	msg         string
	cause       error
	attachments []string
}

// Hex returns a new [Curse].
func Hex(msg string, args ...any) Curse {
	return newHex(nil, msg, args...)
}

// HexWrap returns a new [Curse] with a knwon cause.
func HexWrap(cause error, msg string, args ...any) Curse {
	return newHex(cause, msg, args...)
}

func newHex(cause error, msg string, args ...any) hexCurse {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	return hexCurse{
		errId: uuid.New().String(),
		msg:   msg,
	}
}

// Attached returns the list off attached inforrmation.
func (h hexCurse) Attached() []string {
	return h.attachments
}

// Attach attaches new info to the curse.
func (h hexCurse) Attach(info string) Curse {
	h.attachments = append(h.attachments, info)
	return h
}

// Unwrap returns the cause of the error, i.e. the wrapped
// error.
func (h hexCurse) Unwrap() error {
	return h.cause
}

// Is returns true if the target is the same type and has
// the same ID as the receiving error.
func (h hexCurse) Is(target error) bool {
	if h2, ok := target.(hexCurse); ok {
		return h.errId == h2.errId
	}
	return false
}

// Error returns the error message along with any attached
// information. It satisfies Go's error interface.
func (h hexCurse) Error() string {
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

var _ error = hexCurse{}
var _ Curse = hexCurse{}
