package sourcery

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Effect is the result of an incantation. It contains the
// data or error and a flag for ending the spell, skipping
// the remaining incantations. It also contains a reference
// to the previous effect which will only be nil for the
// initial effect (i.e. it's a linked list).
type Effect struct {
	prior    *Effect
	data     any
	error    error
	endSpell bool
}

// JsonIllumator returns an Incantation that unmarshalls
// the input Effect's JSON byte array data into object,
// which is returned in the resultant effect. The object
// may be a concrete value or a pointer to one. The
// incantation sets error if the input data is not a byte
// array or if unmarshalling fails.
func JsonIllumator(object any) Incantation {
	return func(input Effect) Effect {
		bytes, ok := input.Summon().([]byte)
		if !ok {
			return input.Curse("Wrong type, expected []bytes")
		}

		var e error

		if reflect.ValueOf(object).Kind() == reflect.Ptr {
			e = json.Unmarshal(bytes, object)
		} else {
			e = json.Unmarshal(bytes, &object)
		}

		return input.Judge(object, e)
	}
}

// Judge returns a new Effect with err as the error value
// if it's not nil, else uses data as the result data.
func (effect *Effect) Judge(data any, err error) Effect {
	if err != nil {
		return Effect{
			prior: effect,
			error: err,
		}
	}

	return Effect{
		prior: effect,
		data:  data,
	}
}

// Demystify returns a new Effect with dataOrError as the
// error value if it's an error, else uses dataOrError as
// the result data.
func (effect *Effect) Demystify(dataOrError any) Effect {
	if e, ok := dataOrError.(error); ok {
		return Effect{
			prior: effect,
			error: e,
		}
	}

	return Effect{
		prior: effect,
		data:  dataOrError,
	}
}

// Bestow returns a new Effect using the given data as the
// result data.
func (effect *Effect) Bestow(data any) Effect {
	return Effect{
		prior: effect,
		data:  data,
	}
}

// Forsake returns a new Effect that uses the prior
// Effect's data as the result data.
func (effect *Effect) Forsake() Effect {
	return Effect{
		prior: effect,
		data:  effect.Summon(),
	}
}

// Purify returns a new Effect that sets the data as nil.
func (effect *Effect) Purify() Effect {
	return Effect{
		prior: effect,
	}
}

// Sin returns a new Effect that sets the error to err.
func (effect *Effect) Sin(err error) Effect {
	return Effect{
		prior: effect,
		error: err,
	}
}

// Curse returns a new Effect by creating a new error from
// the given message and optional arguments.
func (effect *Effect) Curse(
	message string,
	args ...any,
) Effect {
	return Effect{
		prior: effect,
		error: fmt.Errorf(message, args...),
	}
}

// Dispel returns a new Effect that will end the spell with
// nil as the result value.
func (effect *Effect) Dispel() Effect {
	return Effect{
		prior:    effect,
		endSpell: true,
	}
}

// Dispel returns a new Effect that ends the spell and uses
// the prior Effect's data as the result data.
func (effect *Effect) Bless() Effect {
	return Effect{
		prior:    effect,
		data:     effect.Summon(),
		endSpell: true,
	}
}

// Honor returns a new Effect that ends the spell with the
// passed data as the result data.
func (effect *Effect) Honor(data any) Effect {
	return Effect{
		prior:    effect,
		endSpell: true,
		data:     data,
	}
}

// Prior returns the prior effect.
func (effect *Effect) Prior() Effect {
	return *(effect.prior)
}

// Cursed returns true if the effect is an error.
func (effect *Effect) Cursed() bool {
	return effect.error != nil
}

// Flaw returns the error if set, else returns nil.
func (effect *Effect) Flaw() error {
	return effect.error
}

// Summon returns the value if set, else returns nil.
func (effect *Effect) Summon() any {
	return effect.data
}

// Dispelled returns true if the spell ends with this
// incantation.
func (effect *Effect) Dispelled() bool {
	return effect.endSpell
}
