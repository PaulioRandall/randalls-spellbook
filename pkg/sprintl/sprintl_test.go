package sprintl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Test_Sprintl_Fmt_1(t *testing.T) {
	// Test Fmt handles multiple formatted lines.

	act := Lines(
		"SELECT",
		"  %s",
		"FROM",
		"  %s",
	).
		Fmt(2, "Name").
		Fmt(4, "Users").
		String()

	exp := joinLines(
		"SELECT",
		"  Name",
		"FROM",
		"  Users",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Rep_1(t *testing.T) {
	// Test Rep accepts individual values.

	values := []any{
		"Name",
		"Age",
		"Height",
	}

	act := Lines(
		"(",
		"  %s",
		")",
	).
		Rep(2, values...).
		String()

	exp := joinLines(
		"(",
		"  Name",
		"  Age",
		"  Height",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Rep_2(t *testing.T) {
	// Test Rep accepts empty values.

	values := []any{
		// Empty.
	}

	act := Lines(
		"(",
		"  %s",
		")",
	).
		Rep(2, values...).
		String()

	exp := joinLines(
		"(",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Rep_3(t *testing.T) {
	// Test Rep accepts args as values.

	values := []any{
		[]any{1, "Name"},
		[]any{2, "Age"},
		[]any{3, "Height"},
	}

	act := Lines(
		"(",
		"  %d: %s",
		")",
	).
		Rep(2, values...).
		String()

	exp := joinLines(
		"(",
		"  1: Name",
		"  2: Age",
		"  3: Height",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Gen_1(t *testing.T) {
	// Test Gen LineGenerator max iterations.

	var gen LineGenerator = func(lf LineFormatter) (string, bool) {
		n := lf.Index() + 1
		s := lf.Fmt(n)
		return s, true
	}

	act := Lines(
		"(",
		"  %d",
		")",
	).
		Gen(2, 3, gen).
		String()

	exp := joinLines(
		"(",
		"  1",
		"  2",
		"  3",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Gen_2(t *testing.T) {
	// Test Gen LineGenerator returning false.

	var gen LineGenerator = func(lf LineFormatter) (string, bool) {
		if lf.Index() > 2 {
			return "", false
		}

		n := lf.Index() + 1
		s := lf.Fmt(n)
		return s, true
	}

	act := Lines(
		"(",
		"  %d",
		")",
	).
		Gen(2, 10, gen).
		String()

	exp := joinLines(
		"(",
		"  1",
		"  2",
		"  3",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Join_1(t *testing.T) {
	// Test Join applies the delimiter correctly.

	values := []any{
		"Name",
		"Age",
		"Height",
	}

	act := Lines(
		"(",
		"  %s",
		")",
	).
		Rep(2, values...).
		Join(",").
		String()

	exp := joinLines(
		"(",
		"  Name,",
		"  Age,",
		"  Height",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Marry_1(t *testing.T) {
	// Test Marry applies the delimiter correctly.

	values := []any{
		"name",
		"age",
		"height",
	}

	act := Lines(
		"%s",
	).
		Rep(1, values...).
		Marry("• ", "and ").
		String()

	exp := joinLines(
		"• name",
		"• and age",
		"• and height",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Trim_1(t *testing.T) {
	// Test Trim applies correctly.

	act := Lines(
		"  %s %s %s  ", // Line 1
		"  %s  ",       // Line 2
	).
		Fmt(1, "A", "B", "C").
		Trim().
		Dup(2, 3, "D").
		Trim().
		String()

	exp := joinLines(
		"A B C",
		"D",
		"D",
		"D",
	)

	require.Equal(t, exp, act)
}
