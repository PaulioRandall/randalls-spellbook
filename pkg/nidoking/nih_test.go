package nidoking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var nihBob = NeedleInHaystack{
	LineStart: 10,
	Start:     12,
	End:       15,
	LineEnd:   16,
	Mode:      ModeString,
	Needle:    "bob",
	Pattern:   "bob",
	Haystack: `
		alice,
		bob,
		charlie
	`,
}

// 😁 uses 4 bytes
var nihJen = NeedleInHaystack{
	LineStart: 6,
	Start:     10,
	End:       13,
	LineEnd:   14,
	Mode:      ModeString,
	Needle:    "jen",
	Pattern:   "jen",
	Haystack:  "alice\n😁jen?\ncharlie",
}

func Test_NeedleInHaystack_LineIndex_1(t *testing.T) {
	// NeedleInHaystack.LineIndex happy path.
	require.Equal(t, 2, nihBob.LineIndex())
}

func Test_NeedleInHaystack_LineNumber_1(t *testing.T) {
	// NeedleInHaystack.LineNumber happy path.
	require.Equal(t, 3, nihBob.LineNumber())
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

func Test_NeedleInHaystack_ReplaceLine_1(t *testing.T) {
	// NeedleInHaystack.ReplaceLine happy path.
	exp := Replacement{
		Nih:   nihBob,
		Start: 10,
		End:   23,
		Haystack: `
		alice,
****dave****,
		charlie
	`,
	}
	require.Equal(t, exp, nihBob.ReplaceLine("****dave****,"))
}

func Test_NeedleInHaystack_ReplaceJoin_1(t *testing.T) {
	// NeedleInHaystack.ReplaceJoin happy path.
	columns := []string{
		"name",
		"level",
		"role",
	}

	haystack := `
		SELECT
			{{columns}}
		FROM
			players
	`

	nih := Find(haystack, "{{columns}}", 0)
	rep := nih.ReplaceJoin(columns, ",")

	exp := Replacement{
		Nih:   nih,
		Start: 10,
		End:   36,
		Haystack: `
		SELECT
			name,
			level,
			role
		FROM
			players
	`,
	}

	require.Equal(t, exp, rep)
}

func Test_NeedleInHaystack_ReplaceRepeat_1(t *testing.T) {
	// NeedleInHaystack.ReplaceRepeat happy path.
	haystack := `
		INSERT INTO players (
			name,
			level,
			role
		)
		VALUES (
			{{values}}
		)
	`

	nih := Find(haystack, "{{values}}", 0)
	rep := nih.ReplaceRepeat("?", 3, ",")

	exp := Replacement{
		Nih:   nih,
		Start: 67,
		End:   83,
		Haystack: `
		INSERT INTO players (
			name,
			level,
			role
		)
		VALUES (
			?,
			?,
			?
		)
	`,
	}

	require.Equal(t, exp, rep)
}

func Test_NeedleInHaystack_ReplaceInline_1(t *testing.T) {
	// NeedleInHaystack.ReplaceInline happy path.
	exp := Replacement{
		Nih:   nihBob,
		Start: 12,
		End:   16,
		Haystack: `
		alice,
		dave,
		charlie
	`,
	}
	require.Equal(t, exp, nihBob.ReplaceInline("dave"))
}

func Test_NeedleInHaystack_ReplaceInlineJoin_1(t *testing.T) {
	// NeedleInHaystack.ReplaceInlineJoin happy path.
	columns := []string{
		"name",
		"level",
		"role",
	}

	haystack := `
		SELECT {{columns}}
		FROM players
	`

	nih := Find(haystack, "{{columns}}", 0)
	rep := nih.ReplaceInlineJoin(columns, ", ")

	exp := Replacement{
		Nih:   nih,
		Start: 10,
		End:   27,
		Haystack: `
		SELECT name, level, role
		FROM players
	`,
	}

	require.Equal(t, exp, rep)
}

func Test_NeedleInHaystack_ReplaceInlineRepeat_1(t *testing.T) {
	// NeedleInHaystack.ReplaceInlineRepeat happy path.
	haystack := `
		SELECT
			name,
			level,
			role
		FROM
			players
		WHERE
			name IN [{{q_mark}}]
	`

	nih := Find(haystack, "{{q_mark}}", 0)
	rep := nih.ReplaceInlineRepeat("?", 3, ", ")

	exp := Replacement{
		Nih:   nih,
		Start: 75,
		End:   82,
		Haystack: `
		SELECT
			name,
			level,
			role
		FROM
			players
		WHERE
			name IN [?, ?, ?]
	`,
	}

	require.Equal(t, exp, rep)
}

func Test_NeedleInHaystack_RemoveLine_1(t *testing.T) {
	// NeedleInHaystack.RemoveLine happy path.
	exp := Replacement{
		Nih:   nihBob,
		Start: 10,
		End:   10,
		Haystack: `
		alice,
		charlie
	`,
	}
	require.Equal(t, exp, nihBob.RemoveLine())
}

func Test_NeedleInHaystack_FindNext_1(t *testing.T) {
	// NeedleInHaystack.FindNext happy path for ModeString.

	haystack := `
		alice,
		bob,
		charlie,
		alice, bob, charlie
	`

	nih := Find(haystack, "bob", 0)
	nih = nih.FindNext()

	exp := NeedleInHaystack{
		LineStart: 28,
		Start:     37,
		End:       40,
		LineEnd:   49,
		Mode:      ModeString,
		Needle:    "bob",
		Pattern:   "bob",
		Haystack:  haystack,
	}

	require.Equal(t, exp, nih)

	nih = nih.FindNext()
	require.Equal(t, NeedleInHaystack{}, nih)
}

func Test_NeedleInHaystack_FindNext_2(t *testing.T) {
	// NeedleInHaystack.FindNext happy path for ModeRegexp.

	haystack := `
		alice,
		bob,
		charlie
	`

	exp1 := NeedleInHaystack{
		LineStart: 1,
		LineEnd:   9,
		Start:     4,
		End:       8,
		Mode:      ModeRegexp,
		Needle:    "lice",
		Pattern:   "li.?e",
		Haystack:  haystack,
	}

	nih := Match(haystack, "li.?e", 0)
	require.Equal(t, exp1, nih)

	exp2 := NeedleInHaystack{
		LineStart: 17,
		Start:     23,
		End:       26,
		LineEnd:   26,
		Mode:      ModeRegexp,
		Needle:    "lie",
		Pattern:   "li.?e",
		Haystack:  haystack,
	}

	nih = nih.FindNext()
	require.Equal(t, exp2, nih)

	nih = nih.FindNext()
	require.Equal(t, NeedleInHaystack{}, nih)
}
