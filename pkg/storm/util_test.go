package storm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_genList_1(t *testing.T) {
	given := []string{
		"Alice",
		"Bob",
		"Charlie",
	}

	exp := joinLines(
		"  0: Alice,",
		"  1: Bob,",
		"  2: Charlie",
	)

	s, e := genList(
		given,
		func(i int, item string) (string, error) {
			return fmt.Sprintf("  %d: %s", i, item), nil
		},
	)

	require.Equal(t, nil, e)
	require.Equal(t, exp, s)
}
