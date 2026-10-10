package thonk

import (
	"encoding/json"
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

var (
	// ErrNotFunc occurs when validation fails because the
	// passed value being validated is not a function.
	ErrNotFunc = sin.Template(
		"Expected a function but given %s",
	)

	// ErrTooManyOutputs occurs when validation fails because
	// the function has too many outputs.
	ErrTooManyOutputs = sin.Template(
		"Function has too many outputs: given %d but want either nothing, a value (T), an error (error), or a value then an error (T, error)",
	)

	// ErrBadJsonString occurs when attempting to parse
	// arguments provided as a JSON array but the JSON is
	// invalid.
	ErrBadJsonString = sin.Err(
		"Failed to parse JSON string",
	)

	// ErrArgMismatch occurs when the number of arguments
	// being set does not match the number of input
	// parameters the function expects.
	ErrArgMismatch = sin.Err(
		"Mismatch between expected and given arguments",
	)
)

// Validate checks if the passed function may be used with
// a [Thonk]. It returns nil if it can, and an error if
// not.
func Validate(f any) error {
	typ := reflect.TypeOf(f)

	if typ.Kind() != reflect.Func {
		return ErrNotFunc.Fmt(typ.Kind())
	}

	if typ.NumOut() > 2 {
		return ErrTooManyOutputs.Fmt(typ.NumOut())
	}

	return nil
}

// Thonk is a configurable Thunk limited to functions that
// return either nothing, a single value (T), and error
// (error), or a value then an error (T, error).
type Thonk struct {
	// The function to be called.
	Func any

	// The input arguments to the function.
	Args []reflect.Value
}

// WithNoArgs returns a new Thonk with no arguments. An
// error is returned if f fails validation.
func WithNoArgs(f any) (Thonk, error) {
	var result Thonk

	e := Validate(f)
	if e != nil {
		return result, e
	}

	result.Func = f
	return result, nil
}

// WithJsonArgs returns a new Thonk with arguments parsed
// from a JSON array, provided as a string. An error is
// returned if f fails validation, jsonStr cannot be
// parsed, or there's a mismatch between arguments and
// the function's inputs.
func WithJsonArgs(f any, jsonStr string) (Thonk, error) {
	var empty Thonk

	th, e := WithNoArgs(f)
	if e != nil {
		return empty, e
	}

	return th.WithJsonArgs(jsonStr)
}

// WithNoArgs returns a copy of the Thonk with all
// arguments removed.
func (th Thonk) WithNoArgs() Thonk {
	th.Args = nil
	return th
}

// WithJsonArgs returns a copy of the Thonk with arguments
// parsed from a JSON array, provided as a string. An error
// is returned if jsonStr cannot be parsed or there's a
// mismatch between arguments and the function's inputs.
func (th Thonk) WithJsonArgs(jsonStr string) (Thonk, error) {
	empty := Thonk{}

	args, e := parseJsonArgs(th.Func, jsonStr)
	if e != nil {
		return empty, e
	}

	th.Args = args
	return th, nil
}

func parseJsonArgs(
	f any,
	jsonStr string,
) ([]reflect.Value, error) {
	typ := reflect.TypeOf(f)

	var jsonArgs []json.RawMessage
	e := json.Unmarshal([]byte(jsonStr), &jsonArgs)
	if e != nil {
		return nil, ErrBadJsonString.Wrap(e)
	}

	e = validateArgsLength(typ, len(jsonArgs))
	if e != nil {
		return nil, ErrArgMismatch.Wrap(e)
	}

	if len(jsonArgs) == 0 {
		return nil, nil
	}

	args, e := parseJsonMessages(typ, jsonArgs)
	if e != nil {
		return nil, ErrBadJsonString.Wrap(e)
	}

	return args, nil
}

func validateArgsLength(typ reflect.Type, size int) error {
	if !typ.IsVariadic() && size != typ.NumIn() {
		return sin.Fmt(
			"Expected %d arguments, given %d",
			typ.NumIn(),
			size,
		)
	}

	if typ.IsVariadic() && size < typ.NumIn()-1 {
		return sin.Fmt(
			"Expected %d or more arguments, given %d",
			typ.NumIn(),
			size,
		)
	}

	return nil
}

func parseJsonMessages(
	funcTyp reflect.Type,
	jsonArgs []json.RawMessage,
) ([]reflect.Value, error) {
	args := make([]reflect.Value, len(jsonArgs))

	for i, arg := range jsonArgs {
		val := createPointerValueToParam(funcTyp, i)

		e := json.Unmarshal(arg, val.Interface())
		if e == nil {
			// Elem() because we're using ValueOf pointer.
			args[i] = val.Elem()
			continue
		}

		return nil, sin.Fmt(
			"JSON array argument index %d",
			i,
		).Wrap(e)
	}

	return args, nil
}

func createPointerValueToParam(
	funcTyp reflect.Type,
	idx int,
) reflect.Value {
	lastParamIdx := funcTyp.NumIn() - 1
	var paramTyp reflect.Type

	if funcTyp.IsVariadic() && idx >= lastParamIdx {
		paramTyp = funcTyp.In(lastParamIdx).Elem()
	} else {
		paramTyp = funcTyp.In(idx)
	}

	return reflect.New(paramTyp)
}

// Call invokes the thunk with its set arguments. If the
// function has 1 output that satisfies the error interface
// then it will be returned as both the value and the
// error. Call doesn't recover from panics.
func (th Thonk) Call() (any, error) {
	val := reflect.ValueOf(th.Func)
	results := val.Call(th.Args)

	switch len(results) {
	case 0:
		return nil, nil
	case 1:
		v := results[0].Interface()
		if err, ok := v.(error); ok {
			return v, err
		}
		return v, nil
	default: // 2 return values
		// Type check done during parsing.
		err, _ := results[1].Interface().(error)
		return results[0].Interface(), err
	}
}

// CallRecover does the same as [Thonk.Call] except it
// recovers from a panic and returns the recovered value as
// a third return value. In most cases the user should
// check if the third recover value is nil before checking
// if the error is nil.
func (th Thonk) CallRecover() (v any, e error, r any) {
	defer func() {
		r = recover()
	}()

	v, e = th.Call()
	return v, e, nil
}
