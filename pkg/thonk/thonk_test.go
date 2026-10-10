package thonk

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Validate_1(t *testing.T) {
	// When validating a function
	// where f is not a function
	// then a named error is returned

	e := Validate(0)
	require.ErrorIs(t, e, ErrNotFunc)
}

func Test_Validate_2(t *testing.T) {
	// When validating a function
	// where f has too many output parameters
	// then a named error is returned

	f := func() (bool, int, error) {
		return true, 2, nil
	}

	e := Validate(f)
	require.ErrorIs(t, e, ErrTooManyOutputs)
}

func Test_Validate_3(t *testing.T) {
	// When validating a function
	// where f has a valid function signature
	// then no errors are returned

	var e error

	f0 := func() { return }
	e = Validate(f0)
	require.NoError(t, e)

	f1 := func() int { return 2 }
	e = Validate(f1)
	require.NoError(t, e)

	f2 := func() (int, error) { return 2, nil }
	e = Validate(f2)
	require.NoError(t, e)
}

func Test_WithNoArgs_1(t *testing.T) {
	// When creating a Thonk without arguments
	// where a valid function is passed
	// then a new Thonk is returned
	// with its function field set
	// and its arguments field empty (nil)

	thunk, e := WithNoArgs(func() {})

	require.NoError(t, e)
	require.NotNil(t, thunk.Func)
	require.Nil(t, thunk.Args)
}

func Test_WithJsonArgs_1(t *testing.T) {
	// When creating a Thonk with JSON arguments
	// where bad JSON is given
	// then a named error is returned

	_, e := WithJsonArgs(
		func() {},
		"",
	)

	require.ErrorIs(t, e, ErrBadJsonString)
}

func Test_WithJsonArgs_2(t *testing.T) {
	// When creating a Thonk with JSON arguments
	// that is non-variadic
	// where the JSON has less values than the amount
	// required by the function
	// then a named error is returned

	_, e := WithJsonArgs(
		func(v any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_WithJsonArgs_3(t *testing.T) {
	// When creating a Thonk with JSON arguments
	// that is non-variadic
	// where the JSON has more values than the amount
	// required by the function
	// then a named error is returned

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
	// When creating a Thonk with JSON arguments
	// that is variadic
	// where the JSON has less values than the amount of
	// non-variadic parameters required by the function
	// then a named error is returned

	_, e := WithJsonArgs(
		func(v any, more ...any) {},
		`
			[]
		`,
	)

	require.ErrorIs(t, e, ErrArgMismatch)
}

func Test_WithJsonArgs_5(t *testing.T) {
	// When creating a Thonk with JSON arguments
	// that is non-variadic
	// where the number of JSON values matches the amount of
	// inputs required by the function
	// then no error is returned

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
	// When creating a Thonk with JSON arguments
	// that is variadic
	// where the number of JSON values matches the amount of
	// non-variadic inputs required by the function
	// then no error is returned

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
	// When creating a Thonk with JSON arguments
	// that is non-variadic
	// where the JSON contains values and the function
	// accepts no parameeters
	// then no error is returned
	// and Thonk's argument field is empty

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
	// When creating a Thonk with JSON arguments
	// that is non-variadic
	// where the function requires typed inputs
	// and the JSON values match the types
	// then no error is returned
	// and Thonk's argument field contains the parsed JSON
	// values

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
	// When creating a Thonk with JSON arguments
	// that is variadic
	// where the JSON provides only the non-variadic
	// arguments
	// then no error is returned
	// and Thonk's argument field contains the parsed JSON
	// values

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
	// When creating a Thonk with JSON arguments
	// that is variadic
	// where the JSON provides several values for the
	// function's variadic input
	// then no error is returned
	// and Thonk's argument field contains the parsed JSON
	// values, including the variadic ones

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
	// there wasn't any issue with parsing structs and
	// arrays.

	// When creating a Thonk with JSON arguments
	// where the function accepts objects
	// the JSON provides the values for those objects
	// then no error is returned
	// and Thonk's argument field contains the parsed JSON
	// objects

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

func Test_Thonk_Call_1(t *testing.T) {
	// When calling a Thonk
	// with many input parameters
	// then the function is called with the arguments
	// and in the correct order

	var v1 int
	var v2 string
	var v3 []float64

	thunk := Thonk{
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

func Test_Thonk_Call_2(t *testing.T) {
	// When calling a Thonk
	// with many output parameters
	// then the function is called
	// and the output values are returned

	var BadToTheBone = errors.New("Bad to the bone")

	thunk := Thonk{
		Func: func() (string, error) {
			return "abc", BadToTheBone
		},
	}

	v, e := thunk.Call()

	require.ErrorIs(t, e, BadToTheBone)
	require.Equal(t, "abc", v)
}

func Test_Thonk_CallRecover_1(t *testing.T) {
	// When calling a Thonk that recovers
	// where the function panics
	// then the Thonk recovers
	// and the panic value returned

	thunk := Thonk{
		Func: func() {
			panic("Moo")
		},
	}

	_, _, r := thunk.CallRecover()
	require.Equal(t, "Moo", r)
}
