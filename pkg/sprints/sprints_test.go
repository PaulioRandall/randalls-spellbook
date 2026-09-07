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

var p = Player{
	Name:       "Bob",
	Level:      69,
	difficulty: "Easy",
}

func Test_Sprints_1(t *testing.T) {
	// Simple, no options.

	act := String(p)

	exp := joinLines(
		`type Player struct {`,
		`	Name: "Bob",`,
		`	Level: int(69),`,
		`}`,
	)

	require.Equal(t, exp, act)
}

func Test_Sprints_2(t *testing.T) {
	// OptionShowUnexported option.

	act := String(p, OptionShowUnexported)

	exp := joinLines(
		`type Player struct {`,
		`	Name: "Bob",`,
		`	Level: int(69),`,
		`	difficulty: <unexported>,`,
		`}`,
	)

	require.Equal(t, exp, act)
}
