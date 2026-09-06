package nidoking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FindNeedle_1(t *testing.T) {
	// FindNeedle finds the first needle given zero start.

	haystack := `
		alice,
		bob,
		charlie
	`

	act := FindNeedle(haystack, "bob", 0)
	exp := NeedleInHaystack{
		LineIndex: 2,
		LineStart: 10,
		LineEnd:   16,
		Start:     12,
		End:       15,
		Needle:    "bob",
		Haystack:  haystack,
	}

	require.Equal(t, exp, act)
}

func Test_FindNeedle_2(t *testing.T) {
	// FindNeedle finds the first needle given non-zero
	// start.

	haystack := `
		alice,
		bob,
		charlie
	`

	act := FindNeedle(haystack, "bob", 5)
	exp := NeedleInHaystack{
		LineIndex: 2,
		LineStart: 10,
		LineEnd:   16,
		Start:     12,
		End:       15,
		Needle:    "bob",
		Haystack:  haystack,
	}

	require.Equal(t, exp, act)
}
