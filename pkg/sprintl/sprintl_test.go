package sprintl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Test_Sprintl_F_1(t *testing.T) {
	// Test F handles multiple formatted lines.

	act := Lines(
		"SELECT",
		"  %s",
		"FROM",
		"  %s",
	).
		F(2, "Name").
		F(4, "Users").
		String()

	exp := joinLines(
		"SELECT",
		"  Name",
		"FROM",
		"  Users",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_R_1(t *testing.T) {
	// Test R accepts individual values.

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
		R(2, "", values...).
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

func Test_Sprintl_R_2(t *testing.T) {
	// Test R accepts empty values.

	values := []any{
		// Empty.
	}

	act := Lines(
		"(",
		"  %s",
		")",
	).
		R(2, "", values...).
		String()

	exp := joinLines(
		"(",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_R_3(t *testing.T) {
	// Test R accepts args as values.

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
		R(2, "", values...).
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

func Test_Sprintl_R_4(t *testing.T) {
	// Test R accepts delim values.

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
		R(2, ",", values...).
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

func Test_Sprintl_G_1(t *testing.T) {
	// Test G LineFormatter max iterations.

	var gen LineFormatter = func(f Formatter) (string, bool) {
		n := f.Index() + 1
		s := f.Fmt(n)
		return s, true
	}

	act := Lines(
		"(",
		"  %d",
		")",
	).
		G(2, ",", 3, gen).
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

func Test_Sprintl_G_2(t *testing.T) {
	// Test G LineFormatter returning false.

	var gen LineFormatter = func(f Formatter) (string, bool) {
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
		G(2, ",", 10, gen).
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
