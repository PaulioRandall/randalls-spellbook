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
	Role       *Role
	Class      **Class
}

type Role struct {
	Name string
}

type Class struct {
	Name string
}

var role = Role{
	Name: "Wizzard",
}

var ptrClazz *Class = nil
var clazz **Class = &ptrClazz

var p = Player{
	Name:       "Bob",
	Level:      69,
	difficulty: "Easy",
	Role:       &role,
	Class:      clazz,
}

func Test_Sprints_1(t *testing.T) {
	// Happy path.

	act := String(p)

	exp := joinLines(
		`Player {`,
		`	Name: "Bob",`,
		`	Level: int(69),`,
		`	Role: *Role{...},`,
		`	Class: *⁎Class{},`,
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
		`	Role: *Role{...},`,
		`	Class: *⁎Class{},`,
		`}`,
	)

	require.Equal(t, exp, act)
}

func Test_Sprints_3(t *testing.T) {
	// Long strings are clipped.

	type Paragraph struct {
		Content string
	}

	para := Paragraph{
		Content: joinLines(
			"No TV and no beer make Homer go crazy.",
			"No TV and no beer make Homer go crazy.",
			"No TV and no beer make Homer go crazy.",
			"No TV and no beer make Homer go crazy.",
			"No TV and no beer make Homer go crazy.",
		),
	}

	act := String(para)

	exp := joinLines(
		`Paragraph {`,
		`	Content: "No TV and no beer make Homer g...",`,
		`}`,
	)

	require.Equal(t, exp, act)
}

func Test_Sprints_4(t *testing.T) {
	// When input is a pointer it is dereferenced.
	// Struct name is preceeded by correct number of '*'

	type Thing struct {
		Name string
	}

	ptrThing := &Thing{
		Name: "Thingy",
	}

	act := String(&ptrThing)

	exp := joinLines(
		`**Thing {`,
		`	Name: "Thingy",`,
		`}`,
	)

	require.Equal(t, exp, act)
}

func Test_Sprints_5(t *testing.T) {
	// When input is an array, only its type and
	// length are printed.

	type Thing struct {
		Empty    []string
		NotEmpty []string
		Nil      []string
	}

	thing := Thing{
		Empty: []string{},
		NotEmpty: []string{
			"Thing1",
			"Thing2",
		},
		Nil: nil,
	}

	act := String(&thing)

	exp := joinLines(
		`*Thing {`,
		`	Empty: [0]string{},`,
		`	NotEmpty: [2]string{...},`,
		`	Nil: []string,`,
		`}`,
	)

	require.Equal(t, exp, act)
}
