package curse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func joinLines(lines ...string) string {
	return strings.Join(lines, "\n")
}

func Test_Stack_1(t *testing.T) {
	// GIVEN an error chain
	// WHEN  producing a stack representation of the chain
	// THEN  the expected string representation is returned

	e1 := fmt.Errorf("Message 1")
	e2 := fmt.Errorf("Message 2: %w", e1)
	e3 := fmt.Errorf("Message 3: %w", e2)

	act := Stack(e3, false)
	exp := joinLines(
		"Message 3",
		"\t⤷ Message 2",
		"\t⤷ Message 1",
	)

	require.Equal(t, exp, act)
}

func Test_Stack_2(t *testing.T) {
	// GIVEN an error chain
	// WHEN  producing a stack representation of the chain
	//       in reverse
	// THEN  the expected string representation is returned

	e1 := fmt.Errorf("Message 1")
	e2 := fmt.Errorf("Message 2: %w", e1)
	e3 := fmt.Errorf("Message 3: %w", e2)

	act := Stack(e3, true)
	exp := joinLines(
		"Message 1",
		"\t⮤ Message 2",
		"\t⮤ Message 3",
	)

	require.Equal(t, exp, act)
}

func Test_RawStack_1(t *testing.T) {
	// GIVEN an error chain
	// WHEN  producing a stack representation of the chain
	// THEN  the expected string representation is returned
	//       without message trimming

	e1 := fmt.Errorf("Message 1")
	e2 := fmt.Errorf("Message 2: %w", e1)
	e3 := fmt.Errorf("Message 3: %w", e2)

	act := RawStack(e3, false)
	exp := joinLines(
		"Message 3: Message 2: Message 1",
		"\t⤷ Message 2: Message 1",
		"\t⤷ Message 1",
	)

	require.Equal(t, exp, act)
}

func Test_StackErr_1(t *testing.T) {
	// Just make sure it doesn't infinitely recurse.

	e1 := fmt.Errorf("Message 1")
	e2 := fmt.Errorf("Message 2: %w", e1)
	e3 := fmt.Errorf("Message 3: %w", e2)

	err := StackErr(e3)
	_ = err.Error()
}
