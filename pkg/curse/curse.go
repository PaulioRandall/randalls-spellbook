package curse

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

// Curse is an error with a Message and optional Cause. If
// created via [TemplateCurse.Fmt] then [Curse.Is] will
// return true if called with itself or the original
// [TemplateCurse].
type Curse struct {
	errId   string
	Message string
	Cause   error
}

// Err creates a new curse with the given message.
func Err(message string) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: message,
	}
}

// Fmt creates a new curse with the given message and
// formatting arguments.
func Fmt(message string, args ...any) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: fmt.Sprintf(message, args...),
	}
}

// Plunder returns a new curse using the message of the
// passed error as the curse's message. It does not wrap
// the passed error. It will panic if the passed error is
// nil.
func Plunder(e error) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: e.Error(),
	}
}

// WrapIn wraps the curse in the passed symptom curse or
// template curse, replacing any existing cause, then
// returns the symptom.
func (cu Curse) WrapIn[T Curse | TemplateCurse](symptom T) T {
	if tc, ok := any(symptom).(TemplateCurse); ok {
		tc.curse.Cause = cu
		return any(tc).(T)
	}

	sc, _ := any(symptom).(Curse)
	sc.Cause = cu
	return any(sc).(T)
}

// WrapErr wraps the curse in a new curse with the given
// message.
func (cu Curse) WrapErr(message string) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: message,
		Cause:   cu,
	}
}

// WrapFmt wraps the curse in a new curse with the given
// message and formatting arguments.
func (cu Curse) WrapFmt(message string, args ...any) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: fmt.Sprintf(message, args...),
		Cause:   cu,
	}
}

// Wrap wraps the cause error replacing any existing
// cause.
func (cu Curse) Wrap(cause error) Curse {
	cu.Cause = cause
	return cu
}

// Unwrap returns the cause of the error, i.e. the wrapped
// error.
func (cu Curse) Unwrap() error {
	return cu.Cause
}

// Is returns true if the target is a [Curse] or
// [TemplateCurse] and has the same internal ID as the
// receiving curse.
func (cu Curse) Is(target error) bool {
	if tc, ok := target.(TemplateCurse); ok {
		return cu.errId == tc.curse.errId
	}

	if cu2, ok := target.(Curse); ok {
		return cu.errId == cu2.errId
	}

	return false
}

// Error returns the error message, satisfying Go's error
// interface.
func (cu Curse) Error() string {
	return cu.Message
}

// Stack returns the result of passing the Curse to
// [Stack].
func (cu Curse) Stack(reversed bool) string {
	return Stack(cu, reversed)
}

// RawStack returns the result of passing the Curse to
// [RawStack].
func (cu Curse) RawStack(reversed bool) string {
	return RawStack(cu, reversed)
}

// TemplateCurse is a [Curse] with a formattable message
// designed to be used as named package errors. Use
// [Template] to create one. [TemplateCurse.Fmt] should be
// called with the correct formatting arguments when
// returning the error. If the exported error needs no
// formatting then use [Err] to create the package error.
type TemplateCurse struct {
	curse Curse
}

// Template creates a [TemplateCurse] designed for use as
// named package errors tht may be coompared. When an error
// occurs, [TemplateCurse.Fmt] should be called with the
// relevant formatting arguments.
//
//	var ErrUserNotFound = curse.Template(
//		"User with ID '%d' not found",
//	)
//
//	func GetUserDetails(id int) (string, error) {
//		// ...
//		return false, ErrUserNotFound.Fmt(id)
//	}
func Template(message string) TemplateCurse {
	return TemplateCurse{
		curse: Err(message),
	}
}

// Error returns the error message, satisfying Go's error
// interface.
func (tc TemplateCurse) Error() string {
	log.Println("WARNING: Use of unformatted TemplateCurse")
	return tc.curse.Message
}

// Fmt formats the template's error message, which must be
// using the passed args and fmt.Sprintf. Calling
// [Curse.Is] with the receiving template will
// return true for all curses created from the receiving
// template.
func (tc TemplateCurse) Fmt(args ...any) Curse {
	cu := tc.curse
	cu.Message = fmt.Sprintf(cu.Message, args...)
	return cu
}

var _ error = Curse{}
var _ error = TemplateCurse{}
