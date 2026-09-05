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

func Test_Sprintl_Range_1(t *testing.T) {
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
		Rep(2, "", values...).
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

func Test_Sprintl_Range_2(t *testing.T) {
	// Test Rep accepts empty values.

	values := []any{
		// Empty.
	}

	act := Lines(
		"(",
		"  %s",
		")",
	).
		Rep(2, "", values...).
		String()

	exp := joinLines(
		"(",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Range_3(t *testing.T) {
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
		Rep(2, "", values...).
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

func Test_Sprintl_Range_4(t *testing.T) {
	// Test Rep accepts delim values.

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
		Rep(2, ",", values...).
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

func Test_Sprintl_Gen_1(t *testing.T) {
	// Test Gen LineFormatter max iterations.

	var gen LineFormatter = func(f LineFormat) (string, bool) {
		n := f.Index() + 1
		s := f.Fmt(n)
		return s, true
	}

	act := Lines(
		"(",
		"  %d",
		")",
	).
		Gen(2, ",", 3, gen).
		String()

	exp := joinLines(
		"(",
		"  1,",
		"  2,",
		"  3",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_Gen_2(t *testing.T) {
	// Test Gen LineFormatter returning false.

	var gen LineFormatter = func(f LineFormat) (string, bool) {
		if f.Index() > 2 {
			return "", false
		}

		n := f.Index() + 1
		s := f.Fmt(n)
		return s, true
	}

	act := Lines(
		"(",
		"  %d",
		")",
	).
		Gen(2, ",", 10, gen).
		String()

	exp := joinLines(
		"(",
		"  1,",
		"  2,",
		"  3",
		")",
	)

	require.Equal(t, exp, act)
}
