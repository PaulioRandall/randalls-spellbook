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

type ValErrThunk struct {
	Func any
	Args []any
}

func ValidateValErrFunc(f any) error {
	typ := reflect.TypeOf(f)

	if typ.Kind() != reflect.Func {
		return ErrNotFunc.Fmt(typ.Kind())
	}

	if typ.NumOut() > 2 {
		return ErrTooManyOutputs.Fmt(typ.NumOut())
	}

	return nil
}

func NoArgs(f any) (ValErrThunk, error) {
	var result ValErrThunk

	e := ValidateValErrFunc(f)
	if e != nil {
		return result, e
	}

	result.Func = f
	return result, nil
}

func WithJsonArgs(f any, jsonStr string) (ValErrThunk, error) {
	var empty ValErrThunk

	vet, e := NoArgs(f)
	if e != nil {
		return empty, e
	}

	return vet.WithJsonArgs(jsonStr)
}

// WithNoArgs returns a copy of the ValErrThunk with all
// arguments removed.
func (vet ValErrThunk) NoArgs() ValErrThunk {
	vet.Args = nil
	return vet
}

// WithJsonArgs returns a copy of the ValErrThink with
// arguments replaced by those parsed from the passed
// jsonStr. The JSON must be an array of values that map to
// the function's parameters. Ordering matters!
func (vet ValErrThunk) WithJsonArgs(jsonStr string) (ValErrThunk, error) {
	empty := ValErrThunk{}

	args, e := parseJsonArgs(vet.Func, jsonStr)
	if e != nil {
		return empty, e
	}

	vet.Args = args
	return vet, nil
}

func parseJsonArgs(f any, jsonStr string) ([]any, error) {
	typ := reflect.TypeOf(f)

	var jsonArgs []json.RawMessage
	e := json.Unmarshal([]byte(jsonStr), &jsonArgs)
	if e != nil {
		return nil, ErrBadJsonString.Wraps(e)
	}

	e = validateArgsLength(typ, len(jsonArgs))
	if e != nil {
		return nil, ErrArgMismatch.Wraps(e)
	}

	if len(jsonArgs) == 0 {
		return nil, nil
	}

	args, e := parseJsonMessages(typ, jsonArgs)
	if e != nil {
		return nil, ErrBadJsonString.Wraps(e)
	}

	return args, nil
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

func parseJsonMessages(
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
