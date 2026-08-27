package spellbook

import (
	"fmt"
)

// Effect is the result of an incantation. It contains some
// value or an error. Each effect has a name, but it can be
// empty. It also contains a reference to the previous
// effect which will only be nil for the initial effect
// (yes, it's a linked list).
type Effect interface {
	// Name returns the effect name. It may be empty.
	Name() string

	// NameAs copies and returns the effect with a new name.
	NameAs(name string) Effect

	// Prior returns the prior effect.
	Prior() Effect

	// PriorAs copies and returns the effect with the prior
	// set as the passed prior Effect.
	PriorAs(prior Effect) Effect

	// Result returns the result value if set, else returns
	// nil.
	Result() any

	// Error returns the error if set, else returns nil.
	Error() error

	// Cursed returns true if the Effect contains an error.
	Cursed() bool

	// Dispelled returns true if the spell ends with this
	// incantation.
	Dispelled() bool

	// Dispel copies and returns the effect with the end
	// spell flag set as true.
	Dispel() Effect

	// Values returns the result value and a nil error if no
	// error is set, else returns a nil value and the error.
	Values() (any, error)
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
func newEffect[T any](value T, err error) *effect {
	ef := &effect{}

	if err != nil {
		var r T
		ef.result = r
		ef.error = err
	} else {
		ef.result = value
	}

	return ef
}

// Bless returns a new effect with the passed value set as
// the result.
func Bless[T any](value T) *effect {
	return &effect{
		result: value,
	}
}

// Bequeath returns a new effect with the passed ef result
// used as the new result.
func Bequeath(ef *effect) *effect {
	return &effect{
		result: ef.Result(),
	}
}

// Curse returns a new effect with a new error created from
// the given message and optional arguments.
func Curse[T any](message string, args ...any) *effect {
	var r T
	return &effect{
		result: r,
		error:  fmt.Errorf(message, args...),
	}
}

// Cursed returns a new effect with a the passed err set
// as the error.
func Cursed[T any](err error) *effect {
	var r T
	return &effect{
		result: r,
		error:  err,
	}
}

// Judge returns a new effect with the error value set
// if valueOrErr is an error, else sets the result value to
// valueOrErr.
func Judge[T any](valueOrErr any) *effect {
	r, _ := valueOrErr.(T)
	err, _ := valueOrErr.(error)
	return newEffect(r, err)
}

// Choose returns a new effect with the error set if err is
// not nil, else the result value is set.
func Choose[T any](value T, err error) *effect {
	return newEffect(value, err)
}

// Name satisfies the Effect interface.
func (ef *effect) Name() string {
	return ef.name
}

// NameAs satisfies the Effect interface.
func (ef *effect) NameAs(name string) Effect {
	cp := *ef
	cp.name = name
	return &cp
}

// Prior satisfies the Effect interface.
func (ef *effect) Prior() Effect {
	return ef.prior
}

// PriorAs satisfies the Effect interface.
func (ef *effect) PriorAs(prior Effect) Effect {
	cp := *ef
	cp.prior = prior
	return &cp
}

// Result satisfies the Effect interface.
func (ef *effect) Result() any {
	return ef.result
}

// Error satisfies the Effect interface.
func (ef *effect) Error() error {
	return ef.error
}

// Cursed satisfies the Effect interface.
func (ef *effect) Cursed() bool {
	return ef.error != nil
}

// Dispelled satisfies the Effect interface.
func (ef *effect) Dispelled() bool {
	return ef.endSpell
}

// Dispel satisfies the Effect interface.
func (ef *effect) Dispel() Effect {
	cp := *ef
	cp.endSpell = true
	return &cp
}

// Values satisfies the Effect interface.
func (ef *effect) Values() (any, error) {
	if ef.error != nil {
		return nil, ef.error
	}
	return ef.result, nil
}

// DemystifyEffect returns the result value, as the return
// type, and a nil error if no error is set, else returns
// the return types empty value and the error. If the
// result value cannot be cast to the generic type then a
// panic ensues.
func DemystifyEffect[R any](ef Effect) (R, error) {
	result, ok := ef.Result().(R)
	if !ok {
		panic("Not even a Wizard² can demystify the Effect to your uttered type")
	}

	if ef.Cursed() {
		var empty R
		return empty, ef.Error()
	}

	return result, nil
}

// SeekNamedEffect iterates the chain of Effects, starting
// with the passed argument ef, and returns the first
// found with the specified name, or nil if no matching
// Effect is found.
func SeekNamedEffect(ef Effect, name string) Effect {
	prior := ef

	for prior != nil {
		if prior.Name() == name {
			return prior
		}

		prior = prior.Prior()
	}

	return nil
}
