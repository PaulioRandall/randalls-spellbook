package nidorino

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Fmt_1(t *testing.T) {
	act := Given(`DROP TABLE IF EXISTS {{table}}`).
		Fmt("table", "players").
		String()

	exp := `DROP TABLE IF EXISTS players`

	require.Equal(t, exp, act)
}

func Test_FmtRepeat_1(t *testing.T) {
	act := Given(`
		SELECT
			name,
			level,
			role
		FROM
			players
		WHERE
			name in [{{q_marks}}]
	`).
		FmtRepeat("q_marks", ", ", "?", 4).
		String()

	exp := `
		SELECT
			name,
			level,
			role
		FROM
			players
		WHERE
			name in [?, ?, ?, ?]
	`

	require.Equal(t, exp, act)
}

func Test_InlineJoin_1(t *testing.T) {
	columns := []string{
		"name",
		"level",
		"role",
	}

	act := Given(`
		SELECT
			{{concat}} AS overview
		FROM
			players
	`).
		FmtJoin("concat", " || '-' || ", columns...).
		String()

	exp := `
		SELECT
			name || '-' || level || '-' || role AS overview
		FROM
			players
	`

	require.Equal(t, exp, act)
}

func Test_Join_1(t *testing.T) {
	columns := []string{
		"name",
		"level",
		"role",
	}

	act := Given(`
		SELECT
			{{columns}}
		FROM
			players
	`).
		Join("columns", ",", columns...).
		String()

	exp := `
		SELECT
			name,
			level,
			role
		FROM
			players
	`

	require.Equal(t, exp, act)
}

func Test_Repeat_1(t *testing.T) {
	act := Given(`
		INSERT INTO players (
			name,
			level,
			role
		)
		VALUES (
			{{q_mark}}
		)
	`).
		Repeat("q_mark", ",", "?", 4).
		String()

	exp := `
		INSERT INTO players (
			name,
			level,
			role
		)
		VALUES (
			?,
			?,
			?,
			?
		)
	`

	require.Equal(t, exp, act)
}

// TODO: Map(key string, func(i int) string)
// TODO: Reduce(key string, func(i int, acc string) string)
// TODO: CopyLines(start, end, to int)
