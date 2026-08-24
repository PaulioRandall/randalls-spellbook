package sourcery

import (
	"fmt"
)

func Example() {
	spellbook := ConjureSpellbook()

	spellbook.Scribe("fireball", newFireball)
	spellbook.Scribe("fireball", inflictDamage)

	spellbook.Cast("fireball", 123)
}

type Fireball struct {
	baseDamage int
}

func newFireball(input Effect) Effect {
	dmg, ok := input.Summon().(int)
	if !ok {
		return input.Curse("Wrong type, expected an int")
	}

	fb := Fireball{
		baseDamage: dmg,
	}

	return input.Bestow(fb)
}

func inflictDamage(input Effect) Effect {
	fb, ok := input.Summon().(Fireball)
	if !ok {
		return input.Curse("Wrong type, expected a Fireball")
	}

	fmt.Printf(
		"Fireball inflicts %d fire damage.\n",
		fb.baseDamage,
	)

	return input.Forsake()
}
