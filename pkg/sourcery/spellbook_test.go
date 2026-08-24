package sourcery

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func mockIncantation(
	data any,
	err error,
) (Incantation, *int) {
	count := 0

	f := func(effect Effect) Effect {
		count++
		return effect.Judge(data, err)
	}

	return f, &count
}

func requireSubset[T any](
	t *testing.T,
	list, subset []T,
) {
	for i, _ := range subset {
		require.Equal(
			t,
			reflect.ValueOf(list[i]).Pointer(),
			reflect.ValueOf(subset[i]).Pointer(),
		)
	}
}

func Test_Scribe(t *testing.T) {
	spellbook := ConjureSpellbook()
	incant, _ := mockIncantation(nil, nil)

	spellbook.Scribe("spellname", incant)

	exp := []Incantation{
		incant,
	}

	requireSubset(t, exp, spellbook.spells["spellname"])
}

func Test_Seek(t *testing.T) {
	spellbook := ConjureSpellbook()
	incantA, _ := mockIncantation("A", nil)
	incantB, _ := mockIncantation("B", nil)

	spellbook.Scribe("spellname", incantA)
	spellbook.Scribe("spellname", incantB)

	act := spellbook.Seek("spellname")
	exp := Spell{
		incantA,
		incantB,
	}

	requireSubset(t, exp, act)
}

func Test_Cast_1(t *testing.T) {
	var effect Effect
	spellbook := ConjureSpellbook()
	incant, countPtr := mockIncantation("data", nil)

	spellbook.Scribe("spellname", incant)
	effect = spellbook.Cast("spellname", nil)
	require.Equal(t, false, effect.Cursed())
	require.Equal(t, "data", effect.Summon())

	require.Equal(t, 1, *countPtr)
}

func Test_Cast_2(t *testing.T) {
	var effect Effect
	spellbook := ConjureSpellbook()
	incant, countPtr := mockIncantation("data", nil)

	spellbook.Scribe("spellname", incant)
	spellbook.Scribe("spellname", incant)
	spellbook.Scribe("spellname", incant)

	effect = spellbook.Cast("spellname", nil)
	require.Equal(t, false, effect.Cursed())
	require.Equal(t, "data", effect.Summon())

	require.Equal(t, 3, *countPtr)
}

func Test_Cast_3(t *testing.T) {
	var effect Effect
	spellbook := ConjureSpellbook()

	incant, countPtr := mockIncantation("data", nil)
	spellbook.Scribe("spellname", incant)

	effect = spellbook.Cast("spellname", nil)
	require.Equal(t, false, effect.Cursed())
	require.Equal(t, "data", effect.Summon())

	effect = spellbook.Cast("spellname", nil)
	require.Equal(t, false, effect.Cursed())
	require.Equal(t, "data", effect.Summon())

	effect = spellbook.Cast("spellname", nil)
	require.Equal(t, false, effect.Cursed())
	require.Equal(t, "data", effect.Summon())

	require.Equal(t, 3, *countPtr)
}

func Test_Cast_4(t *testing.T) {
	var effect Effect
	spellbook := ConjureSpellbook()

	err := errors.New("error")
	incant, _ := mockIncantation(nil, err)
	spellbook.Scribe("spellname", incant)

	effect = spellbook.Cast("spellname", nil)
	require.Equal(t, true, effect.Cursed())
	require.Equal(t, err, effect.Flaw())
	require.Equal(t, nil, effect.Summon())
}
