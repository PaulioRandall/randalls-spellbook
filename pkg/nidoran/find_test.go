package nidoran

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Find_1(t *testing.T) {
	// Find finds the first needle given zero start.

	haystack := `
		alice,
		bob,
		charlie
	`

	act := Find(haystack, "bob", 0)
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

func Test_Find_2(t *testing.T) {
	// Find finds the first needle given non-zero
	// start.

	haystack := `
		bob,
		alice,
		bob,
		charlie
	`

	act := Find(haystack, "bob", 10)
	exp := NeedleInHaystack{
		LineIndex: 3,
		LineStart: 17,
		LineEnd:   23,
		Start:     19,
		End:       22,
		Needle:    "bob",
		Haystack:  haystack,
	}

	require.Equal(t, exp, act)
}
