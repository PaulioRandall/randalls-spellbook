package spellbook

import (
	"encoding/json"
	"errors"
	"reflect"
)

type DemystifyData func([]byte) (any, error)

type Incantation func(data any) (any, error)

type Spellbook struct {
	spells map[string][]Incantation
}

func New() *Spellbook {
	return &Spellbook{
		spells: map[string][]Incantation{},
	}
}

func (spellbook *Spellbook) ScribeIncantation(
	spellname string,
	incantation Incantation,
) {
	spellbook.spells[spellname] = append(
		spellbook.spells[spellname],
		incantation,
	)
}

func (spellbook *Spellbook) SeekIncantations(
	spellname string,
) []Incantation {
	for name, incantations := range spellbook.spells {
		if name == spellname {
			return incantations
		}
	}

	return nil
}

func (spellbook *Spellbook) CastSpell(
	spellname string,
	data any,
) (any, error) {
	incantations := spellbook.SeekIncantations(spellname)

	if incantations == nil {
		return nil, errors.New("The spell bears no incantations within this spellbook.")
	}

	return chant(incantations, data)
}

func chant(
	incantations []Incantation,
	data any,
) (any, error) {
	var e error

	for _, incant := range incantations {
		data, e = incant(data)
		if e != nil {
			return nil, e
		}
	}

	return data, nil
}

func DemystifyToNil(_ []byte) (any, error) {
	return nil, nil
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
	return func(bytes []byte) (any, error) {
		return demystifyToObject(bytes, demystifyer)
	}
}

func demystifyToObject(
	bytes []byte,
	objectExample any,
) (any, error) {
	objectType := reflect.TypeOf(objectExample)
	object := reflect.Zero(objectType).Interface()

	e := json.Unmarshal(bytes, &object)
	if e != nil {
		return nil, e
	}

	return object, nil
}

func Invoke(spell string) {
	// TODO
}
