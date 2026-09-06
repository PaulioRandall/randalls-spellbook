package nidoking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var nihBob = NeedleInHaystack{
	LineIndex: 2,
	LineStart: 10,
	LineEnd:   16,
	Start:     12,
	End:       15,
	Needle:    "bob",
	Haystack: `
		alice,
		bob,
		charlie
	`,
}

func Test_NeedleInHaystack_LineNum_1(t *testing.T) {
	// NeedleInHaystack.LineNum happy path.
	require.Equal(t, 3, nihBob.LineNum())
}

func Test_NeedleInHaystack_LineText_1(t *testing.T) {
	// NeedleInHaystack.LineText happy path.
	require.Equal(t, "		bob,", nihBob.LineText())
}

func Test_NeedleInHaystack_InlineStart_1(t *testing.T) {
	// NeedleInHaystack.InlineStart happy path.
	require.Equal(t, 2, nihBob.InlineStart())
}

func Test_NeedleInHaystack_InlineEnd_1(t *testing.T) {
	// NeedleInHaystack.InlineEnd happy path.
	require.Equal(t, 5, nihBob.InlineEnd())
}

func Test_NeedleInHaystack_ReplaceNeedle_1(t *testing.T) {
	// NeedleInHaystack.ReplaceNeedle happy path.
	exp := `
		alice,
		dave,
		charlie
	`
	require.Equal(t, exp, nihBob.ReplaceNeedle("dave"))
}

func Test_NeedleInHaystack_ReplaceLine_1(t *testing.T) {
	// NeedleInHaystack.ReplaceLine happy path.
	exp := `
		alice,
****dave****,
		charlie
	`
	require.Equal(t, exp, nihBob.ReplaceLine("****dave****,"))
}

func Test_NeedleInHaystack_FindNext_1(t *testing.T) {
	// NeedleInHaystack.FindNext happy path.

	haystack := `
		alice,
		bob,
		charlie,
		alice, bob, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)

	exp := NeedleInHaystack{
		LineIndex: 4,
		LineStart: 28,
		LineEnd:   49,
		Start:     37,
		End:       40,
		Needle:    "bob",
		Haystack:  haystack,
	}

	nih = nih.FindNext()
	require.Equal(t, exp, nih)

	nih = nih.FindNext()
	require.Equal(t, NeedleInHaystack{}, nih)
}

func Test_NeedleInHaystack_ReplaceNeedleFindNext_1(t *testing.T) {
	// NeedleInHaystack.ReplaceNeedleFindNext happy path.

	haystack := `
		alice,
		bob,
		charlie,
		alice, bob, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	nih = nih.ReplaceNeedleFindNext("dave")

	haystack = `
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
		Haystack:  haystack,
	}

	require.Equal(t, exp, nih)
}

func Test_NeedleInHaystack_ReplaceLineFindNext_1(t *testing.T) {
	// NeedleInHaystack.ReplaceLineFindNext happy path.

	haystack := `
		alice,
		bob,
		charlie,
		alice, bob, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	nih = nih.ReplaceLineFindNext("****dave****,")

	haystack = `
		alice,
****dave****,
		charlie,
		alice, bob, charlie
	`

	exp := NeedleInHaystack{
		LineIndex: 4,
		LineStart: 35,
		LineEnd:   56,
		Start:     44,
		End:       47,
		Needle:    "bob",
		Haystack:  haystack,
	}

	require.Equal(t, exp, nih)
}
