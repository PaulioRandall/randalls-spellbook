package spellbook

import (
	"errors"
	"fmt"
)

func effectExample() {
	var err error = errors.New("error message")
	var ef, prior Effect

	ef = Bless("data")
	ef = Curse[string]("error %s", "message")
	ef = Cursed[string](err)
	ef = Judge[string]("data")
	ef = Judge[string](err)
	ef = Choose("data", err)

	prior = Bless("abc")
	ef = Bless("efg").
		PriorAs(prior).
		NameAs("3rd, 4th, and 5th letters").
		Dispel()

	_ = ef.Name()      // "3rd, 4th, and 5th letters"
	_ = ef.Prior()     // Effect{result: "abc"}
	_ = ef.Result()    // "efg"
	_ = ef.Cursed()    // false
	_ = ef.Dispelled() // true
	_, _ = ef.Values() // "efg", nil

	ef = Curse[string]("error %s", "message")

	_ = ef.Cursed()    // true
	_ = ef.Error()     // error{message:"error message"}
	_, _ = ef.Values() // nil, error{message:"error message"}
}

func spellbookExample() {
	spellbook := Conjure()

	spellbook.Enscribe("fireball", newFireball)
	spellbook.Enscribe("fireball", inflictDamage)

	spellbook.Cast("fireball", nil, 123)
}

type Fireball struct {
	baseDamage int
}

func newFireball(_ any, data any) Effect {
	dmg := Demystify[int](data)

	fb := Fireball{
		baseDamage: dmg,
	}

	return Bless(fb).NameAs("Fireball")
}

func inflictDamage(_ any, data any) Effect {
	fb := Demystifyf[Fireball](
		data,
		"Wrong type, expected a Fireball",
	)

	fmt.Printf(
		"Fireball inflicts %d fire damage.\n",
		fb.baseDamage,
	)

	return Bless(fb)
}
