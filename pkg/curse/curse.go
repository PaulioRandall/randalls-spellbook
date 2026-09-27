package curse

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

// TODO: Rename 'ProtoCurse' to 'Template' and rename
//       functions accordingly.
// TODO: Consider renaming package to 'sin' but leaving
//       main error type as 'Curse'.

// Curse is an error with a Message and Cause. If
// returned by [ProtoCurse.Fmt] then [Curse.Is] will return
// true if called with itself or the original [ProtoCurse].
type Curse struct {
	errId   string
	Message string
	Cause   error
}

// Err creates a new Curse with the given message.
func Err(message string) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: message,
	}
}

// Fmt creates a new Curse with the given message and
// formatting arguments.
func Fmt(message string, args ...any) Curse {
	return Curse{
		errId:   uuid.New().String(),
		Message: fmt.Sprintf(message, args...),
	}
}

// Plunder returns a new Curse using the message of the
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
// proto curse, replacing any existing cause, then returns
// the symptom.
func (cu Curse) WrapIn[T Curse | ProtoCurse](symptom T) T {
	if pc, ok := any(symptom).(ProtoCurse); ok {
		pc.curse.Cause = cu
		return any(pc).(T)
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

// Is returns true if the target is the same type and has
// the same ID as the receiving error.
func (cu Curse) Is(target error) bool {
	if pc, ok := target.(ProtoCurse); ok {
		return cu.errId == pc.curse.errId
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

// ProtoCurse is a [Curse] with a formattable message
// designed to be used as named exported package errors.
// Use [Proto] to create one. [ProtoCurse.Fmt] should be
// called with the correct formatting arguments when
// returning the error. If the exported error needs no
// formatting then use [Err] instead.
type ProtoCurse struct {
	curse Curse
}

// Proto creates a [ProtoCurse] used as named exported
// errors for comparison. When an error occurs, the
// [ProtoCurse.Fmt] should be called with the correct
// formatting arguments.
//
//	var ErrParseNumericBool = curse.Proto(
//		"Numeric bool must be 0 or 1, given %d",
//	)
//
//	func parseNumericBool(n int) (bool, error) {
//		if n == 0 {
//			return false, nil
//		}
//
//		if n == 1 {
//			return true, nil
//		}
//
//		return false, ErrParseNumericBool.Fmt(n)
//	}
func Proto(message string) ProtoCurse {
	return ProtoCurse{
		curse: Err(message),
	}
}

// Error returns the error message, satisfying Go's error
// interface.
func (pc ProtoCurse) Error() string {
	log.Println("WARNING: Use of unformatted ProtoCurse")
	return pc.curse.Message
}

// Fmt formats the error message, which must be format
// string, using fmt.Sprintf and returns a curse.
//
// This is used with template curses
func (pc ProtoCurse) Fmt(args ...any) Curse {
	cu := pc.curse
	cu.Message = fmt.Sprintf(cu.Message, args...)
	return cu
}

var _ error = Curse{}
var _ error = ProtoCurse{}
