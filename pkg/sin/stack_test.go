package sin

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Test_StackString_1(t *testing.T) {
	// GIVEN an error chain
	// WHEN  producing a stack representation of the chain
	// THEN  the expected string representation is returned

	e1 := fmt.Errorf("Message 1")
	e2 := fmt.Errorf("Message 2: %w", e1)
	e3 := fmt.Errorf("Message 3: %w", e2)

	act := StackString(e3)
	exp := joinLines(
		"Message 3: Message 2: Message 1",
		"\t⤷ Message 2: Message 1",
		"\t⤷ Message 1",
	)

	require.Equal(t, exp, act)
}

func Test_ReverseStackString_1(t *testing.T) {
	// GIVEN an error chain
	// WHEN  producing a stack representation of the chain
	//       in reverse
	// THEN  the expected string representation is returned

	e1 := fmt.Errorf("Message 1")
	e2 := fmt.Errorf("Message 2: %w", e1)
	e3 := fmt.Errorf("Message 3: %w", e2)

	act := ReverseStackString(e3)
	exp := joinLines(
		"Message 1",
		"\t⮤ Message 2: Message 1",
		"\t⮤ Message 3: Message 2: Message 1",
	)

	require.Equal(t, exp, act)
}

func Test_AsStackErr_1(t *testing.T) {
	// Just make sure it doesn't infinitely recurse.

	e1 := fmt.Errorf("Message 1")
	e2 := fmt.Errorf("Message 2: %w", e1)
	e3 := fmt.Errorf("Message 3: %w", e2)

	err := Stack(e3)
	_ = err.Error()
}
