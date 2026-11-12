package sourcery

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

	// ErrOutputNotError occurs when validation fails because
	// the function's second output parameter is not an
	// error.
	ErrOutputNotError = sin.Template(
		"Expected an error as second output type but got %s",
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

// ValidateFunction returns nil if the passed function may
// be used with [Call]. It returns an error describing why
// if not.
//
// The function must return either nothing, a value (T), an
// error (error), or a value then an error (T, error).
func ValidateFunction(f any) error {
	typ := reflect.TypeOf(f)

	if typ.Kind() != reflect.Func {
		return ErrNotFunc.Fmt(typ.Kind())
	}

	if typ.NumOut() < 2 {
		return nil
	}

	if typ.NumOut() > 2 {
		return ErrTooManyOutputs.Fmt(typ.NumOut())
	}

	errType := typ.Out(1)
	if errType.Name() != "error" {
		return ErrOutputNotError.Fmt(errType.Name())
	}

	return nil
}

// callRecover is the same as [call] except it recovers
// from panics and returns the value as the third output.
func callRecover(f any, jsonArgs string) (v any, e error, r any) {
	defer func() {
		r = recover()
	}()

	v, e = call(f, jsonArgs)
	return v, e, nil
}

// call invokes the function f with the arguments provided
// as a JSON string. The function must return either
// nothing, a value (T), an error (error), or a value then
// an error (T, error).
//
// If f fails to validate via [ValidateFunction] then an
// error is returned. If jsonArgs cannot be passed into an
// array of arguments matching the order and types of the
// function parameters then an error is returned.
//
// If the function has 1 output that satisfies the error
// interface then it will be returned as both the value and
// the error.
func call(f any, jsonArgs string) (any, error) {
	if e := ValidateFunction(f); e != nil {
		return nil, e
	}

	args, e := parseJsonArgs(f, jsonArgs)
	if e != nil {
		return nil, e
	}

	val := reflect.ValueOf(f)
	res := val.Call(args)

	// Since we validated the function we know this will the
	// returned values will match one of these cases.
	switch len(res) {
	case 0:
		return nil, nil
	case 1:
		v := res[0].Interface()
		if err, ok := v.(error); ok {
			return v, err
		}
		return v, nil
	case 2:
		v := res[0].Interface()
		anyErr := res[1].Interface()
		if anyErr == nil {
			return v, nil
		}
		return v, anyErr.(error)
	default:
		panic("Sanity check! Unreachable code")
	}
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
