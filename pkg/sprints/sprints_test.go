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
	Role       Role
}

type Role struct {
	Name      string
	Strength  uint8
	Stamina   uint8
	Intellect uint8
}

var p = Player{
	Name:       "Bob",
	Level:      69,
	difficulty: "Easy",
	Role:       Role{},
}

func Test_Sprints_1(t *testing.T) {
	// Simple, no options.

	act := String(p)

	exp := joinLines(
		`Player {`,
		`	Name: "Bob",`,
		`	Level: int(69),`,
		`	Role: Role{...},`,
		`}`,
	)

	require.Equal(t, exp, act)
}

func Test_Sprints_2(t *testing.T) {
	// OptionShowUnexported option.

	act := String(p, OptionShowUnexported)

	exp := joinLines(
		`Player {`,
		`	Name: "Bob",`,
		`	Level: int(69),`,
		`	difficulty: <unexported>,`,
		`	Role: Role{...},`,
		`}`,
	)

	require.Equal(t, exp, act)
}
