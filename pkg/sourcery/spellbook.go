package sourcery

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
type Incantation func(Effect) Effect

// Spell is a series of incantations that are to be invoked
// in order.
type Spell []Incantation

// Spellbook stores and enables casting of spells.
type Spellbook struct {
	spells map[string]Spell
}

// ConjureSpellbook a new Spellbook.
func ConjureSpellbook() *Spellbook {
	return &Spellbook{
		spells: map[string]Spell{},
	}
}

// Scribe adds a new incantation to the end of a spell.
// Incantations may appear multiple times within a spell.
func (spellbook *Spellbook) Scribe(
	spellname string,
	incantation Incantation,
) {
	spellbook.spells[spellname] = append(
		spellbook.spells[spellname],
		incantation,
	)
}

// Seek finds and returns a Spell, or nil if no
// incantations for spellname exist.
func (spellbook *Spellbook) Seek(
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
	data any,
) Effect {
	spell := spellbook.Seek(spellname)
	return castSpell(spell, data)
}

func castSpell(
	spell Spell,
	data any,
) Effect {
	input := Effect{data: data}
	output := input

	for _, incant := range spell {
		output = incant(input)

		if output.Cursed() {
			return output
		}

		if output.Dispelled() {
			break
		}

		input = output
	}

	return output
}
