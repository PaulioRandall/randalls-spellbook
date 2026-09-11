package nidoking

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
		Mode:      ModeString,
		Needle:    "bob",
		Pattern:   "bob",
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
		Mode:      ModeString,
		Needle:    "bob",
		Pattern:   "bob",
		Haystack:  haystack,
	}

	require.Equal(t, exp, act)
}

func Test_Find_3(t *testing.T) {
	// Match panics when passed pattern with linefeed or
	// carriage return.

	haystack := `
		alice,
		bob,
		charlie
	`

	require.Panics(t, func() {
		Find(haystack, "[a-z]\n[0-9]", 0)
	})

	require.Panics(t, func() {
		Find(haystack, "[a-z]\r[0-9]", 0)
	})
}

func Test_Match_1(t *testing.T) {
	// Match matches the first needle given zero start.

	haystack := `
		alice,
		bob,
		charlie
	`

	act := Match(haystack, "li.?e", 0)
	exp := NeedleInHaystack{
		LineIndex: 1,
		LineStart: 1,
		LineEnd:   9,
		Start:     4,
		End:       8,
		Mode:      ModeRegexp,
		Needle:    "lice",
		Pattern:   "li.?e",
		Haystack:  haystack,
	}

	require.Equal(t, exp, act)
}

func Test_Match_2(t *testing.T) {
	// Match matches the second needle given non-zero start.

	haystack := `
		alice,
		bob,
		charlie
	`

	act := Match(haystack, "li.?e", 8)
	exp := NeedleInHaystack{
		LineIndex: 3,
		LineStart: 17,
		LineEnd:   26,
		Start:     23,
		End:       26,
		Mode:      ModeRegexp,
		Needle:    "lie",
		Pattern:   "li.?e",
		Haystack:  haystack,
	}

	require.Equal(t, exp, act)
}

func Test_Match_4(t *testing.T) {
	// Match panics when passed pattern with multi-line
	// activated.

	haystack := `
		alice,
		bob,
		charlie
	`

	require.Panics(t, func() {
		Match(haystack, "abc\nxyz", 0)
	})

	require.Panics(t, func() {
		Match(haystack, "abc\rxyz", 0)
	})

	require.Panics(t, func() {
		Match(haystack, "abc(?im)xyz", 0)
	})

	require.Panics(t, func() {
		Match(haystack, "abc(?is)xyz", 0)
	})

	require.NotPanics(t, func() {
		Match(haystack, "abc(?i)xyz", 0)
	})
}
