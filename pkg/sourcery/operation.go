package sourcery

import (
	"encoding/json"
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

var (
	ErrNotFunc = curse.Proto(
		"Expected a function but given %s",
	)

	ErrTooManyOutputs = curse.Proto(
		"Too many output parameters, given %d: either return nothing, a value (T), an error (error), or a value then an error (T, error)",
	)

	ErrBadJsonString = curse.Err(
		"Failed to parse JSON string",
	)

	ErrArgMismatch = curse.Err(
		"Mismatch between expected and given arguments",
	)
)

type Operation struct {
	Func any
	Args []any
}

func NewOperation(f any) (Operation, error) {
	result := Operation{}
	typ := reflect.TypeOf(f)

	if typ.Kind() != reflect.Func {
		return result, ErrNotFunc.Fmt(typ.Kind())
	}

	if typ.NumOut() > 2 {
		return result, ErrTooManyOutputs.Fmt(typ.NumOut())
	}

	result.Func = f
	return result, nil
}

// WithJsonArgs accepts a JSON string containing an array
// of values to be used as operation arguments.
func (op Operation) WithJsonArgs(
	jsonStr string,
) (Operation, error) {
	empty := Operation{}
	typ := reflect.TypeOf(op.Func)

	var jsonArgs []json.RawMessage
	e := json.Unmarshal([]byte(jsonStr), &jsonArgs)
	if e != nil {
		return empty, ErrBadJsonString.Wraps(e)
	}

	e = validateArgsLength(typ, len(jsonArgs))
	if e != nil {
		return empty, ErrArgMismatch.Wraps(e)
	}

	op.Args, e = parseJsonArgs(typ, jsonArgs)
	if e != nil {
		return empty, ErrBadJsonString.Wraps(e)
	}

	return op, nil
}

func validateArgsLength(typ reflect.Type, size int) error {
	if !typ.IsVariadic() && size != typ.NumIn() {
		return curse.Fmt(
			"Expected %d arguments, given %d",
			typ.NumIn(),
			size,
		)
	}

	if typ.IsVariadic() && size < typ.NumIn()-1 {
		return curse.Fmt(
			"Expected %d or more arguments, given %d",
			typ.NumIn(),
			size,
		)
	}

	return nil
}

func parseJsonArgs(
	funcTyp reflect.Type,
	jsonArgs []json.RawMessage,
) ([]any, error) {
	args := make([]any, len(jsonArgs), len(jsonArgs))

	for i, arg := range jsonArgs {
		val := createPointerValueToParam(funcTyp, i)

		e := json.Unmarshal(arg, val.Interface())
		if e == nil {
			// Elem() because we're using ValueOf pointer.
			args[i] = val.Elem().Interface()
			continue
		}

		return nil, curse.Fmt(
			"JSON array argument index %d",
			i,
		).Wraps(e)
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

func (op Operation) Invoke() (any, error) {
	// NEXT: Call and return output based upon func outputs.
	return nil, nil
}
