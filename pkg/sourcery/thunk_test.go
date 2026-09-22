package sourcery

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateValErrFunc_1(t *testing.T) {
	// GIVEN f is not a function.
	// WHEN calling ValidateValErrFunc.
	// THEN returns named error.

	e := ValidateValErrFunc(0)
	require.ErrorIs(t, e, ErrNotFunc)
}

func Test_ValidateValErrFunc_2(t *testing.T) {
	// GIVEN f has 3 or more outputs.
	// WHEN calling ValidateValErrFunc.
	// THEN returns named error.

	f := func() (bool, int, error) {
		return true, 2, nil
	}

	e := ValidateValErrFunc(f)
	require.ErrorIs(t, e, ErrTooManyOutputs)
}

func Test_ValidateValErrFunc_3(t *testing.T) {
	// GIVEN f has 0, 1, or 2 outputs.
	// WHEN calling ValidateValErrFunc.
	// THEN returns nil error.

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

func Test_WithNoArgs_1(t *testing.T) {
	// GIVEN valid func.
	// WHEN calling WithNoArgs.
	// THEN returns ValErrFunc with Func set and Args as nil.

	thunk, e := WithNoArgs(func() {})

	require.NoError(t, e)
	require.NotNil(t, thunk.Func)
	require.Nil(t, thunk.Args)
}

func Test_WithJsonArgs_1(t *testing.T) {
	// GIVEN bad JSON string.
	// WHEN calling WithJsonArgs.
	// THEN returns named error.

	_, e := WithJsonArgs(
		func() {},
		"",
	)

	require.ErrorIs(t, e, ErrBadJsonString)
}

func Test_WithJsonArgs_2(t *testing.T) {
	// GIVEN too few args for non-variadic func.
	// WHEN calling WithJsonArgs.
	// THEN returns named error.

	_, e := WithJsonArgs(
		func(v any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_WithJsonArgs_3(t *testing.T) {
	// GIVEN too many args for non-variadic func.
	// WHEN calling WithJsonArgs.
	// THEN returns named error.

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
	// GIVEN too few args for variadic func.
	// WHEN calling WithJsonArgs.
	// THEN returns named error.

	_, e := WithJsonArgs(
		func(v any, more ...any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_WithJsonArgs_5(t *testing.T) {
	// GIVEN valid number of args for non-variadic func.
	// WHEN calling WithJsonArgs.
	// THEN returns no error.

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
	// GIVEN valid number of args for variadic func.
	// WHEN calling WithJsonArgs.
	// THEN returns no error.

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
	// GIVEN no args for func with no parameters.
	// WHEN calling WithJsonArgs.
	// THEN returns ValErrFunc with no args.

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
	// GIVEN valid arguments for non-variadic func.
	// WHEN calling WithJsonArgs.
	// THEN sets ValErrFunction.Args with the arguments.

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
	require.Equal(t, "abc", thunk.Args[0].Interface())
	require.Equal(t, 123, thunk.Args[1].Interface())
	require.Equal(t, 2, len(thunk.Args))
}

func Test_WithJsonArgs_9(t *testing.T) {
	// GIVEN no variadic args for variadic func.
	// WHEN calling WithJsonArgs.
	// THEN sets ValErrFunction.Args with the arguments.

	thunk, e := WithJsonArgs(
		func(s string, numbers ...float64) {},
		`
			[
				"abc"
			]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, "abc", thunk.Args[0].Interface())
	require.Equal(t, 1, len(thunk.Args))
}

func Test_WithJsonArgs_10(t *testing.T) {
	// GIVEN many variadic args for variadic func.
	// WHEN calling WithJsonArgs.
	// THEN sets ValErrFunction.Args with the arguments.

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
	require.Equal(t, "abc", thunk.Args[0].Interface())
	require.Equal(t, 1.11, thunk.Args[1].Interface())
	require.Equal(t, 2.22, thunk.Args[2].Interface())
	require.Equal(t, 3.33, thunk.Args[3].Interface())
	require.Equal(t, 4, len(thunk.Args))
}

func Test_WithJsonArgs_11(t *testing.T) {
	// Test not really needed but wanted extra confidence
	// there wasn't any issue with paarsing structs and
	// arrays.

	// GIVEN func has struct and array parameters.
	// WHEN calling WithJsonArgs.
	// THEN sets ValErrFunction.Args with the arguments.

	type DummyArg struct {
		One int
		Two string
	}

	thunk, e := WithJsonArgs(
		func(
			dummy DummyArg,
			dummies []DummyArg,
		) {
		},
		`
			[
				{
					"One": 123,
					"Two": "abc",
					"Ignored": true
				},
				[
					{
						"One": 4
					},
					{
						"One": 5
					}
				]
			]
		`,
	)

	require.NoError(t, e)

	expStruct := DummyArg{One: 123, Two: "abc"}
	require.Equal(t, expStruct, thunk.Args[0].Interface())

	expArray := []DummyArg{
		DummyArg{One: 4},
		DummyArg{One: 5},
	}
	require.Equal(t, expArray, thunk.Args[1].Interface())

	require.Equal(t, 2, len(thunk.Args))
}

func Test_ValErrFunc_Call_1(t *testing.T) {
	// GIVEN func with many parameters.
	// WHEN calling ValErrFunc.CallRecover.
	// THEN function is called with expected arguments.

	var v1 int
	var v2 string
	var v3 []float64

	thunk := ValErrFunc{
		Func: func(arg1 int, arg2 string, arg3 ...float64) {
			v1 = arg1
			v2 = arg2
			v3 = arg3
		},
		Args: []reflect.Value{
			reflect.ValueOf(123),
			reflect.ValueOf("abc"),
			reflect.ValueOf(1.11),
			reflect.ValueOf(2.22),
			reflect.ValueOf(3.33),
		},
	}

	_, _ = thunk.Call()

	require.Equal(t, 123, v1)
	require.Equal(t, "abc", v2)

	exp3 := []float64{1.11, 2.22, 3.33}
	require.Equal(t, exp3, v3)
}

func Test_ValErrFunc_Call_2(t *testing.T) {
	// GIVEN func with outputs.
	// WHEN calling ValErrFunc.Call.
	// THEN returns values returned by the function.

	var BadToTheBone = errors.New("Bad to the bone")

	thunk := ValErrFunc{
		Func: func() (string, error) {
			return "abc", BadToTheBone
		},
	}

	v, e := thunk.Call()

	require.ErrorIs(t, e, BadToTheBone)
	require.Equal(t, "abc", v)
}

func Test_ValErrFunc_CallRecover_1(t *testing.T) {
	// GIVEN func that panics.
	// WHEN calling ValErrFunc.CallRecover.
	// Then recovers and returns recovered value.

	thunk := ValErrFunc{
		Func: func() {
			panic("Moo")
		},
	}

	_, _, r := thunk.CallRecover()
	require.Equal(t, "Moo", r)
}
