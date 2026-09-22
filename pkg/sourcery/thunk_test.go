package sourcery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateValErrFunc_1(t *testing.T) {
	// Given f is not a function.
	// Return named error.

	e := ValidateValErrFunc(0)
	require.ErrorIs(t, e, ErrNotFunc)
}

func Test_ValidateValErrFunc_2(t *testing.T) {
	// Given f has 3 or more outputs.
	// Return named error.

	f := func() (bool, int, error) {
		return true, 2, nil
	}

	e := ValidateValErrFunc(f)
	require.ErrorIs(t, e, ErrTooManyOutputs)
}

func Test_ValidateValErrFunc_3(t *testing.T) {
	// Given f has 0, 1, or 2 outputs.
	// Return nil error.

	var e error

	f0 := func() { return }
	e = ValidateValErrFunc(f0)
	require.NoError(t, e)

	f1 := func() int { return 2 }
	e = ValidateValErrFunc(f1)
	require.NoError(t, e)

	f2 := func() (int, error) { return 2, nil }
	e = ValidateValErrFunc(f2)
	require.NoError(t, e)
}

func Test_NoArgs_1(t *testing.T) {
	// Given valid func.
	// Returns ValErrThunk with Func and nil Args.

	thunk, e := NoArgs(func() {})

	require.NoError(t, e)
	require.NotNil(t, thunk.Func)
	require.Nil(t, thunk.Args)
}

func Test_WithJsonArgs_1(t *testing.T) {
	// Given bad JSON string.
	// Returns named error.

	_, e := WithJsonArgs(
		func() {},
		"",
	)

	require.ErrorIs(t, e, ErrBadJsonString)
}

func Test_WithJsonArgs_2(t *testing.T) {
	// Given too few args for non-variadic func.
	// Returns named error.

	_, e := WithJsonArgs(
		func(v any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_WithJsonArgs_3(t *testing.T) {
	// Given too many args for non-variadic func.
	// Returns named error.

	_, e := WithJsonArgs(
		func() {},
		`
			[
				123
			]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_WithJsonArgs_4(t *testing.T) {
	// Given too few args for variadic func.
	// Returns named error.

	_, e := WithJsonArgs(
		func(v any, more ...any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_WithJsonArgs_5(t *testing.T) {
	// Given valid number of args for non-variadic func.
	// Returns no error.

	_, e := WithJsonArgs(
		func(v1 any, v2 any) {},
		`
			[
				"abc",
				123
			]
		`,
	)

	require.NoError(t, e)
}

func Test_WithJsonArgs_6(t *testing.T) {
	// Given valid number of args for variadic func.
	// Returns no error.

	_, e := WithJsonArgs(
		func(v1 any, v2 any, more ...any) {},
		`
			[
				"abc",
				123
			]
		`,
	)

	require.NoError(t, e)
}

func Test_WithJsonArgs_7(t *testing.T) {
	// Given no args for func with no parameters.
	// Returns ValErrThunk with no args.

	thunk, e := WithJsonArgs(
		func() {},
		`
			[]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, 0, len(thunk.Args))
}

func Test_WithJsonArgs_8(t *testing.T) {
	// Given valid arguments for non-variadic func.
	// Sets ValErrFunction.Args with them.

	thunk, e := WithJsonArgs(
		func(s string, i int) {},
		`
			[
				"abc",
				123
			]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, "abc", thunk.Args[0])
	require.Equal(t, 123, thunk.Args[1])
	require.Equal(t, 2, len(thunk.Args))
}

func Test_WithJsonArgs_9(t *testing.T) {
	// Given no variadic args for variadic func.
	// Sets ValErrFunction.Args with them.

	thunk, e := WithJsonArgs(
		func(s string, numbers ...float64) {},
		`
			[
				"abc"
			]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, "abc", thunk.Args[0])
	require.Equal(t, 1, len(thunk.Args))
}

func Test_WithJsonArgs_10(t *testing.T) {
	// Given many variadic args for variadic func.
	// Sets ValErrFunction.Args with them.

	thunk, e := WithJsonArgs(
		func(s string, numbers ...float64) {},
		`
			[
				"abc",
				1.11,
				2.22,
				3.33
			]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, "abc", thunk.Args[0])
	require.Equal(t, 1.11, thunk.Args[1])
	require.Equal(t, 2.22, thunk.Args[2])
	require.Equal(t, 3.33, thunk.Args[3])
	require.Equal(t, 4, len(thunk.Args))
}
