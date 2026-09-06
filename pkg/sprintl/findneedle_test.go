package sprintl

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

func Test_FindNeedles_1(t *testing.T) {
	// FindNeedles finds all needles.

	haystack := `
		alice,
		bob,
		charlie,
		alice, bob, charlie
	`

	act := FindNeedles(haystack, "bob", 0)
	exp := []NeedleInHaystack{
		NeedleInHaystack{
			LineIndex: 2,
			LineStart: 10,
			LineEnd:   16,
			Start:     12,
			End:       15,
			Needle:    "bob",
			Haystack:  haystack,
		},
		NeedleInHaystack{
			LineIndex: 4,
			LineStart: 28,
			LineEnd:   49,
			Start:     37,
			End:       40,
			Needle:    "bob",
			Haystack:  haystack,
		},
	}

	require.Equal(t, exp, act)
}

func Test_NeedleInHaystack_LineNum_1(t *testing.T) {
	// NeedleInHaystack.LineNum happy path.

	haystack := `
		alice,
		bob,
		charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	require.Equal(t, 3, nih.LineNum())
}

func Test_NeedleInHaystack_Line_1(t *testing.T) {
	// NeedleInHaystack.Line happy path.

	haystack := `
		alice,
		bob,
		charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	require.Equal(t, "		bob,", nih.Line())
}

func Test_NeedleInHaystack_StartInline_1(t *testing.T) {
	// NeedleInHaystack.InlineStart happy path.

	haystack := `
		alice,
		bob,
		charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	require.Equal(t, 2, nih.InlineStart())
}

func Test_NeedleInHaystack_EndInline_1(t *testing.T) {
	// NeedleInHaystack.InlineEnd happy path.

	haystack := `
		alice,
		bob,
		charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	require.Equal(t, 5, nih.InlineEnd())
}

func Test_NeedleInHaystack_Replace_1(t *testing.T) {
	// NeedleInHaystack.Replace happy path.

	haystack := `
		alice,
		bob,
		charlie
	`

	nih := FindNeedle(haystack, "bob", 0)

	exp := `
		alice,
		dave,
		charlie
	`

	require.Equal(t, exp, nih.Replace("dave"))
}

func Test_NeedleInHaystack_ReplaceFindNext_1(t *testing.T) {
	// NeedleInHaystack.ReplaceFindNext happy path.

	haystack := `
		alice,
		bob,
		charlie,
		alice, bob, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)

	expHaystack := `
		alice,
		dave,
		charlie,
		alice, bob, charlie
	`
	exp := NeedleInHaystack{
		LineIndex: 4,
		LineStart: 29,
		LineEnd:   50,
		Start:     38,
		End:       41,
		Needle:    "bob",
		Haystack:  expHaystack,
	}

	nih = nih.ReplaceFindNext("dave")
	require.Equal(t, exp, nih)
	require.Equal(
		t,
		NeedleInHaystack{},
		nih.ReplaceFindNext("dave"),
	)
}
