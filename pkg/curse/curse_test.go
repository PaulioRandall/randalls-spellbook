package curse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Curse_WrapIn_1(t *testing.T) {
	// GIVEN a Curse symptom
	// WHEN  calling Curse.WrapIn
	// THEN  cause is wrapped by the Curse symptom
	// AND   that symptom is returned

	symptom := Curse{
		Message: "Symptom",
	}

	cause := Curse{
		Message: "Cause",
	}

	act := cause.WrapIn(symptom)

	exp := Curse{
		Message: "Symptom",
		Cause:   cause,
	}

	require.Equal(t, exp, act)
}

func Test_Curse_WrapIn_2(t *testing.T) {
	// GIVEN a ProtoCurse symptom
	// WHEN  calling Curse.WrapIn
	// THEN  cause is wrapped by the ProtoCurse symptom
	// AND   that symptom is returned

	symptom := ProtoCurse{
		curse: Curse{
			Message: "Symptom: %s",
		},
	}

	cause := Curse{
		Message: "Cause",
	}

	act := cause.WrapIn(symptom)

	exp := ProtoCurse{
		curse: Curse{
			Message: "Symptom: %s",
			Cause:   cause,
		},
	}

	require.Equal(t, exp, act)
}

func Test_Curse_Is_1(t *testing.T) {
	// GIVEN two curses independently created curses
	// THEN  they must be considered different in meaning

	a := Err("A")
	b := Err("B")

	require.Equal(t, false, a.Is(b))
	require.Equal(t, false, b.Is(a))
}

func Test_Curse_Is_2(t *testing.T) {
	// GIVEN a curse
	// THEN  it must be considered meaningfully the same when
	//       compared with itself

	a := Err("A")

	require.Equal(t, true, a.Is(a))
}

func Test_Curse_Is_3(t *testing.T) {
	// GIVEN a Curse B that was create from ProtoCurse A
	// THEN  they must be considered meaningfully the same

	p := Proto("%s")
	a := p.Fmt("A")
	b := p.Fmt("B")

	require.Equal(t, true, a.Is(p))
	require.Equal(t, true, b.Is(p))
	require.Equal(t, true, a.Is(b))
	require.Equal(t, true, b.Is(a))
}
