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

func Test_NeedleInHaystack_RuneLineStart_1(t *testing.T) {
	// NeedleInHaystack.RuneLineStart happy path.
	require.Equal(t, 6, nihJen.RuneLineStart())
}

func Test_NeedleInHaystack_RuneLineEnd_1(t *testing.T) {
	// NeedleInHaystack.RuneLineEnd happy path.
	require.Equal(t, 11, nihJen.RuneLineEnd())
}

func Test_NeedleInHaystack_RuneStart_1(t *testing.T) {
	// NeedleInHaystack.RuneStart happy path.
	require.Equal(t, 7, nihJen.RuneStart())
}

func Test_NeedleInHaystack_RuneEnd_1(t *testing.T) {
	// NeedleInHaystack.RuneEnd happy path.
	require.Equal(t, 10, nihJen.RuneEnd())
}

func Test_NeedleInHaystack_RuneInlineStart_1(t *testing.T) {
	// NeedleInHaystack.RuneInlineStart happy path.
	require.Equal(t, 1, nihJen.RuneInlineStart())
}

func Test_NeedleInHaystack_RuneInlineEnd_1(t *testing.T) {
	// NeedleInHaystack.RuneInlineEnd happy path.
	require.Equal(t, 4, nihJen.RuneInlineEnd())
}

func Test_NeedleInHaystack_Replace_1(t *testing.T) {
	// NeedleInHaystack.Replace happy path.
	exp := `
		alice,
		dave,
		charlie
	`
	require.Equal(t, exp, nihBob.Replace("dave"))
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

func Test_NeedleInHaystack_RemoveLine_1(t *testing.T) {
	// NeedleInHaystack.RemoveLine happy path.
	exp := `
		alice,
		charlie
	`
	require.Equal(t, exp, nihBob.RemoveLine())
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

func Test_NeedleInHaystack_ReplaceFindNext_1(t *testing.T) {
	// NeedleInHaystack.ReplaceFindNext happy path.

	haystack := `
		alice,
		bob,
		charlie,
		alice, bob, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	nih = nih.ReplaceFindNext("dave")

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

// 😁 uses 4 bytes
var nihJen = NeedleInHaystack{
	LineIndex: 1,
	LineStart: 6,
	LineEnd:   14,
	Start:     10,
	End:       13,
	Needle:    "jen",
	Haystack:  "alice\n😁jen?\ncharlie",
}
