package sprints

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinLines(s ...string) string {
	return strings.Join(s, "\n")
}

type Player struct {
	Name       string
	Level      int
	difficulty string
}

func Test_Sprints_1(t *testing.T) {
	p := Player{
		Name:       "Bob",
		Level:      69,
		difficulty: "Hard",
	}

	act := Sprints(p)

	exp := joinLines(
		`type Player struct {`,
		`	Name: "Bob",`,
		`	Level: int(69),`,
		`	difficulty: <unexported>,`,
		`}`,
	)

	require.Equal(t, exp, act)
}
