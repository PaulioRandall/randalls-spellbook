package sprintl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Test_Sprintl_F_1(t *testing.T) {
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
		R(2, values...).
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
	values := []any{
		// Empty.
	}

	act := Lines(
		"(",
		"  %s",
		")",
	).
		R(2, values...).
		String()

	exp := joinLines(
		"(",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_RF_1(t *testing.T) {
	values := [][]any{
		[]any{1, "Name"},
		[]any{2, "Age"},
		[]any{3, "Height"},
	}

	act := Lines(
		"(",
		"  %d: %s",
		")",
	).
		RF(2, values...).
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

func Test_Sprintl_J_1(t *testing.T) {
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
		J(2, ",", values...).
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

func Test_Sprintl_JF_1(t *testing.T) {
	values := [][]any{
		[]any{1, "Name"},
		[]any{2, "Age"},
		[]any{3, "Height"},
	}

	act := Lines(
		"(",
		"  %d: %s",
		")",
	).
		JF(2, ",", values...).
		String()

	exp := joinLines(
		"(",
		"  1: Name,",
		"  2: Age,",
		"  3: Height",
		")",
	)

	require.Equal(t, exp, act)
}

func Test_Sprintl_G_1(t *testing.T) {
	var gen LineGenerator
	gen = func(i int, line string) (string, bool) {
		if i >= 3 {
			return "", false
		}

		i++
		s := fmt.Sprintf(line, i)

		if i < 3 {
			s += ","
		}

		return s, true
	}

	act := Lines(
		"(",
		"  %d",
		")",
	).
		G(2, gen).
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

func Test_Sprintl_GN_1(t *testing.T) {
	var gen LineGenerator
	gen = func(i int, line string) (string, bool) {
		if i%2.0 == 0 {
			return "", false
		}

		i++
		s := fmt.Sprintf(line, i)
		return s, true
	}

	act := Lines(
		"(",
		"  %d",
		")",
	).
		GN(2, 7, gen).
		String()

	exp := joinLines(
		"(",
		"  2",
		"  4",
		"  6",
		")",
	)

	require.Equal(t, exp, act)
}
