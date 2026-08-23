package spellbook

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func mockIncantation(
	result any,
	err error,
) (Incantation, *int) {
	count := 0

	f := func(any) (any, error) {
		count++
		return result, err
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

func Test_ScribeIncantation(t *testing.T) {
	spellbook := New()
	incant, _ := mockIncantation(nil, nil)

	spellbook.ScribeIncantation("spellname", incant)

	exp := []Incantation{
		incant,
	}

	requireSubset(t, exp, spellbook.spells["spellname"])
}

func Test_SeekIncantation(t *testing.T) {
	spellbook := New()
	incantA, _ := mockIncantation("A", nil)
	incantB, _ := mockIncantation("B", nil)

	spellbook.ScribeIncantation("spellname", incantA)
	spellbook.ScribeIncantation("spellname", incantB)

	act := spellbook.SeekIncantations("spellname")
	exp := []Incantation{
		incantA,
		incantB,
	}

	requireSubset(t, exp, act)
}

func Test_CastSpell_1(t *testing.T) {
	var e error
	var result any
	spellbook := New()
	incant, countPtr := mockIncantation("data", nil)

	spellbook.ScribeIncantation("spellname", incant)
	result, e = spellbook.CastSpell("spellname", nil)
	require.Equal(t, error(nil), e)
	require.Equal(t, "data", result)

	require.Equal(t, 1, *countPtr)
}

func Test_CastSpell_2(t *testing.T) {
	var e error
	var result any
	spellbook := New()
	incant, countPtr := mockIncantation("data", nil)

	spellbook.ScribeIncantation("spellname", incant)
	spellbook.ScribeIncantation("spellname", incant)
	spellbook.ScribeIncantation("spellname", incant)

	result, e = spellbook.CastSpell("spellname", nil)
	require.Equal(t, error(nil), e)
	require.Equal(t, "data", result)

	require.Equal(t, 3, *countPtr)
}

func Test_CastSpell_3(t *testing.T) {
	var e error
	var result any
	spellbook := New()

	incant, countPtr := mockIncantation("data", nil)
	spellbook.ScribeIncantation("spellname", incant)

	result, e = spellbook.CastSpell("spellname", nil)
	require.Equal(t, error(nil), e)
	require.Equal(t, "data", result)

	result, e = spellbook.CastSpell("spellname", nil)
	require.Equal(t, error(nil), e)
	require.Equal(t, "data", result)

	result, e = spellbook.CastSpell("spellname", nil)
	require.Equal(t, error(nil), e)
	require.Equal(t, "data", result)

	require.Equal(t, 3, *countPtr)
}

func Test_CastSpell_4(t *testing.T) {
	var e error
	var result any
	spellbook := New()

	err := errors.New("error")
	incant, _ := mockIncantation(nil, err)
	spellbook.ScribeIncantation("spellname", incant)

	result, e = spellbook.CastSpell("spellname", nil)
	require.Equal(t, err, e)
	require.Equal(t, nil, result)
}
