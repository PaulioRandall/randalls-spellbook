package spellbook

// Spellbook stores and enables casting of spells.
type Spellbook struct {
	spells map[string]Spell
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
		output = incant(ctx, input.Result()).PriorAs(input)

		if output.Cursed() || output.Dispelled() {
			return output
		}

		input = output
	}

	return output
}
