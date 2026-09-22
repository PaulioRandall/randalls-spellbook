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

type ValErrFunc struct {
	Func any
	Args []reflect.Value
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

func WithNoArgs(f any) (ValErrFunc, error) {
	var result ValErrFunc

	e := ValidateValErrFunc(f)
	if e != nil {
		return result, e
	}

	result.Func = f
	return result, nil
}

func WithJsonArgs(f any, jsonStr string) (ValErrFunc, error) {
	var empty ValErrFunc

	vet, e := WithNoArgs(f)
	if e != nil {
		return empty, e
	}

	return vet.WithJsonArgs(jsonStr)
}

// WithNoArgs returns a copy of the ValErrFunc with all
// arguments removed.
func (vet ValErrFunc) WithNoArgs() ValErrFunc {
	vet.Args = nil
	return vet
}

// WithJsonArgs returns a copy of the ValErrThink with
// arguments replaced by those parsed from the passed
// jsonStr. The JSON must be an array of values that map to
// the function's parameters. Ordering matters!
func (vet ValErrFunc) WithJsonArgs(
	jsonStr string,
) (ValErrFunc, error) {
	empty := ValErrFunc{}

	args, e := parseJsonArgs(vet.Func, jsonStr)
	if e != nil {
		return empty, e
	}

	vet.Args = args
	return vet, nil
}

func parseJsonArgs(
	f any,
	jsonStr string,
) ([]reflect.Value, error) {
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
) ([]reflect.Value, error) {
	args := make([]reflect.Value, len(jsonArgs), len(jsonArgs))

	for i, arg := range jsonArgs {
		val := createPointerValueToParam(funcTyp, i)

		e := json.Unmarshal(arg, val.Interface())
		if e == nil {
			// Elem() because we're using ValueOf pointer.
			args[i] = val.Elem()
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

// Call invokes the thunk with its set arguments.
// If the function has 1 output that satisfies the error
// interface then it will be returned as both the value and
// the error. Call doesn't recover from panics.
func (vet ValErrFunc) Call() (any, error) {
	val := reflect.ValueOf(vet.Func)
	results := val.Call(vet.Args)

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

// CallRecover does the same as Call except it recovers
// from a panic and returns the recovered value as a third
// return value.
func (vet ValErrFunc) CallRecover() (v any, e error, r any) {
	defer func() {
		r = recover()
	}()

	v, e = vet.Call()
	return v, e, nil
}
