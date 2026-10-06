package sin

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

// Curse is an error with a Message and optional Cause. If
// created via [TemplateCurse.Fmt] then [Curse.Is] will
// return true if called with itself or any curse created
// via the original [TemplateCurse].
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

// Hijack is the same as [Plunder] except the error chain
// is also plundered.
func Hijack(e error) Curse {
	type unwrapper interface {
		Unwrap() error
	}

	cu := Curse{
		errId:   uuid.New().String(),
		Message: e.Error(),
	}

	if un, ok := e.(unwrapper); ok {
		cu.Cause = un.Unwrap()
	}

	return cu
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
// message. The receiving curse is returned.
func (cu Curse) WrapErr(message string) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: message,
		Cause:   cu,
	}
}

// WrapFmt wraps the curse in a new curse with the given
// message and formatting arguments. The receiving curse is
// returned.
func (cu Curse) WrapFmt(message string, args ...any) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: fmt.Sprintf(message, args...),
		Cause:   cu,
	}
}

// Wrap wraps the cause error replacing any existing
// cause. The receiving curse is returned.
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
	return cu.StackString()
}

// StackString returns the result of passing the Curse to
// [StackString].
func (cu Curse) StackString() string {
	return StackString(cu)
}

// ReverseStackString returns the result of passing the Curse
// to [ReverseStackString].
func (cu Curse) ReverseStackString() string {
	return ReverseStackString(cu)
}

// AsStackError wraps the error in a StackError.
func (cu Curse) AsStackError() error {
	return Stack(cu)
}

// AsReverseStackError wraps the error in a StackError.
func (cu Curse) AsReverseStackError() error {
	return ReverseStack(cu)
}

// TemplateCurse creates [Curse] errors with formattable
// messages. It is designed to be used as named package
// errors. [TemplateCurse.Fmt] should be called at the site
// of an error to create an error that is returned. If the
// exported error needs no formatting then create the
// package error with [Err].
type TemplateCurse struct {
	curse Curse
}

// Template creates a [TemplateCurse] designed for use as
// named package errors. Calling [Curse.Is] with the
// template that created it will return true.
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
