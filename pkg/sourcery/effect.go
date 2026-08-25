package sourcery

import (
	"fmt"
)

// Effect is the result of an incantation. It contains some
// value or an error. There is a flag for ending the spell,
// i.e. skipping the remaining incantations. Each effect
// has a name, but it can be empty. It also contains a
// reference to the previous effect which will only be nil
// for the initial effect (yes, it's a linked list).
type Effect interface {
	// Name returns the effect name. It may be empty.
	Name() string

	// Prior returns the prior effect.
	Prior() Effect

	// Result returns the result value if set, else returns
	// nil.
	Result() any

	// Curse returns the error if set, else returns nil.
	Curse() error

	// IsCursed returns true if the effect is an error.
	IsCursed() bool

	// IsDispelled returns true if the spell ends with this
	// incantation.
	IsDispelled() bool
}

// effect satisfies the Effect interface.
type effect struct {
	prior    Effect
	name     string
	result   any
	error    error
	endSpell bool
}

// newEffect sets the effect's error if err is not
// nil, else sets the result value.
func newEffect(value any, err error) effect {
	ef := effect{}

	if err != nil {
		ef.error = err
	} else {
		ef.result = value
	}

	return ef
}

// Bless returns a new effect with the passed value set as
// the result. If the value is an Effect then its value
// will be used instead. A result value cannot be an
// Effect.
func Bless(value any) effect {
	if ef, ok := value.(Effect); ok {
		value = ef.Result()
	}

	return effect{
		result: value,
	}
}

// Ruin returns a new effect with a the passed err set as
// the error.
func Ruin(err error) effect {
	return effect{
		error: err,
	}
}

// Curse returns a new effect with a new error created from
// the given message and optional arguments.
func Curse(message string, args ...any) effect {
	return effect{
		error: fmt.Errorf(message, args...),
	}
}

// Demystify returns a new effect with the error value set
// if valueOrErr is an error, else sets the result value to
// valueOrErr.
func Demystify(valueOrErr any) effect {
	err, _ := valueOrErr.(error)
	return newEffect(valueOrErr, err)
}

// Judge returns a new effect with the error set if err is
// not nil, else the result value is set.
func Judge(value any, err error) effect {
	return newEffect(value, err)
}

// Name satisfies the Effect interface.
func (ef effect) Name() string {
	return ef.name
}

// Prior satisfies the Effect interface.
func (ef effect) Prior() Effect {
	return ef.prior
}

// Result satisfies the Effect interface.
func (ef effect) Result() any {
	return ef.result
}

// Curse satisfies the Effect interface.
func (ef effect) Curse() error {
	return ef.error
}

// IsCursed satisfies the Effect interface.
func (ef effect) IsCursed() bool {
	return ef.error != nil
}

// IsDispelled satisfies the Effect interface.
func (ef effect) IsDispelled() bool {
	return ef.endSpell
}

// Named copies and returns the effect with a new name.
func (ef effect) Named(name string) effect {
	ef.name = name
	return ef
}

// Dispelled copies and returns the effect with the end
// spell flag set as true.
func (ef effect) Dispelled(valueOrErr any) effect {
	ef.endSpell = true
	return ef
}
