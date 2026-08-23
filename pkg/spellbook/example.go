package spellbook

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

	spellbook.ScribeIncantation("fireball", onFireball)
	spellbook.ScribeIncantation("frostbolt", onFrostbolt)
	spellbook.ScribeIncantation("thunderbolt", onThunderbolt)

	spellbook.CastSpell("fireball", nil)
	spellbook.CastSpell("frostbolt", onFrostbolt)
	spellbook.CastSpell("thunderbolt", nil)
}

func demystifyFrostbolt(bytes []byte) (any, error) {
	frostbolt := Frostbolt{}

	e := json.Unmarshal(bytes, &frostbolt)
	if e != nil {
		return nil, e
	}

	return frostbolt, nil
}

func onFireball(data any) (any, error) {
	fireball, ok := data.(Fireball)
	if !ok {
		panic("Wrong type, expected Fireball")
	}

	println("Fire damage")

	return fireball, nil
}

func onFrostbolt(data any) (any, error) {
	frostbolt, ok := data.(Frostbolt)
	if !ok {
		panic("Wrong type, expected Frostbolt")
	}

	println("Ice damage!")

	return frostbolt, nil
}

func onThunderbolt(_ any) (any, error) {
	println("Lighting damage! (no args)")

	return nil, nil
}
