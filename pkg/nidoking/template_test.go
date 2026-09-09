package nidoking

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
			{{q_marks}}
		)
	`).
		Repeat("q_marks", ",", "?", 4).
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

func Test_Map_1(t *testing.T) {
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
		Map("columns", 6, func(i int) (string, bool) {
			if i >= len(columns) {
				return "", false
			}
			if i+1 >= len(columns) {
				return columns[i], true
			}
			return columns[i] + ",", true
		}).
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

func Test_RepeatLines_1(t *testing.T) {
	act := Given(`
		INSERT INTO players (
			name,
			level,
			role
		)
		{{>>>}}
		VALUES (
			?,
			?,
			?
		)
		{{<<<}}
	`).
		RepeatLines(">>>", "<<<", 3).
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
			?
		)
		VALUES (
			?,
			?,
			?
		)
		VALUES (
			?,
			?,
			?
		)
	`

	require.Equal(t, exp, act)
}

func Test_RemoveLines_1(t *testing.T) {
	act := Given(`
		SELECT
			name,
			level,
			role
		FROM
			players
		{{>>>}}
		WEHRE
			name = ?
		{{<<<}}
	`).
		RemoveLines(">>>", "<<<").
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

func Test_KeepLines_1(t *testing.T) {
	act := Given(`
		SELECT
			name,
			level,
			role
		FROM
			players
		{{>>>}}
		WEHRE
			name = ?
		{{<<<}}
	`).
		KeepLines(">>>", "<<<").
		String()

	exp := `
		SELECT
			name,
			level,
			role
		FROM
			players
		WEHRE
			name = ?
	`

	require.Equal(t, exp, act)
}
