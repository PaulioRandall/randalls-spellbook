package sourcery

import (
	"encoding/json"
	"reflect"
)

type DemystifyData func(bytes []byte) Effect
type Incantation func(data any) Effect

type Spellbook struct {
	spells map[string][]Incantation
}

func New() *Spellbook {
	return &Spellbook{
		spells: map[string][]Incantation{},
	}
}

func (spellbook *Spellbook) Scribe(
	spellname string,
	incantation Incantation,
) {
	spellbook.spells[spellname] = append(
		spellbook.spells[spellname],
		incantation,
	)
}

func (spellbook *Spellbook) Seek(
	spellname string,
) []Incantation {
	for name, incantations := range spellbook.spells {
		if name == spellname {
			return incantations
		}
	}

	return nil
}

func (spellbook *Spellbook) Cast(
	spellname string,
	data any,
) Effect {
	incantations := spellbook.Seek(spellname)

	if incantations == nil {
		return Curse("The spell bears no incantations within this spellbook.")
	}

	return chant(incantations, data)
}

func chant(
	incantations []Incantation,
	data any,
) Effect {
	var effect Effect

	for _, incant := range incantations {
		effect = incant(data)
		if effect.Cursed() {
			return effect
		}
	}

	return effect
}

func DemystifyToNil(_ []byte) Effect {
	return Purify()
}

func parseDemystifyer(demystifyer any) DemystifyData {
	if demystifyer == nil {
		return DemystifyToNil
	}

	f, ok := demystifyer.(DemystifyData)
	if ok {
		return f
	}

	// Assume it's an object from which we can determine the
	// type and instantiate a new instance of.
	return func(bytes []byte) Effect {
		return demystifyToObject(bytes, demystifyer)
	}
}

func demystifyToObject(
	bytes []byte,
	objectExample any,
) Effect {
	objectType := reflect.TypeOf(objectExample)
	object := reflect.Zero(objectType).Interface()

	e := json.Unmarshal(bytes, &object)
	if e != nil {
		return Sin(e)
	}

	return Bestow(object)
}
