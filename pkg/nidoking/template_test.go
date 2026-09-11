package nidoking

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type dbTable struct {
	Name        string
	columnNames []string
}

func (dbt dbTable) IdColumnName() string {
	return dbt.columnNames[0]
}

type dbColumn struct {
	Name            string
	Type            string
	notNull         bool
	unexportedField int
}

func (col dbColumn) Constraints() string {
	if col.notNull {
		return "NOT NULL"
	}
	return ""
}

func (col dbColumn) MethodTooManyInputs(param string) string {
	return ""
}

func (col dbColumn) MethodTooManyOutputs() (string, error) {
	return "", nil
}

func (col dbColumn) MethodTooFewOutputs() {
}

func Test_Template_Inline_1(t *testing.T) {
	act := Given(`DROP TABLE IF EXISTS {{table}}`).
		Inline("table", "players").
		String()

	exp := `DROP TABLE IF EXISTS players`

	require.Equal(t, exp, act)
}

func Test_Template_InlineRepeat_1(t *testing.T) {
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
		InlineRepeat("q_marks", ", ", "?", 4).
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

func Test_Template_InlineJoin_1(t *testing.T) {
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
		InlineJoin("concat", " || '-' || ", columns...).
		String()

	exp := `
		SELECT
			name || '-' || level || '-' || role AS overview
		FROM
			players
	`

	require.Equal(t, exp, act)
}

func Test_Template_InlineMap_1(t *testing.T) {
	testTable := dbTable{
		Name: "players",
		columnNames: []string{
			"name",
			"level",
			"role",
		},
	}

	act := Given(`
		SELECT
			*
		FROM
			{{table.Name}}
		WHERE
			{{table.IdColumnName}} = ?
	`).
		InlineMap("table", testTable).
		String()

	exp := `
		SELECT
			*
		FROM
			players
		WHERE
			name = ?
	`

	require.Equal(t, exp, act)
}

func Test_Template_ListJoin_1(t *testing.T) {
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
		ListJoin("columns", ",", columns...).
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

func Test_Template_ListRepeat_1(t *testing.T) {
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
		ListRepeat("q_marks", ",", "?", 4).
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

func Test_Template_ListMapper_1(t *testing.T) {
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
		ListMapper("columns", 6, func(i int) (string, bool) {
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

func Test_Template_LinesRepeat_1(t *testing.T) {
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
		LinesRepeat(">>>", "<<<", 3).
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

func Test_Template_LinesRemove_1(t *testing.T) {
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
		LinesRemove(">>>", "<<<").
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

func Test_Template_LinesKeep_1(t *testing.T) {
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
		LinesKeep(">>>", "<<<").
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

func Test_Template_ListMap_1(t *testing.T) {
	// Happy path.

	columns := []dbColumn{
		dbColumn{
			Name:    "name",
			Type:    "TEXT",
			notNull: true,
		},
		dbColumn{
			Name:    "level",
			Type:    "INTEGER",
			notNull: true,
		},
		dbColumn{
			Name:    "role",
			Type:    "TEXT",
			notNull: false,
		},
	}

	act := Given(`
		CREATE TABLE players (
			{{col.Name}} {{col.Type}} {{col.Constraints}}
		)
	`).
		ListMap("col", ",", columns...).
		String()

	exp := `
		CREATE TABLE players (
			name TEXT NOT NULL,
			level INTEGER NOT NULL,
			role TEXT 
		)
	`

	require.Equal(t, exp, act)
}

func Test_Template_ListMap_2(t *testing.T) {
	// Panics given bad input.

	columns := []dbColumn{
		dbColumn{
			Name:    "name",
			Type:    "TEXT",
			notNull: true,
		},
		dbColumn{
			Name:    "level",
			Type:    "INTEGER",
			notNull: true,
		},
		dbColumn{
			Name:    "role",
			Type:    "TEXT",
			notNull: false,
		},
	}

	require.Panics(t, func() {
		tmpl := Given(`{{col.UnknownField}}`)
		tmpl.ListMap("col", ",", columns...)
	})

	require.Panics(t, func() {
		tmpl := Given(`{{col.unexportedField}}`)
		tmpl.ListMap("col", ",", columns...)
	})

	require.Panics(t, func() {
		tmpl := Given(`{{col.MethodTooManyInputs}}`)
		tmpl.ListMap("col", ",", columns...)
	})

	require.Panics(t, func() {
		tmpl := Given(`{{col.MethodTooManyOutputs}}`)
		tmpl.ListMap("col", ",", columns...)
	})

	require.Panics(t, func() {
		tmpl := Given(`{{col.MethodTooFewOutputs}}`)
		tmpl.ListMap("col", ",", columns...)
	})
}
