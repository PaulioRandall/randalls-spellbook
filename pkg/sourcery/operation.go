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
		return empty, ErrBadJsonString.Wrap(e)
	}

	// TODO: Create arg len checking func.
	if !typ.IsVariadic() && len(jsonArgs) != typ.NumIn() {
		return empty, curse.Fmt(
			"Expected %d arguments, given %d",
			typ.NumIn(),
			len(jsonArgs),
		).Wrap(ErrArgMismatch)
	}
	if typ.IsVariadic() && len(jsonArgs) < typ.NumIn()-1 {
		return empty, curse.Fmt(
			"Expected %d or more arguments, given %d",
			typ.NumIn(),
			len(jsonArgs),
		).Wrap(ErrArgMismatch)
	}

	// NEXT: iterate raw messages,
	//       unmarshal each creating new container using
	//       reflect.New (creates pointer to value) then
	//       pass Value.Interface() to json.Unmarshal,
	//       REMEMBER to call Value.Elem() to get value
	//       being pointed at afterwards!!

	return empty, nil
}
