package sourcery

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ValidateFunction_1(t *testing.T) {
	// When validating a function
	// where f is not a function
	// then a named error is returned

	e := ValidateFunction(0)
	require.ErrorIs(t, e, ErrNotFunc)
}

func Test_ValidateFunction_2(t *testing.T) {
	// When validating a function
	// where f has too many output parameters
	// then a named error is returned

	f := func() (bool, int, error) {
		return true, 2, nil
	}

	e := ValidateFunction(f)
	require.ErrorIs(t, e, ErrTooManyOutputs)
}

func Test_ValidateFunction_3(t *testing.T) {
	// When validating a function
	// where f's second output parameter is not an error
	// then a named error is returned

	f := func() (bool, int) {
		return true, 2
	}

	e := ValidateFunction(f)
	require.ErrorIs(t, e, ErrOutputNotError)
}

func Test_ValidateFunction_4(t *testing.T) {
	// When validating a function
	// where f has a valid function signature
	// then no errors are returned

	var e error

	f0 := func() { return }
	e = ValidateFunction(f0)
	require.NoError(t, e)

	f1 := func() int { return 2 }
	e = ValidateFunction(f1)
	require.NoError(t, e)

	f2 := func() (int, error) { return 2, nil }
	e = ValidateFunction(f2)
	require.NoError(t, e)
}

func Test_call_1(t *testing.T) {
	// When calling a function
	// where bad JSON is given
	// then a named error is returned

	_, e := call(
		func() {},
		"",
	)

	require.ErrorIs(t, e, ErrBadJsonString)
}

func Test_call_2(t *testing.T) {
	// When calling a function
	// that is non-variadic
	// where the JSON has less values than the amount
	// required by the function
	// then a named error is returned

	_, e := call(
		func(v any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_call_3(t *testing.T) {
	// When calling a function
	// that is non-variadic
	// where the JSON has more values than the amount
	// required by the function
	// then a named error is returned

	_, e := call(
		func() {},
		`
			[
				123
			]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_call_4(t *testing.T) {
	// When calling a function
	// that is variadic
	// where the JSON has less values than the amount of
	// non-variadic parameters required by the function
	// then a named error is returned

	_, e := call(
		func(v any, more ...any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_call_5(t *testing.T) {
	// When calling a function
	// that is non-variadic
	// where the number of JSON values matches the amount of
	// inputs required by the function
	// then no error is returned

	_, e := call(
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

func Test_call_6(t *testing.T) {
	// When calling a function
	// that is variadic
	// where the number of JSON values matches the amount of
	// non-variadic inputs required by the function
	// then no error is returned

	_, e := call(
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

func Test_call_7(t *testing.T) {
	// When calling a function
	// that is non-variadic
	// where the JSON contains values and the function
	// accepts no parameeters
	// then no error is returned
	// and the function is called
	// and the correct value is returned

	v, e := call(
		func() int { return 1 },
		`
			[]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, 1, v)
}

func Test_call_8(t *testing.T) {
	// When calling a function
	// that is non-variadic
	// where the function requires typed inputs
	// and the JSON values match the types
	// then no error is returned
	// and the function is called
	// and the correct parameters are passed

	var arg1 string
	var arg2 int

	_, e := call(
		func(s string, i int) {
			arg1 = s
			arg2 = i
		},
		`
			[
				"abc",
				123
			]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, "abc", arg1)
	require.Equal(t, 123, arg2)
}

func Test_call_9(t *testing.T) {
	// When calling a function
	// that is variadic
	// where the JSON provides only the non-variadic
	// arguments
	// then no error is returned
	// and the function is called
	// and the correct parameters are passed

	var arg1 string
	var arg2 []float64

	_, e := call(
		func(s string, numbers ...float64) {
			arg1 = s
			arg2 = numbers
		},
		`
			[
				"abc"
			]
		`,
	)

	require.NoError(t, e)
	require.Equal(t, "abc", arg1)
	require.Equal(t, []float64{}, arg2)
}

func Test_call_10(t *testing.T) {
	// When calling a function
	// that is variadic
	// where the JSON provides several values for the
	// function's variadic input
	// then no error is returned
	// and the function is called
	// and the correct parameters are passed

	var arg1 string
	var arg2 []float64

	_, e := call(
		func(s string, numbers ...float64) {
			arg1 = s
			arg2 = numbers
		},
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
	require.Equal(t, "abc", arg1)
	require.Equal(t, 1.11, arg2[0])
	require.Equal(t, 2.22, arg2[1])
	require.Equal(t, 3.33, arg2[2])
	require.Equal(t, 3, len(arg2))
}

func Test_call_11(t *testing.T) {
	// Test not really needed but wanted extra confidence
	// there wasn't any issue with parsing structs and
	// arrays.

	// When calling a function
	// where the function accepts objects
	// the JSON provides the values for those objects
	// then no error is returned
	// and the function is called
	// and the correct parameters are passed

	type DummyArg struct {
		One int
		Two string
	}

	var arg1 DummyArg
	var arg2 []DummyArg

	_, e := call(
		func(
			dummy DummyArg,
			dummies []DummyArg,
		) {
			arg1 = dummy
			arg2 = dummies
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

	expArg1 := DummyArg{One: 123, Two: "abc"}
	require.Equal(t, expArg1, arg1)

	expArg2 := []DummyArg{
		DummyArg{One: 4},
		DummyArg{One: 5},
	}
	require.Equal(t, expArg2, arg2)
	require.Equal(t, 2, len(arg2))
}

func Test_call_12(t *testing.T) {
	// When calling a function
	// with multiple output parameters
	// then the function is called
	// and the output values are returned

	var BadToTheBone = errors.New("Bad to the bone")

	v, e := call(
		func() (string, error) {
			return "abc", BadToTheBone
		},
		"[]",
	)

	require.ErrorIs(t, e, BadToTheBone)
	require.Equal(t, "abc", v)
}

func Test_Thonk_callRecover_1(t *testing.T) {
	// When calling a function with recovery
	// where the function panics
	// then the panic value is returned

	_, _, r := callRecover(
		func() { panic("Moo") },
		"[]",
	)

	require.Equal(t, "Moo", r)
}
