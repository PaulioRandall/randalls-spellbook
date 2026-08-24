package sourcery

import (
	"encoding/json"
)

type Fireball struct {
	baseDamage   int
	burnInterval int
	burnCount    int
	burnDamage   int
}

type Frostbolt struct {
	baseDamage   int
	freezeChance int
}

func Example() {
	spellbook := New()

	spellbook.Scribe("fireball", onFireball)
	spellbook.Scribe("frostbolt", onFrostbolt)
	spellbook.Scribe("thunderbolt", onThunderbolt)

	spellbook.Cast("fireball", nil)
	spellbook.Cast("frostbolt", onFrostbolt)
	spellbook.Cast("thunderbolt", nil)
}

func demystifyFrostbolt(bytes []byte) Effect {
	frostbolt := Frostbolt{}

	e := json.Unmarshal(bytes, &frostbolt)
	if e != nil {
		return Sin(e)
	}

	return Bestow(frostbolt)
}

func onFireball(data any) Effect {
	fireball, ok := data.(Fireball)
	if !ok {
		return Curse("Wrong type, expected Fireball")
	}

	println("Fire damage")

	return Bestow(fireball)
}

func onFrostbolt(data any) Effect {
	frostbolt, ok := data.(Frostbolt)
	if !ok {
		return Curse("Wrong type, expected Frostbolt")
	}

	println("Ice damage!")

	return Bestow(frostbolt)
}

func onThunderbolt(_ any) Effect {
	println("Lighting damage! (no args)")

	return Purify()
}
