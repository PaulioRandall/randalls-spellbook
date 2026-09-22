package sourcery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewOperation_1(t *testing.T) {
	// Given f is not a function.
	// Return named error.

	_, e := NewOperation(0)

	require.ErrorIs(t, e, ErrNotFunc)
}

func Test_NewOperation_2(t *testing.T) {
	// Given f has 3 or more outputs.
	// Return named error.

	f := func() (bool, int, error) {
		return true, 2, nil
	}
	_, e := NewOperation(f)

	require.ErrorIs(t, e, ErrTooManyOutputs)
}

func Test_NewOperation_3(t *testing.T) {
	// Given f has 0, 1, or 2 outputs.
	// Return nil error.

	f0 := func() { return }
	_, e := NewOperation(f0)
	require.NoError(t, e)

	f1 := func() int { return 2 }
	_, e = NewOperation(f1)
	require.NoError(t, e)

	f2 := func() (int, error) { return 2, nil }
	_, e = NewOperation(f2)
	require.NoError(t, e)
}

func Test_NewOperation_4(t *testing.T) {
	// Given valid func.
	// Returns Operation with Func set.

	f := func() {}
	op, e := NewOperation(f)

	require.NoError(t, e)
	require.NotNil(t, op.Func)
}

func Test_Operation_WithJsonArgs_1(t *testing.T) {
	// Given bad JSON string.
	// Returns named error.

	f := func() {}
	op, _ := NewOperation(f)

	_, e := op.WithJsonArgs("")
	require.ErrorIs(t, e, ErrBadJsonString)
}

func Test_Operation_WithJsonArgs_2(t *testing.T) {
	// Given too few args for non-variadic func.
	// Returns named error.

	f := func(v any) {}
	op, _ := NewOperation(f)

	_, e := op.WithJsonArgs(`
		[]
	`)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_Operation_WithJsonArgs_3(t *testing.T) {
	// Given too many args for non-variadic func.
	// Returns named error.

	f := func() {}
	op, _ := NewOperation(f)

	_, e := op.WithJsonArgs(`
		[
			123
		]
	`)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_Operation_WithJsonArgs_4(t *testing.T) {
	// Given too few args for variadic func.
	// Returns named error.

	f := func(v any, more ...any) {}
	op, _ := NewOperation(f)

	_, e := op.WithJsonArgs(`
		[]
	`)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_Operation_WithJsonArgs_5(t *testing.T) {
	// Given valid number of args for non-variadic func.
	// Returns no error.

	f := func(v1 any, v2 any) {}
	op, _ := NewOperation(f)

	_, e := op.WithJsonArgs(`
		[
			123,
			"abc"
		]
	`)

	require.NoError(t, e)
}

func Test_Operation_WithJsonArgs_6(t *testing.T) {
	// Given valid number of args for variadic func.
	// Returns no error.

	f := func(v1 any, v2 any, more ...any) {}
	op, _ := NewOperation(f)

	_, e := op.WithJsonArgs(`
		[
			123,
			"abc"
		]
	`)

	require.NoError(t, e)
}

func Test_Operation_WithJsonArgs_7(t *testing.T) {
	// Given valid arguments for non-variadic func.
	// Sets Operation.Args with them.

	f := func(s string, i int) {}
	op, _ := NewOperation(f)

	act, e := op.WithJsonArgs(`
		[
			"abc",
			123
		]
	`)

	require.NoError(t, e)
	require.Equal(t, "abc", act.Args[0])
	require.Equal(t, 123, act.Args[1])
	require.Equal(t, 2, len(act.Args))
}

func Test_Operation_WithJsonArgs_8(t *testing.T) {
	// Given no variadic args for variadic func.
	// Sets Operation.Args with them.

	f := func(s string, numbers ...float64) {}
	op, _ := NewOperation(f)

	act, e := op.WithJsonArgs(`
		[
			"abc"
		]
	`)

	require.NoError(t, e)
	require.Equal(t, "abc", act.Args[0])
	require.Equal(t, 1, len(act.Args))
}

func Test_Operation_WithJsonArgs_9(t *testing.T) {
	// Given many variadic args for variadic func.
	// Sets Operation.Args with them.

	f := func(s string, numbers ...float64) {}
	op, _ := NewOperation(f)

	act, e := op.WithJsonArgs(`
		[
			"abc",
			1.11,
			2.22,
			3.33
		]
	`)

	require.NoError(t, e)
	require.Equal(t, "abc", act.Args[0])
	require.Equal(t, 1.11, act.Args[1])
	require.Equal(t, 2.22, act.Args[2])
	require.Equal(t, 3.33, act.Args[3])
	require.Equal(t, 4, len(act.Args))
}
