package spellbook

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Bless_1(t *testing.T) {
	act := Bless("abc")
	exp := &effect{
		result: "abc",
	}

	require.Equal(t, exp, act)
}

func Test_Sin_1(t *testing.T) {
	err := errors.New("abc")

	act := Cursed[string](err)
	exp := &effect{
		result: "",
		error:  err,
	}

	require.Equal(t, exp, act)
}

func Test_Curse_1(t *testing.T) {
	act := Curse[string]("Error: %s", "abc")
	exp := &effect{
		result: "",
		error:  errors.New("Error: abc"),
	}

	require.Equal(t, exp, act)
}

func Test_Judge_1(t *testing.T) {
	var act, exp Effect

	act = Judge[string]("abc")
	exp = &effect{
		result: "abc",
	}
	require.Equal(t, exp, act)

	err := errors.New("Error: abc")
	act = Judge[string](err)
	exp = &effect{
		result: "",
		error:  err,
	}
	require.Equal(t, exp, act)
}

func Test_Choose_1(t *testing.T) {
	var act, exp Effect

	act = Choose("abc", nil)
	exp = &effect{
		result: "abc",
	}
	require.Equal(t, exp, act)

	err := errors.New("Error: abc")
	act = Choose("abc", err)
	exp = &effect{
		result: "",
		error:  err,
	}
	require.Equal(t, exp, act)
}

func Test_effect_NameAs_1(t *testing.T) {
	ef := &effect{}

	act := ef.NameAs("Bob")
	exp := &effect{
		name: "Bob",
	}

	require.Equal(t, exp, act)
}

func Test_effect_Dispel_1(t *testing.T) {
	ef := &effect{}

	act := ef.Dispel()
	exp := &effect{
		endSpell: true,
	}

	require.Equal(t, exp, act)
}

func Test_SeekNamedEffect_1(t *testing.T) {
	A := &effect{name: "A", prior: nil}
	B := &effect{name: "B", prior: A}
	C := &effect{name: "C", prior: B}
	D := &effect{name: "D", prior: C}

	act := SeekNamedEffect(D, "B")
	require.Equal(t, B, act)

	act = SeekNamedEffect(B, "D")
	require.Equal(t, nil, act)
}
