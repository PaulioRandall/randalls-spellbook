package spellbook

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Bless_1(t *testing.T) {
	act := Bless("abc")
	exp := effect{
		result: "abc",
	}

	require.Equal(t, exp, act)
}

func Test_Sin_1(t *testing.T) {
	err := errors.New("abc")

	act := Sin(err)
	exp := effect{
		error: err,
	}

	require.Equal(t, exp, act)
}

func Test_Curse_1(t *testing.T) {
	act := Curse("Error: %s", "abc")
	exp := effect{
		error: errors.New("Error: abc"),
	}

	require.Equal(t, exp, act)
}

func Test_Judge_1(t *testing.T) {
	var act, exp Effect

	act = Judge("abc")
	exp = effect{
		result: "abc",
	}
	require.Equal(t, exp, act)

	err := errors.New("Error: abc")
	act = Judge(err)
	exp = effect{
		error: err,
	}
	require.Equal(t, exp, act)
}

func Test_Choose_1(t *testing.T) {
	var act, exp Effect

	act = Choose("abc", nil)
	exp = effect{
		result: "abc",
	}
	require.Equal(t, exp, act)

	err := errors.New("Error: abc")
	act = Choose("abc", err)
	exp = effect{
		error: err,
	}
	require.Equal(t, exp, act)
}

func Test_effect_Named_1(t *testing.T) {
	ef := effect{}

	act := ef.Named("Bob")
	exp := effect{
		name: "Bob",
	}

	require.Equal(t, exp, act)
}

func Test_effect_Dispels_1(t *testing.T) {
	ef := effect{}

	act := ef.Dispels()
	exp := effect{
		endSpell: true,
	}

	require.Equal(t, exp, act)
}
