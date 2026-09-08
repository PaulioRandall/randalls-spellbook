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

// TODO: Repeat(key, delim string, value any, count int)
// TODO: FmtInlineJoin(key, delim string, values ...any)
// TODO: FmtInlineRepeat(key, value, delim string, count int)
// TODO: Map(key string, func(i int) string)
// TODO: Reduce(key string, func(i int, acc string) string)
