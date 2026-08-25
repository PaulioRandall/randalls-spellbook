package spellbook

import (
	"fmt"
)

func SpellbookExample() {
	spellbook := Conjure()

	spellbook.Enscribe("fireball", newFireball)
	spellbook.Enscribe("fireball", inflictDamage)

	spellbook.Cast("fireball", 123)
}

type Fireball struct {
	baseDamage int
}

func newFireball(input Effect) Effect {
	dmg, ok := input.Result().(int)
	if !ok {
		return Curse("Wrong type, expected an int")
	}

	fb := Fireball{
		baseDamage: dmg,
	}

	return Bless(fb).Named("Fireball")
}

func inflictDamage(input Effect) Effect {
	fb, ok := input.Result().(Fireball)
	if !ok {
		return Curse("Wrong type, expected a Fireball")
	}

	fmt.Printf(
		"Fireball inflicts %d fire damage.\n",
		fb.baseDamage,
	)

	return Bless(fb)
}
