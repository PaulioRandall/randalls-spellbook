package sourcery

import (
	"fmt"
)

// Effect is the result of an incantation. It contains the
// value or error and any flags for controlling the next
// incantation in the spell.
type Effect struct {
	endSpell    bool
	ignoreValue bool
	value       any
	error       error
}

// Judge returns an Effect with error as the err value if
// it's not nil, else uses value as the result value.
func Judge(value any, err error) Effect {
	if err != nil {
		return Effect{
			error: err,
		}
	}

	return Effect{
		value: value,
	}
}

// Demystify returns an Effect with valueOrError as the
// error value if it is an error, else uses valueOrError as
// the result value.
func Demystify(valueOrError any) Effect {
	if e, ok := valueOrError.(error); ok {
		return Effect{
			error: e,
		}
	}

	return Effect{
		value: valueOrError,
	}
}

// Bestow returns an Effect using the given value as the
// result value.
func Bestow(value any) Effect {
	return Effect{
		value: value,
	}
}

// Forsake returns an Effect that sets the incantation
// input as the result value, i.e. no new value.
func Forsake() Effect {
	return Effect{
		ignoreValue: true,
	}
}

// Purify returns an Effect that sets nil as the value.
func Purify() Effect {
	return Effect{
		// None.
	}
}

// Sin returns an Effect that sets the error.
func Sin(err error) Effect {
	return Effect{
		error: err,
	}
}

// Curse returns an Effect by creating a new error from
// the given message and optional arguments.
func Curse(message string, args ...any) Effect {
	return Effect{
		error: fmt.Errorf(message, args...),
	}
}

// Dispel returns an Effect that will end the spell with
// nil as the result value.
func Dispel() Effect {
	return Effect{
		endSpell: true,
	}
}

// Dispel returns an Effect that will end the spell using
// the incantation input as the result value.
func Bless() Effect {
	return Effect{
		endSpell:    true,
		ignoreValue: true,
	}
}

// Honor ends the spell with the value as the result.
func Honor(value any) Effect {
	return Effect{
		endSpell: true,
		value:    value,
	}
}

// Cursed returns true if the effect is an error.
func (ef Effect) Cursed() bool {
	return ef.error != nil
}

// Dispelled returns true if the spell ends with this
// incantation.
func (ef Effect) Dispelled() bool {
	return ef.endSpell
}

// Forsaken returns true if the incantation value should
// be ignored and the input data should be used as the
// value.
func (ef Effect) Forsaken() bool {
	return ef.ignoreValue
}

// Flaw returns the error if set, else returns nil.
func (ef Effect) Flaw() error {
	return ef.error
}

// Summon returns the value if set, else returns nil.
func (ef Effect) Summon() any {
	return ef.value
}
