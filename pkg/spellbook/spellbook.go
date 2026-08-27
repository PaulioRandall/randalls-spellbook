package spellbook

import (
	"encoding/json"
)

// Incantation is a single step when casting a spell. It
// may transform data, create a side effect (e.g. write to
// storage), check data and generate an error, or a
// combination of these.
//
// The first thing most incantations will do is parse and
// check the input data. This is the price paid for such a
// decoupled architecture. Of course you can attempt to
// cast values to new types without checking to save a few
// lines, but the runtime errors are on you for that.
//
// In most cases, the output Effect should be constructed
// using the input Effect's functions to ensure future
// incantations and the user of a spell's result can
// inspect the chain of Effects. This is not enforced or
// required in the casting of spells as there may be good
// reason to hide or force garbage collection of prior
// Effects. However, having the chain of Effects produced
// is useful for testing spells and debugging!
type Incantation func(ctx any, prior Effect) Effect

// Spell is a series of incantations that are to be invoked
// in order.
type Spell []Incantation

// Spellbook stores and enables casting of spells.
type Spellbook struct {
	spells map[string]Spell
}

// JsonDemystifyer returns an Incantation that unmarshalls
// the input Effect's JSON byte array data into object,
// which is returned in the resultant effect. The object
// may be a concrete value or a pointer to one. The
// incantation sets error if the input data is not a byte
// array or if unmarshalling fails.
func JsonDemystifyer(object any) Incantation {
	return func(_ any, input Effect) Effect {
		bytes, ok := input.Result().([]byte)
		if !ok {
			return Curse("Wrong type, expected []bytes")
		}

		e := json.Unmarshal(bytes, object)
		return Choose(object, e)
	}
}

// Conjure a new Spellbook.
func Conjure() *Spellbook {
	return &Spellbook{
		spells: map[string]Spell{},
	}
}

// Enscribe adds one or more new incantations to the end of
// a spell. Incantations may appear multiple times within
// the same spell.
func (spellbook *Spellbook) Enscribe(
	spellname string,
	incantations ...Incantation,
) {
	spellbook.spells[spellname] = append(
		spellbook.spells[spellname],
		incantations...,
	)
}

// Transcribe adds or overrides a spell.
func (spellbook *Spellbook) Transcribe(
	spellname string,
	spell ...Incantation,
) {
	spellbook.spells[spellname] = spell
}

// Describe finds and returns a Spell, or nil if no
// incantations for spellname exist.
func (spellbook *Spellbook) Describe(
	spellname string,
) Spell {
	for name, spell := range spellbook.spells {
		if name == spellname {
			return spell
		}
	}

	return nil
}

// Cast a spell invoking its incantations using the passed
// data as the initial input. If the spell has no
// incantations an effect is returned containing the input
// as the result data.
func (spellbook *Spellbook) Cast(
	spellname string,
	ctx any,
	data any,
) Effect {
	spell := spellbook.Describe(spellname)
	var input Effect = Bless(data)
	var output Effect = input

	for _, incant := range spell {
		output = incant(ctx, input).PriorAs(input)

		if output.Cursed() || output.Dispelled() {
			return output
		}

		input = output
	}

	return output
}
