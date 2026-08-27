package effect

import (
	"fmt"
)

// UntypedEffect represents an untyped result or an error,
// similar to Rust's Result enum without the types. Effect
// is the typed version of UntypedEffect.
//
// It contains a reference to the previous effect so
// effects can be chained to create a [linked list] of
// prior effects representing a chain of operations.
type UntypedEffect interface {
	// Name returns the effect name. It may be empty.
	Name() string

	// Prior returns the prior effect. Note that the prior
	// effect may carry a different result type than this
	// one; use TransmuteEffect to recover it in typed form.
	Prior() UntypedEffect

	// UntypedResult returns the result as 'any' type since
	// we don't know what the type is unless cast to an
	// Effect.
	UntypedResult() any

	// Error returns the error if set, else returns nil.
	Error() error

	// Cursed returns true if the Effect contains an error.
	Cursed() bool

	// UntypedValues returns the result value and a nil error
	// if no error is set, else returns nil and the error.
	UntypedValues() (any, error)
}

// Effect is the typed version of UntypedEffect. Functions
// here will return an Effect unless we cannot determine
// the type at compile time, UntypedEffect will be used
// instead. This maximises use that minimises direct type
// casting.
type Effect[T any] interface {
	UntypedEffect

	// NameAs copies and returns the effect with a new name.
	NameAs(name string) Effect[T]

	// PriorAs copies and returns the effect with the prior
	// set as the passed prior Effect. The prior may carry a
	// different result type.
	PriorAs(prior UntypedEffect) Effect[T]

	// Result returns the result value if set, else returns
	// the zero value of T.
	Result() T

	// Values returns the result value and a nil error if no
	// error is set, else returns the zero value of T and the
	// error.
	Values() (T, error)
}

// effect satisfies the TypedEffect interface.
type effect[T any] struct {
	prior  UntypedEffect
	name   string
	result T
	error  error
}

// newEffect sets the effect's error if err is not
// nil, else sets the result value.
func newEffect[T any](value T, err error) *effect[T] {
	ef := &effect[T]{}

	if err != nil {
		ef.error = err
	} else {
		ef.result = value
	}

	return ef
}

// Bless returns a new effect with the passed value set as
// the result. If the value is itself an Effect[T] then
// its result will be used instead. A result value cannot
// be an Effect or UntypedEffect.
func Bless[T any](value T) *effect[T] {
	if te, ok := any(value).(Effect[T]); ok {
		value = te.Result()
	}

	return &effect[T]{
		result: value,
	}
}

// Curse returns a new effect with a new error created from
// the given message and optional arguments. The result
// type T must be specified explicitly since no value is
// provided, e.g. Curse[string]("alakazam").
func Curse[T any](message string, args ...any) *effect[T] {
	return &effect[T]{
		error: fmt.Errorf(message, args...),
	}
}

// Cursed returns a new effect with the passed err set
// as the error. The result type T must be specified
// explicitly since no value is provided, e.g.
// Cursed[string](err).
func Cursed[T any](err error) *effect[T] {
	return &effect[T]{
		error: err,
	}
}

// Judge returns a new effect with the error value set
// if valueOrErr is an error, else sets the result value to
// valueOrErr cast to T.
func Judge[T any](valueOrErr any) *effect[T] {
	if err, ok := valueOrErr.(error); ok {
		return &effect[T]{error: err}
	}

	value, _ := valueOrErr.(T)
	return &effect[T]{result: value}
}

// Choose returns a new effect with the error set if err is
// not nil, else the result value is set.
func Choose[T any](value T, err error) *effect[T] {
	return newEffect(value, err)
}

// Name satisfies the UntypedEffect and Effect interfaces.
func (ef *effect[T]) Name() string {
	return ef.name
}

// NameAs satisfies the UntypedEffect and Effect interfaces.
func (ef *effect[T]) NameAs(name string) Effect[T] {
	cp := *ef
	cp.name = name
	return &cp
}

// Prior satisfies the UntypedEffect and Effect interfaces.
func (ef *effect[T]) Prior() UntypedEffect {
	return ef.prior
}

// PriorAs satisfies the UntypedEffect and Effect
// interfaces.
func (ef *effect[T]) PriorAs(prior UntypedEffect) Effect[T] {
	cp := *ef
	cp.prior = prior
	return &cp
}

// Result satisfies the UntypedEffect and Effect
// interfaces.
func (ef *effect[T]) UntypedResult() any {
	return ef.result
}

// Result satisfies the UntypedEffect and Effect
// interfaces.
func (ef *effect[T]) Result() T {
	return ef.result
}

// Error satisfies the UntypedEffect and Effect
// interfaces.
func (ef *effect[T]) Error() error {
	return ef.error
}

// Cursed satisfies the UntypedEffect and Effect
// interfaces.
func (ef *effect[T]) Cursed() bool {
	return ef.error != nil
}

// UntypedValues satisfies the TypedEffect and Effect
// interfaces.
func (ef *effect[T]) UntypedValues() (any, error) {
	if ef.error != nil {
		return nil, ef.error
	}
	return ef.result, nil
}

// Values satisfies the Effect interface.
func (ef *effect[T]) Values() (T, error) {
	if ef.error != nil {
		var empty T
		return empty, ef.error
	}
	return ef.result, nil
}

// Demystify casts the UntypedEffect ue to Effect[T].
// If the underlying result type isn't T then the resultant
// Effect shall contain an error stating so, with ue set as
// the prior Effect.
func Demystify[T any](ue UntypedEffect) Effect[T] {
	te, ok := ue.(Effect[T])
	if !ok {
		result := Curse[T]("Not even a Wizard² can transmute the Effect result to your uttered type")
		return result.PriorAs(ue)
	}

	return te
}

// SeekNamedEffect iterates the chain of UntypedEffects,
// starting with the passed argument ue, and returns the
// first found with the specified name, or nil if no
// matching UntypedEffect is found. Since links in the
// chain may carry different result types, the return value
// is UntypedEffect. If you're confident you know what the
// type is, use Demystify to cast to an Effect.
func SeekNamedEffect(
	ue UntypedEffect,
	name string,
) UntypedEffect {
	prior := ue

	for prior != nil {
		if prior.Name() == name {
			return prior
		}

		prior = prior.Prior()
	}

	return nil
}
