package effect

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Bless_1(t *testing.T) {
	act := Bless("abc")
	exp := &effect[string]{
		result: "abc",
	}

	require.Equal(t, exp, act)
}

func Test_Sin_1(t *testing.T) {
	err := errors.New("abc")

	act := Cursed[string](err)
	exp := &effect[string]{
		error: err,
	}

	require.Equal(t, exp, act)
}

func Test_Curse_1(t *testing.T) {
	act := Curse[string]("Error: %s", "abc")
	exp := &effect[string]{
		error: errors.New("Error: abc"),
	}

	require.Equal(t, exp, act)
}

func Test_Judge_1(t *testing.T) {
	var act, exp Effect[string]

	act = Judge[string]("abc")
	exp = &effect[string]{
		result: "abc",
	}
	require.Equal(t, exp, act)

	err := errors.New("Error: abc")
	act = Judge[string](err)
	exp = &effect[string]{
		error: err,
	}
	require.Equal(t, exp, act)
}

func Test_Choose_1(t *testing.T) {
	var act, exp Effect[string]

	act = Choose("abc", nil)
	exp = &effect[string]{
		result: "abc",
	}
	require.Equal(t, exp, act)

	err := errors.New("Error: abc")
	act = Choose("abc", err)
	exp = &effect[string]{
		error: err,
	}
	require.Equal(t, exp, act)
}

func Test_effect_NameAs_1(t *testing.T) {
	ef := &effect[string]{}

	act := ef.NameAs("Bob")
	exp := &effect[string]{
		name: "Bob",
	}

	require.Equal(t, exp, act)
}

func Test_Demystify_1(t *testing.T) {
	var te Effect[string] = &effect[string]{result: "abc"}
	var ue UntypedEffect = te

	var act Effect[string] = Demystify[string](ue)
	require.Equal(t, te, act)
}

func Test_SeekNamedEffect_1(t *testing.T) {
	A := &effect[string]{name: "A", prior: nil}
	B := &effect[string]{name: "B", prior: A}
	C := &effect[string]{name: "C", prior: B}
	D := &effect[string]{name: "D", prior: C}

	act := SeekNamedEffect(D, "B")
	require.Equal(t, B, act)

	act = SeekNamedEffect(B, "D")
	require.Equal(t, nil, act)
}
